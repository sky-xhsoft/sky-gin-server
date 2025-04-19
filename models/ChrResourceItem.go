// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: ChrResourceItem.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/14
// Project Description:
// ----------------------------------------------------------------------------

package models

// 资源明细
type ChrResourceItem struct {
	Base
	ChrResourceId *uint   `gorm:"column:CHR_RESOURCE_ID" json:"chrResourceId"`
	ProjectId     uint    `gorm:"column:CHR_PROJECT_ID" json:"projectId"`
	Name          string  `gorm:"column:NAME" json:"name"`
	Type          string  `gorm:"column:TYPE" json:"type"` // RTMP / VIDEO
	RtmpUrl       string  `gorm:"column:RTMP_URL" json:"rtmpUrl"`
	CutTimes      *int    `gorm:"column:CUT_TIMES" json:"cutTimes"`
	VideoUrl      string  `gorm:"column:VIDEO_URL" json:"videoUrl"`
	VideoParam    string  `gorm:"column:VIDEO_PARAM" json:"videoParam"`
	VideoFileType string  `gorm:"column:VIDEO_FILE_TYPE" json:"videoFileType"`
	VideoFileSize float64 `gorm:"column:VIDEO_FILE_SIZE" json:"videoFileSize"`
	HeadImg       string  `gorm:"column:HEAD_IMG" json:"headImg"`
}

func (ChrResourceItem) TableName() string {
	return "chr_resource_item"
}
