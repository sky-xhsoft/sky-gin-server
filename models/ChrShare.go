// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: ChrShare.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/19
// Project Description:
// ----------------------------------------------------------------------------

package models

type ChrShare struct {
	Base
	Key               string `gorm:"column:KEY;type:varchar(255);uniqueIndex" json:"key"`
	Password          string `gorm:"column:PASSWORD;type:varchar(255)" json:"password"`
	ChrProjectID      *uint  `gorm:"column:CHR_PROJECT_ID;type:int" json:"chrProjectId"`
	ChrResourceID     *uint  `gorm:"column:CHR_RESOURCE_ID;type:int" json:"chrResourceId"`
	ChrResourceItemID *uint  `gorm:"column:CHR_RESOURCE_ITEM_ID;type:int" json:"chrResourceItemId"`
	SysDiskFileID     *uint  `gorm:"column:SYS_DISK_FILE_ID;type:int" json:"sysDiskFileId"`
}

// 指定表名（如果你想手动指定）
func (ChrShare) TableName() string {
	return "chr_share"
}
