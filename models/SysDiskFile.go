// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: SysDiskFile.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/14
// Project Description:
// ----------------------------------------------------------------------------

package models

type SysDiskFile struct {
	ID           uint   `gorm:"primaryKey;column:ID"`
	SysCompanyId uint   `gorm:"column:SYS_COMPANY_ID"`
	ProjectId    uint   `gorm:"column:PROJECT_ID"`
	ParentId     uint   `gorm:"column:PARENT_ID"`
	FileName     string `gorm:"column:FILE_NAME"`
	FileType     string `gorm:"column:FILE_TYPE"` // F or D
	FilePath     string `gorm:"column:FILE_PATH"`
	FileExt      string `gorm:"column:FILE_EXT"`
	FileSize     int64  `gorm:"column:FILE_SIZE"`
	MimeType     string `gorm:"column:MIME_TYPE"`
	Hash         string `gorm:"column:HASH"`
	IsShared     string `gorm:"column:IS_SHARED"`

	Base
}

func (SysDiskFile) TableName() string {
	return "sys_disk_file"
}
