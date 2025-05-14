// ----------------------------------------------------------------------------
// Project Name: sky-gin-server
// File Name: SysColumn.go
// Author: xhsoftware-skyzhou
// Created On: 2025/5/15
// Project Description:
// ----------------------------------------------------------------------------

package models

type SysColumn struct {
	Base
	SysTableId    uint   `gorm:"column:SYS_TABLE_ID" json:"SYS_TABLE_ID"`
	Name          string `gorm:"column:NAME" json:"NAME"`
	DisplayName   string `gorm:"column:DISPLAY_NAME" json:"DISPLAY_NAME"`
	FullName      string `gorm:"column:FULL_NAME" json:"FULL_NAME"`
	OrderNo       int    `gorm:"column:ORDERNO" json:"ORDERNO"`
	SubMethod     string `gorm:"column:SUBMETHOD" json:"SUBMETHOD"`
	ColType       string `gorm:"column:COL_TYPE" json:"COL_TYPE"`
	ColLength     int    `gorm:"column:COL_LENGTH" json:"COL_LENGTH"`
	ColPrecision  int    `gorm:"column:COL_PRECISION" json:"COL_PRECISION"`
	Mask          string `gorm:"column:MASK" json:"MASK"`
	NullAble      string `gorm:"column:NULL_ABLE" json:"NULL_ABLE"`
	IsDK          string `gorm:"column:IS_DK" json:"IS_DK"`
	IsQuery       string `gorm:"column:IS_QUERY" json:"IS_QUERY"`
	IsUppercase   string `gorm:"column:IS_UPPERCASE" json:"IS_UPPERCASE"`
	SetValueType  string `gorm:"column:SET_VALUE_TYPE" json:"SET_VALUE_TYPE"`
	RefTableId    uint   `gorm:"column:REF_TABLE_ID" json:"REF_TABLE_ID"`
	RefColumnId   uint   `gorm:"column:REF_COLUMN_ID" json:"REF_COLUMN_ID"`
	ColumnSQL     string `gorm:"column:COLUMN_SQL" json:"COLUMN_SQL"`
	RefOnDelete   string `gorm:"column:REF_ON_DELETE" json:"REF_ON_DELETE"`
	SEQ           string `gorm:"column:SEQ" json:"SEQ"`
	SysDictId     int    `gorm:"column:SYS_DICT_ID" json:"SYS_DICT_ID"`
	DefaultValue  string `gorm:"column:DEFAULT_VALUE" json:"DEFAULT_VALUE"`
	RegExpression string `gorm:"column:REG_EXPRESSION" json:"REG_EXPRESSION"`
	ErrMsg        string `gorm:"column:ERR_MSG" json:"ERR_MSG"`
	DisplayType   string `gorm:"column:DISPLAY_TYPE" json:"DISPLAY_TYPE"`
	DisplayCols   int    `gorm:"column:DISPLAY_COLS" json:"DISPLAY_COLS"`
	DisplayRows   int    `gorm:"column:DISPLAY_ROWS" json:"DISPLAY_ROWS"`
	Props         string `gorm:"column:PROPS" json:"PROPS"`
	IsShowTitle   string `gorm:"column:IS_SHOW_TITLE" json:"IS_SHOW_TITLE"`
	Description   string `gorm:"column:DESCRIPTION" json:"DESCRIPTION"`
	ShowColumnId  int    `gorm:"column:SHOW_COLUMN_ID" json:"SHOW_COLUMN_ID"`
	ShowColumnVal string `gorm:"column:SHOW_COLUMN_VAL" json:"SHOW_COLUMN_VAL"`
}

func (SysColumn) TableName() string {
	return "sys_column"
}
