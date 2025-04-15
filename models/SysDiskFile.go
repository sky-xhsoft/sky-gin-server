// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: SysDiskFile.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/14
// Project Description:
// ----------------------------------------------------------------------------

package models

type SysDiskFile struct {
	ProjectId uint   `gorm:"column:PROJECT_ID" json:"projectId"`
	ParentId  uint   `gorm:"column:PARENT_ID" json:"parentId"`
	FileName  string `gorm:"column:FILE_NAME" json:"fileName"`
	FileType  string `gorm:"column:FILE_TYPE" json:"fileType"` // F: 文件, D: 目录
	FilePath  string `gorm:"column:FILE_PATH" json:"filePath"`
	FileExt   string `gorm:"column:FILE_EXT" json:"fileExt"`
	FileSize  int64  `gorm:"column:FILE_SIZE" json:"fileSize"`
	MimeType  string `gorm:"column:MIME_TYPE" json:"mimeType"`
	Hash      string `gorm:"column:HASH" json:"hash"`
	IsShared  string `gorm:"column:IS_SHARED" json:"isShared"`

	Base
}

func (SysDiskFile) TableName() string {
	return "sys_disk_file"
}
