// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: SysUser.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/11
// Project Description:
// ----------------------------------------------------------------------------

package models

type SysUser struct {
	Base

	TrueName string `gorm:"column:TRUE_NAME;type:varchar(255)" json:"trueName"` // 真实姓名
	Username string `gorm:"column:USERNAME;type:varchar(255)" json:"username"`  // 用户名称
	Password string `gorm:"column:PASSWORD;type:varchar(255)" json:"-"`         // 密码（不返回前端）
	Phone    string `gorm:"column:PHONE;type:varchar(20)" json:"phone"`         // 手机号
	Email    string `gorm:"column:EMAIL;type:varchar(255)" json:"email"`        // 邮箱
	Language string `gorm:"column:LANGUAGE;type:varchar(255)" json:"language"`  // 语言
}

// 指定表名（可选，如果未启用 GORM 的 `SingularTable` 策略）
func (SysUser) TableName() string {
	return "sys_user"
}
