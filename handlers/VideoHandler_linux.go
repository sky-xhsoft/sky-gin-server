// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: VideoHandler_windows.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/15
// Project Description:
// ----------------------------------------------------------------------------

//go:build linux
// +build linux

package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/sky-xhsoft/sky-gin-server/core"
	"github.com/sky-xhsoft/sky-gin-server/models"
	"github.com/sky-xhsoft/sky-gin-server/pkg/ecode"
	"github.com/sky-xhsoft/sky-gin-server/pkg/ossUtil"
	"github.com/sky-xhsoft/sky-gin-server/pkg/utils"
	"gorm.io/gorm"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type VideoHandler struct {
	db *gorm.DB
}

func (h *VideoHandler) HandlerName() string {
	return "VideoHandler"
}

func init() {
	Register("VideoHandler", &VideoHandler{})
}

func (h *VideoHandler) SetOption(ctx *core.AppContext) {
	h.db = ctx.DB
}

var cutProcesses = make(map[uint]*exec.Cmd) // 资源ID -> 进程
var cutProcessesLock sync.Mutex

var processedFiles = make(map[string]bool)
var mu sync.Mutex // 保证并发安全

// StartCut 启动直播流分段录制并监听目录上传切片至 OSS
func (h *VideoHandler) StartCut(c *gin.Context) {
	//tx := utils.GetTx(c, h.db)
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

	outputDir := fmt.Sprintf("tmp/outcut/%d", *item.ChrResourceId)
	outputTemplate := filepath.Join(outputDir, "%Y-%m-%d_%H-%M-%S.mp4")

	if err := utils.EnsureDir(outputDir); err != nil {
		ecode.Resp(c, ecode.ErrServer, "无法创建输出目录")
		return
	}

	go func(rtmpUrl, outputTemplate, cutTime string, rid *uint, pid uint, db *gorm.DB) {

		defer func() {
			cutProcessesLock.Lock()
			delete(cutProcesses, *rid)
			cutProcessesLock.Unlock()
			h.db.Model(&models.ChrResource{}).
				Where("ID = ?", *item.ChrResourceId).
				Update("CUT_STATUS", 0)
		}()

		log.Println("开始直播切片")

		cmd := exec.Command("ffmpeg",
			"-i", rtmpUrl,
			"-c", "copy",
			"-f", "segment",
			"-segment_time", cutTime,
			"-reset_timestamps", "1",
			"-strftime", "1",
			outputTemplate,
		)

		// 为进程设置独立的进程组，便于后续整体杀死
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

		if err := cmd.Start(); err != nil {
			log.Printf("ffmpeg 启动失败: %v", err)
			return
		}
		h.db.Model(&models.ChrResource{}).
			Where("ID = ?", *item.ChrResourceId).
			Update("CUT_STATUS", 1)

		cutProcessesLock.Lock()
		cutProcesses[*rid] = cmd
		cutProcessesLock.Unlock()

		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Printf("创建文件监控失败: %v", err)
			return
		}
		defer watcher.Close()

		if err := watcher.Add(outputDir); err != nil {
			log.Printf("监听目录失败: %v", err)
			return
		}

		//var lastFilePath string
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.CloseWrite == fsnotify.CloseWrite && strings.HasSuffix(event.Name, ".mp4") {
					log.Println("新切片创建：" + event.Name)

					mu.Lock()
					if processedFiles[event.Name] {
						mu.Unlock()
						continue
					}
					processedFiles[event.Name] = true
					mu.Unlock()

					// 上传前一个文件
					//if lastFilePath != "" {
					go uploadFileToOSS(event.Name, rid, pid, db, c)
					//}
					//lastFilePath = event.Name
				}
			case err := <-watcher.Errors:
				log.Printf("监控错误: %v", err)
			}
		}

	}(item.RtmpUrl, outputTemplate, strconv.Itoa(*item.CutTimes), item.ChrResourceId, item.ProjectId, h.db)

	ecode.SuccessResp(c, "切片任务已启动")
}

func uploadFileToOSS(filePath string, rid *uint, pid uint, db *gorm.DB, c *gin.Context) {
	log.Printf("准备上传上一个切片: %s", filePath)

	f, err := os.Open(filePath)
	if err != nil {
		log.Printf("打开切片失败: %v", err)
		return
	}
	defer f.Close()

	ossClient := ossUtil.GetClient()
	fileRecord, err := ossClient.UploadLocalFile(f, filepath.Base(filePath), rid)
	if err != nil {
		log.Printf("切片上传失败: %v", err)
		return
	}

	param, _ := json.Marshal(fileRecord)
	sizeInMB := float64(fileRecord.FileSize) / 1024.0 / 1024.0
	sizeRounded := math.Round(sizeInMB*10) / 10

	item := models.ChrResourceItem{
		ChrResourceId: rid,
		ProjectId:     pid,
		Name:          filepath.Base(filePath),
		Type:          "VIDEO",
		VideoUrl:      fileRecord.OSSURL,
		VideoFileSize: sizeRounded,
		VideoFileType: fileRecord.FileType,
		VideoParam:    string(param),
	}
	models.FillCreateMeta(c, &item) // context 可选
	db.Create(&item)
}

