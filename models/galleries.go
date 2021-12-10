package models

import "github.com/jinzhu/gorm"

// image container resources that visitors view
type Galler struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}
