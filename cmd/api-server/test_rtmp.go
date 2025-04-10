// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: test_rtmp.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/8
// Project Description:
// ----------------------------------------------------------------------------
package main

import (
	"log"
	"os"
	"os/exec"
)

func main_trmp() {
	// 你的 RTMP 地址
	rtmpURL := "rtmp://pull-hssh.vzan.com/v/94041117_694257248267285199?zbid=94041117&tpid=157700674"

	// 输出文件模板（会自动生成多个分段）
	outputTemplate := "recordings/%Y-%m-%d_%H-%M-%S.mp4"

	// 创建输出目录
	os.MkdirAll("recordings", 0755)

	// 构造 ffmpeg 命令
	cmd := exec.Command("ffmpeg",
		"-i", rtmpURL, // 输入流
		"-c", "copy", // 不重新编码，节省资源
		"-f", "segment", // 使用分段模式
		"-segment_time", "60", // 每 60 秒一个切片
		"-reset_timestamps", "1", // 重置时间戳，避免播放问题
		"-strftime", "1", // 启用时间戳命名
		outputTemplate, // 输出文件名模板
	)

	// 输出 ffmpeg 日志到终端
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Println("开始拉流并切片...")

	// 运行 ffmpeg 进程
	if err := cmd.Run(); err != nil {
		log.Fatalf("ffmpeg 执行失败: %v", err)
	}

	log.Println("录制结束。")
}