func (h *VideoHandler) StopCut(c *gin.Context) {
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

	var item models.ChrResourceItem
	if err := h.db.Where("CHR_RESOURCE_ID = ? AND TYPE = ? AND IS_ACTIVE = 'Y'", resourceID, "RTMP").First(&item).Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrRequest, "未找到资源")
		return
	}

	cutProcessesLock.Lock()
	cmd, exists := cutProcesses[*item.ChrResourceId]
	cutProcessesLock.Unlock()

	if !exists || cmd == nil || cmd.Process == nil {
		ecode.Resp(c, ecode.ErrRequest, "未找到运行中的切片任务")
		return
	}

	// 杀掉整个进程组，负号表示进程组
	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
		ecode.Resp(c, ecode.ErrServer, "停止切片失败: "+err.Error())
		return
	}

	h.db.Model(&models.ChrResource{}).
		Where("ID = ?", *item.ChrResourceId).
		Update("CUT_STATUS", 0)

	// 清理映射
	cutProcessesLock.Lock()
	delete(cutProcesses, uint(rid))
	cutProcessesLock.Unlock()

	ecode.SuccessResp(c, "切片任务已停止")
}

func waitForCompleteWrite(path string, checkInterval time.Duration, maxWait time.Duration) bool {
	var prevSize int64 = -1
	start := time.Now()

	for time.Since(start) < maxWait {
		fi, err := os.Stat(path)
		if err != nil {
			return false
		}
		currSize := fi.Size()
		if currSize > 0 && currSize == prevSize {
			return true // 文件大小稳定
		}
		prevSize = currSize
		time.Sleep(checkInterval)
		start = time.Now()
	}
	return false
}

// GenerateHLS 调用 ffmpeg 生成 HLS 并返回播放路径（支持 html 播放）
func (h *VideoHandler) GenerateHLS(c *gin.Context) {
	tx := utils.GetTx(c, h.db)
	resourceID := c.Query("resourceId")
	if resourceID == "" {
		ecode.Resp(c, ecode.ErrInvalidParam, "缺少资源组ID")
		return
	}

	var item models.ChrResourceItem
	if err := tx.Where("CHR_RESOURCE_ID = ? AND TYPE = ? AND IS_ACTIVE = 'Y'", resourceID, "RTMP").First(&item).Error; err != nil {
		c.Error(err)
		ecode.Resp(c, ecode.ErrRequest, "未找到资源")
		return
	}

	// 使用资源ID生成唯一播放文件名
	streamDir := fmt.Sprintf("./tmp/hls_output/%d", *item.ChrResourceId)
	m3u8Path := filepath.Join(streamDir, "stream.m3u8")

	// 创建输出目录
	if err := utils.EnsureDir(streamDir); err != nil {
		ecode.Resp(c, ecode.ErrServer, "无法创建输出目录")
		return
	}

	go func(rtmpUrl, outputDir, outputM3U8 string) {
		log.Println("开始直播转码")
		cmd := exec.Command("ffmpeg",
			"-fflags", "nobuffer",
			"-flags", "low_delay",
			"-strict", "experimental",
			"-analyzeduration", "0",
			"-probesize", "32",
			"-i", rtmpUrl,
			"-preset", "ultrafast",
			"-tune", "zerolatency",
			"-g", "25",
			"-sc_threshold", "0",
			"-c:v", "libx264",
			"-crf", "18", // 高质量（保留分辨率+尽量无损）
			"-c:a", "aac",
			"-ar", "44100",
			"-f", "hls",
			"-hls_time", "1", // 每片1秒
			"-hls_list_size", "10",
			"-hls_flags", "delete_segments+split_by_time", // 删除旧片段
			"-hls_segment_filename", filepath.Join(outputDir, "seg%d.ts"),
			outputM3U8,
		)

		if err := cmd.Run(); err != nil {
			fmt.Printf("HLS 异步转码失败: %v\n", err)
		} else {
			log.Println("HLS 转码完成:", outputM3U8)
		}
	}(item.RtmpUrl, streamDir, m3u8Path)

	ecode.SuccessResp(c, gin.H{
		"hls":     fmt.Sprintf("/static/hls_output/%d/stream.m3u8", *item.ChrResourceId),
		"preview": fmt.Sprintf("/api/resource/play?id=%d", *item.ChrResourceId),
	})
}

// 预览直播效果
func (h *VideoHandler) PlayPage(c *gin.Context) {
	id := c.Query("id")

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="utf-8">
		<title>HLS 播放</title>
	</head>
	<body>
		<video id="video"  controls autoplay muted></video>
		<script src="https://cdn.jsdelivr.net/npm/hls.js@latest"></script>
		<script>
			if(Hls.isSupported()) {
				var video = document.getElementById('video');
				var hls = new Hls();
				hls.loadSource("/static/hls_output/%s/stream.m3u8");
				hls.attachMedia(video);
				hls.on(Hls.Events.MANIFEST_PARSED, function() {
					video.play();
				});
			} else if (video.canPlayType('application/vnd.apple.mpegurl')) {
				video.src = "/static/hls_output/%s/stream.m3u8";
				video.addEventListener('loadedmetadata', function() {
					video.play();
				});
			}
		</script>
	</body>
	</html>`, id, id))
}
