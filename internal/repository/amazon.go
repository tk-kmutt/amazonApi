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
	Asin      string `gorm:"unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
	IsDelete  soft_delete.DeletedAt `gorm:"softDelete:flag"`
}
