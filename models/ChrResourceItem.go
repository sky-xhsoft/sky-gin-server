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
	SysCompanyId  *uint  `gorm:"column:SYS_COMPANY_ID" json:"sysCompanyId"`
	ChrResourceId *uint  `gorm:"column:CHR_RESOURCE_ID" json:"chrResourceId"`
	Name          string `gorm:"column:NAME" json:"name"`
	Type          string `gorm:"column:TYPE" json:"type"` // RTMP / VIDEO
	RtmpUrl       string `gorm:"column:RTMP_URL" json:"rtmpUrl"`
	CutTimes      *int   `gorm:"column:CUT_TIMES" json:"cutTimes"`
	VideoUrl      string `gorm:"column:VIDEO_URL" json:"videoUrl"`
	VideoParam    string `gorm:"column:VIDEO_PARAM" json:"videoParam"`
	IsActive      string `gorm:"column:IS_ACTIVE" json:"isActive"`
}

func (ChrResourceItem) TableName() string {
	return "chr_resource_item"
}
