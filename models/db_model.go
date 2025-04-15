package models

import (
	"github.com/sky-xhsoft/sky-gin-server/pkg/utils"
)

type Base struct {
	ID           uint           `gorm:"primaryKey;column:ID;type:int unsigned" json:"ID"`
	SysCompanyId uint           `gorm:"column:SYS_COMPANY_ID;type:int unsigned" json:"sysCompanyId"`
	CreateBy     string         `gorm:"column:CREATE_BY;type:varchar(80)" json:"createBy"`
	CreateTime   utils.JsonTime `gorm:"column:CREATE_TIME;type:datetime;not null;autoCreateTime" json:"createTime"`
	UpdateBy     string         `gorm:"column:UPDATE_BY;type:varchar(80)" json:"updateBy"`
	UpdateTime   utils.JsonTime `gorm:"column:UPDATE_TIME;type:datetime;not null;autoUpdateTime" json:"updateTime"`
	IsActive     string         `gorm:"column:IS_ACTIVE;type:char(1);default:Y;not null" json:"isActive"`
}
