package models

import "time"

type Base struct {
	ID           uint      `gorm:"primaryKey;column:ID;type:int unsigned"`                   // 主键自增
	SysCompanyId uint      `gorm:"column:SYS_COMPANY_ID;type:int unsigned"`                  // 所属公司（非空）
	CreateBy     string    `gorm:"column:CREATE_BY;type:varchar(80)"`                        // 创建人（非空）
	CreateTime   time.Time `gorm:"column:CREATE_TIME;type:datetime;not null;autoCreateTime"` // 自动记录时间
	UpdateBy     string    `gorm:"column:UPDATE_BY;type:varchar(80)"`                        // 更新人（非空）
	UpdateTime   time.Time `gorm:"column:UPDATE_TIME;type:datetime;not null;autoUpdateTime"` // 自动更新时间
	IsActive     string    `gorm:"column:IS_ACTIVE;type:char(1);default:Y;not null"`         // 是否有效
}
