package repository

import (
	"gorm.io/plugin/soft_delete"
	"time"
)

type AmazonItems struct {
	Name      string
	Maker     string
	Price     int
	Comment   string
	Url       string
	Asin      string
	CreatedAt time.Time
	UpdatedAt time.Time
	ACTIVE    soft_delete.DeletedAt `gorm:"softDelete:flag"`
}
