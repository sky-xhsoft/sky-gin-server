// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: hash_test.go
// Author: xhsoftware-skyzhou
// Created On: 2025/4/11
// Project Description:
// ----------------------------------------------------------------------------

package hash

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	password, _ := HashPassword("abc123")
	fmt.Println(password)
}
