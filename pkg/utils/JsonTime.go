// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: JsonTime.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/15
// Project Description:
// ----------------------------------------------------------------------------

package utils

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

type JsonTime time.Time

func (jt JsonTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", time.Time(jt).Format(TimeFormat))
	return []byte(formatted), nil
}

func (jt *JsonTime) UnmarshalJSON(data []byte) error {
	t, err := time.Parse(`"`+TimeFormat+`"`, string(data))
	if err != nil {
		return err
	}
	*jt = JsonTime(t)
	return nil
}

func (jt JsonTime) String() string {
	return time.Time(jt).Format(TimeFormat)
}

// ✅ 关键：实现 driver.Valuer 接口
func (t JsonTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}

// ✅ 可选：实现 sql.Scanner 接口（用于从数据库读值）
func (t *JsonTime) Scan(value interface{}) error {
	if val, ok := value.(time.Time); ok {
		*t = JsonTime(val)
		return nil
	}
	return fmt.Errorf("cannot scan value %v into JsonTime", value)
}
