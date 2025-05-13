// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: SysApi.go
// Author: xhsoftware-skyzhou
// Created On: 2025/5/14
// Project Description:
// ----------------------------------------------------------------------------

package models

type SysApi struct {
	Base
	Name       string `gorm:"column:NAME;type:varchar(255);not null"`                       // 名称（非空）
	Path       string `gorm:"column:PATH;type:varchar(255);not null;index:idx_path_method"` // 路径（非空）
	Method     string `gorm:"column:METHOD;type:varchar(10);not null"`                      // 方法（非空）
	Handle     string `gorm:"column:HANDLE;type:varchar(255);not null"`                     // 处理函数（非空）
	Permission string `gorm:"column:PERMISSION;type:varchar(2);not null"`
	ReqDemo    string `gorm:"column:REQDEMO;type:text" json:"ReqDemo"`
	ResDemo    string `gorm:"column:RESDEMO;type:text" json:"ResDemo"`
	ReqFields  string `gorm:"column:REQFIELDS;type:text" json:"reqFields"`
	ResFields  string `gorm:"column:RESFIELDS;type:text" json:"resFields"`
}

func (SysApi) TableName() string {
	return "sys_api"
}
