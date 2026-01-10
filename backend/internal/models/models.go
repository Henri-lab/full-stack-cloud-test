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
