package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID         uint64     `gorm:"primaryKey"`
	Username   string     `gorm:"size:64"`
	Password   string     `gorm:"size:255"`
	Notes      []Note     `gorm:"foreignKey:UserID"`
	CreditCard CreditCard `gorm:"foreignKey:UserID"`
	Role       Role       `json:"role"`
}

type Note struct {
	gorm.Model
	ID      uint64 `gorm:"primaryKey"`
	Name    string `gorm:"size:255"`
	Content string `gorm:"type:text"`
	UserID  uint64 `gorm:"index"`
}

type CreditCard struct {
	gorm.Model
	Number string
	UserID uint64 `gorm:"primaryKey"`
}
