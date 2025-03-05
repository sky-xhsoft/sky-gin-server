package models

type SysRoutes struct {
	Base
	Path   string `gorm:"column:PATH;type:varchar(255);not null;index:idx_path_method"`  // 路径（非空）
	Method string `gorm:"column:METHOD;type:varchar(10);not null;index:idx_path_method"` // 方法（非空）
	Handle string `gorm:"column:HANDLE;type:varchar(255);not null"`                      // 处理函数（非空）
}
