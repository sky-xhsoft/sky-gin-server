// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: ossUtil.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/14
// Project Description:
// ----------------------------------------------------------------------------

package ossUtil

import (
	"context"
	"fmt"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/gin-gonic/gin"
	"github.com/sky-xhsoft/sky-gin-server/config"
	"gorm.io/gorm"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

type OSSClient struct {
	client  *oss.Client
	config  *oss.Config
	host    string
	bucket  string
	baseUrl string
}

// FileRecord 文件记录模型
type FileRecord struct {
	gorm.Model
	FileName     string `gorm:"size:255"`
	FileKey      string `gorm:"size:255;uniqueIndex"`
	FileSize     int64
	FileType     string `gorm:"size:100"`
	OSSURL       string `gorm:"size:500"`
	UploadStatus string `gorm:"size:50"` // uploading, completed, failed
	PartNumber   int    // 分片序号，普通上传为0
	UploadID     string `gorm:"size:255"` // 分片上传ID
	ETag         string `gorm:"size:255"` // 文件ETag
}

var defaultClient *OSSClient

// Init 初始化默认 OSS 客户端实例（使用配置）
func Init(cfg *config.Config) error {
	// 初始化OSS客户端
	config := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.Oss.AccessKeyId, cfg.Oss.AccessKeySecret)).
		WithRegion(cfg.Oss.Region)

	client := oss.NewClient(config)

	defaultClient = &OSSClient{
		client:  client,
		config:  config,
		bucket:  cfg.Oss.BucketName,
		baseUrl: cfg.Oss.BaseUrl,
	}
	return nil
}

// GetClient 获取默认客户端
func GetClient() *OSSClient {
	return defaultClient
}

// UploadFile 上传文件到OSS（小文件直接上传）
func (o *OSSClient) UploadFile(ctx *gin.Context, fileHeader *multipart.FileHeader, customKey ...string) (*FileRecord, error) {
	// 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// 生成文件key
	fileKey := o.generateFileKey(fileHeader.Filename, customKey...)

	// 创建OSS请求
	putObjectRequest := &oss.PutObjectRequest{
		Bucket: &o.bucket,
		Key:    &fileKey,
		Body:   file,
	}

	// 执行上传
	_, err = o.client.PutObject(context.TODO(), putObjectRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to OSS: %v", err)
	}

	// 构建文件URL
	fileURL := fmt.Sprintf("https://%s/%s", o.baseUrl, fileKey)

	// 保存文件记录到数据库
	fileRecord := &FileRecord{
		FileName:     fileHeader.Filename,
		FileKey:      fileKey,
		FileSize:     fileHeader.Size,
		FileType:     fileHeader.Header.Get("Content-Type"),
		OSSURL:       fileURL,
		UploadStatus: "completed",
		PartNumber:   0, // 0表示普通上传
	}

	return fileRecord, nil
}

// generateFileKey 生成文件在OSS中的key
func (o *OSSClient) generateFileKey(originalName string, customKey ...string) string {
	ext := filepath.Ext(originalName)
	_ = strings.TrimSuffix(originalName, ext)

	if len(customKey) > 0 && customKey[0] != "" {
		return fmt.Sprintf("%s/%d%s", customKey[0], time.Now().UnixNano(), ext)
	}

	return fmt.Sprintf("uploads/%d/%d%s", time.Now().Year(), time.Now().UnixNano(), ext)
}
