package repository

import "gorm.io/plugin/soft_delete"

// 論理削除
const (
	ACTIVE   soft_delete.DeletedAt = 0
	INACTIVE                       = 1
)
