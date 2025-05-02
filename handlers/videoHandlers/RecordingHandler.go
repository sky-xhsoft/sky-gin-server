// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: RecordingHandler.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/17
// Project Description:
// ----------------------------------------------------------------------------

package videoHandlers

import (
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/sky-xhsoft/sky-gin-server/core"
	"github.com/sky-xhsoft/sky-gin-server/handlers"
	"github.com/sky-xhsoft/sky-gin-server/models"
	"github.com/sky-xhsoft/sky-gin-server/pkg/ecode"
	"github.com/sky-xhsoft/sky-gin-server/pkg/utils"
	"gorm.io/gorm"
	"io"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type RecordingHandler struct {
	db *gorm.DB
}

func (h *RecordingHandler) HandlerName() string {
	return "RecordingHandler"
}

func init() {
	handlers.Register("RecordingHandler", &RecordingHandler{})
}

func (h *RecordingHandler) SetOption(ctx *core.AppContext) {
	h.db = ctx.DB
}

var recordingProcesses = make(map[uint]*exec.Cmd) // 资源ID -> 进程
var recordingProcessesLock sync.Mutex
var recordingProcessesStdin = make(map[uint]io.WriteCloser)

var recordingProcessesCancel = make(map[uint]context.CancelFunc)

// StartCut 启动切片任务，且支持使用 context 取消两个同步协程
func (h *RecordingHandler) StartRecording(c *gin.Context) {
	resourceID := c.Query("resourceId")
	if resourceID == "" {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少资源组ID")
		return
	}

	var item models.ChrResourceItem
	if err := h.db.Where("CHR_RESOURCE_ID = ? AND TYPE = ? AND IS_ACTIVE = 'Y'", resourceID, "RTMP").First(&item).Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrRequest, "未找到资源")
		return
	}

	outputDir := fmt.Sprintf("tmp/outcut/all/%d", *item.ChrResourceId)
	outputTemplate := filepath.Join(outputDir, "ALL_%Y-%m-%d_%H-%M-%S.mp4")
	if err := utils.EnsureDir(outputDir); err != nil {
		ecode.Resp(c, ecode.ErrServer, "无法创建输出目录")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	recordingProcessesCancel[*item.ChrResourceId] = cancel

	go func(rtmpUrl, outputTemplate, cutTime string, rid *uint, pid uint, db *gorm.DB, ctx context.Context) {
		defer func() {
			recordingProcessesLock.Lock()
			delete(recordingProcesses, *rid)
			delete(recordingProcessesCancel, *rid)
			delete(recordingProcessesStdin, *rid)
			recordingProcessesLock.Unlock()
			h.db.Model(&models.ChrResource{}).Where("ID = ?", *rid).Update("CUT_STATUS", 0)
		}()

		log.Println("开始直播切片")
		cmd := exec.Command("ffmpeg",
			"-i", rtmpUrl,
			"-c", "copy",
			"-c:a", "aac", // 强制音频重编码为 AAC
			"-f", "segment",
			"-segment_time", "10800",
			"-reset_timestamps", "1",
			"-strftime", "1",
			outputTemplate,
		)
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

		if err := cmd.Start(); err != nil {
			log.Printf("ffmpeg 启动失败: %v", err)
			return
		}

		h.db.Model(&models.ChrResource{}).Where("ID = ?", *rid).Update("RECORDING_STATUS", 1)
		recordingProcessesLock.Lock()
		recordingProcesses[*rid] = cmd
		recordingProcessesLock.Unlock()

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Println("创建监控失败: %v", err)
			return
		}
		defer watcher.Close()

		if err := watcher.Add(outputDir); err != nil {
			log.Println("监控目录失败: %v", err)
			return
		}

		var pendingFiles []string
		var pendingMu sync.Mutex

		//监控文件是否稳定并上传
		go func() {
			ticker := time.NewTicker(60 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					//上传完成所有文件后退出
					pendingMu.Lock()
					var remain []string
					for _, f := range pendingFiles {
						if isFileStable(f, 3) {
							go uploadFileToOSS(f, rid, pid, db, c)
						} else {
							remain = append(remain, f)
						}
					}
					pendingFiles = remain
					pendingMu.Unlock()
					log.Println("[切片任务] 上传监控退出")
					return
				case <-ticker.C:
					pendingMu.Lock()
					var remain []string
					for _, f := range pendingFiles {
						if isFileStable(f, 4) {
							go uploadFileToOSS(f, rid, pid, db, c)
						} else {
							remain = append(remain, f)
						}
					}
					pendingFiles = remain
					pendingMu.Unlock()
				}
			}
		}()

		for {
			select {
			case <-ctx.Done():
				log.Println("[切片任务] 监控退出")
				return
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create && strings.HasSuffix(event.Name, ".mp4") {
					pendingMu.Lock()
					pendingFiles = append(pendingFiles, event.Name)
					pendingMu.Unlock()
				}
			case err := <-watcher.Errors:
				log.Printf("监控错误: %v", err)
			}
		}
	}(item.RtmpUrl, outputTemplate, strconv.Itoa(*item.CutTimes), item.ChrResourceId, item.ProjectId, h.db, ctx)

	ecode.SuccessResp(c, "切片任务已启动")
}

// 终止直播切片
func (h *RecordingHandler) StopRecording(c *gin.Context) {
	resourceID := c.Query("resourceId")
	if resourceID == "" {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少资源组ID")
		return
	}
	rid, err := strconv.ParseUint(resourceID, 10, 64)
	if err != nil {
		ecode.Resp(c, ecode.ErrInvalidParam, "资源组ID非法")
		return
	}

	recordingProcessesLock.Lock()
	cmd, exists := recordingProcesses[uint(rid)]
	cancelFunc, cancelExists := recordingProcessesCancel[uint(rid)]
	recordingProcessesLock.Unlock()

	if exists && cmd != nil && cmd.Process != nil {
		_ = syscall.Kill(cmd.Process.Pid, syscall.SIGINT)
	}

	if cancelExists {
		time.Sleep(1 * time.Second)
		cancelFunc() // 通知同步协程退出
	}

	h.db.Model(&models.ChrResource{}).
		Where("ID = ?", uint(rid)).
		Update("RECORDING_STATUS", 0)

	ecode.SuccessResp(c, "切片任务已停止")
}
