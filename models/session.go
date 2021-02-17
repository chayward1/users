package models

import "gorm.io/gorm"

type Session struct {
	gorm.Model
	Token   string
	Expires int64
	UserID  uint
}
