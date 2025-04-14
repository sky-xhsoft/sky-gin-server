// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: ChrResource.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/14
// Project Description:
// ----------------------------------------------------------------------------

package models

// 资源组结构体
type ChrResource struct {
	Base
	Name      string `gorm:"column:NAME" json:"name"`
	ProjectId uint   `gorm:"column:CHR_PROJECT_ID" json:"projectId"`
}

func (ChrResource) TableName() string {
	return "chr_resource"
}
