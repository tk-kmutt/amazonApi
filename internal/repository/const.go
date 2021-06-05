package repository

import "gorm.io/plugin/soft_delete"

// 論理削除
const (
	NOT_DELETE soft_delete.DeletedAt = 0
	DELETE                           = 1
)
