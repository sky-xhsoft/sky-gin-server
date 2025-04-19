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
	Name            string            `gorm:"column:NAME" json:"name"`
	ProjectId       uint              `gorm:"column:CHR_PROJECT_ID" json:"projectId"`
	CutStatus       int               `gorm:"column:CUT_STATUS" json:"cutStatus"`
	RecordingStatus int               `gorm:"column:RECORDING_STATUS" json:"recordingStatus"`
	Size            float64           `gorm:"column:SIZE" json:"size"`
	Qty             int               `gorm:"column:QTY" json:"qty"`
	HeadImg         string            `gorm:"column:HEAD_IMG" json:"headImg"`
	Items           []ChrResourceItem `gorm:"-" json:"items"`
}

func (ChrResource) TableName() string {
	return "chr_resource"
}
