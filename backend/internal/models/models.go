package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Username     string         `gorm:"uniqueIndex;not null" json:"username"`
	Email        string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"not null" json:"-"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type Task struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	Status      string         `gorm:"default:'open'" json:"status"`
	CreatorID   uint           `gorm:"not null" json:"creator_id"`
	Creator     User           `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Email struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	Main       string         `gorm:"uniqueIndex;not null" json:"main"`
	Password   string         `gorm:"not null" json:"password"`
	Deputy     string         `json:"deputy"`
	Key2FA     string         `gorm:"column:key_2fa" json:"key_2FA"`
	Banned     bool           `gorm:"default:false" json:"-"`
	Price      int            `gorm:"default:0" json:"-"`
	Sold       bool           `gorm:"default:false" json:"-"`
	NeedRepair bool           `gorm:"default:false" json:"-"`
	Source     string         `gorm:"column:source" json:"-"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
