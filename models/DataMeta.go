// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: DataMeta.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/14
// Project Description:
// ----------------------------------------------------------------------------

package models

import (
	"github.com/gin-gonic/gin"
	"reflect"
	"time"
)

func FillCreateMeta(c *gin.Context, model any) {
	fillMeta(c, model, true)
}

func FillUpdateMeta(c *gin.Context, model any) {
	fillMeta(c, model, false)
}

func fillMeta(c *gin.Context, model any, isCreate bool) {
	user := getUser(c)
	//now := time.Now().Format("2006-01-02 15:04:05")

	v := reflect.ValueOf(model).Elem()
	base := v.FieldByName("Base")
	if !base.IsValid() {
		return
	}

	if isCreate {
		base.FieldByName("CreateBy").SetString(user.Username)

		//设置创建时间
		val := reflect.ValueOf(&base).Elem()
		field := val.FieldByName("CreateTime")
		if field.IsValid() && field.CanSet() && field.Type().AssignableTo(reflect.TypeOf(time.Time{})) {
			field.Set(reflect.ValueOf(time.Now()))
		}

		base.FieldByName("IsActive").SetString("Y")
		base.FieldByName("SysCompanyId").SetUint(uint64(user.SysCompanyId))
	}
	base.FieldByName("UpdateBy").SetString(user.Username)

	//设置更新时间
	val := reflect.ValueOf(&base).Elem()
	field := val.FieldByName("UpdateTime")
	if field.IsValid() && field.CanSet() && field.Type().AssignableTo(reflect.TypeOf(time.Time{})) {
		field.Set(reflect.ValueOf(time.Now()))
	}

}

func getUser(c *gin.Context) *SysUser {
	user, exists := c.Get("User")
	if !exists {
		return &SysUser{Username: "system"}
	}
	if u, ok := user.(*SysUser); ok {
		return u
	}
	return &SysUser{Username: "system"}
}

func FillUpdateMetaMap(c *gin.Context, data map[string]interface{}) {
	userVal, exists := c.Get("User")
	if exists {
		if u, ok := userVal.(*SysUser); ok {
			data["UPDATE_BY"] = u.Username
		}
	}
	data["UPDATE_TIME"] = time.Now()
}
