// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: project.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/14
// Project Description:
// ----------------------------------------------------------------------------

package models

type ChrProject struct {
	Base
	Name        string         `gorm:"column:NAME" json:"name"`
	Type        string         `gorm:"column:TYPE" json:"type"`          // LL/ZY
	IsScreen    string         `gorm:"column:IS_SCREEN" json:"isScreen"` // Y/N
	Prem        string         `gorm:"column:PREM" json:"prem"`          // 默认权限
	Size        float64        `gorm:"column:SIZE" json:"size"`
	Qty         int            `gorm:"column:QTY" json:"qty"`
	HeadImg     string         `gorm:"column:HEAD_IMG" json:"headImg"`
	Resouse     []ChrResource  `gorm:"-" json:"resouse"`
	ProjectUser ChrProjectUser `gorm:"-" json:"projectUser"`
	Files       []SysDiskFile  `gorm:"-" json:"files"`
}

func (ChrProject) TableName() string {
	return "chr_project"
}

type ChrProjectUser struct {
	ProjectId uint   `gorm:"column:CHR_PROJECT_ID" json:"projectId"`
	UserId    uint   `gorm:"column:SYS_USER_ID" json:"userId"`
	Prem      string `gorm:"column:PREM" json:"prem"` // R/D/E/A
	IsOwner   string `gorm:"column:IS_OWNER" json:"isOwner"`
	Base
}

func (ChrProjectUser) TableName() string {
	return "chr_project_user"
}
