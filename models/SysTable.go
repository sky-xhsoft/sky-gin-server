// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: SysTable.go
// Author: xhsoftware-skyzhou
// Created On: 2025/5/15
// Project Description:
// ----------------------------------------------------------------------------

package models

type SysTable struct {
	Base
	Name          string `gorm:"column:NAME;type:varchar(255);not null" json:"NAME"`                 // 表名
	DisplayName   string `gorm:"column:DISPLAY_NAME;type:varchar(255);not null" json:"DISPLAY_NAME"` // 显示名称
	RealTableID   *uint  `gorm:"column:REAL_TABLE_ID" json:"REAL_TABLE_ID"`                          // 实际表ID
	Filter        string `gorm:"column:FILTER;type:varchar(2000)" json:"FILTER"`                     // 过滤SQL
	Mask          string `gorm:"column:MASK;type:char(10)" json:"MASK"`                              // 表单规则
	IsMenu        string `gorm:"column:IS_MENU;type:char(1);default:'N'" json:"IS_MENU"`             // 是否菜单
	DkColumnID    *uint  `gorm:"column:DK_COLUMN_ID" json:"DK_COLUMN_ID"`                            // 显示字段
	ParentTableID *uint  `gorm:"column:PARENT_TABLE_ID" json:"PARENT_TABLE_ID"`                      // 父表ID
	Props         string `gorm:"column:PROPS;type:varchar(2000)" json:"PROPS"`                       // 扩展属性
	Description   string `gorm:"column:DESCRIPTION;type:varchar(2000)" json:"DESCRIPTION"`           // 备注
	DisplayCol    string `gorm:"column:DISPLAY_COL;type:char(1)" json:"DISPLAY_COL"`                 // 显示列数
}

func (SysTable) TableName() string {
	return "sys_table"
}
