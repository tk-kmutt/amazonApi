package repository

import (
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
	IsDelete  bool
}
