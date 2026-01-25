package database

import (
	"fullstack-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Task{},
		&models.Email{},
		&models.EmailImport{},
		&models.EmailFamily{},
		&models.Account{},
		&models.TemporaryUsage{},
		&models.ExclusivePurchase{},
		&models.FamilyGroup{},
		&models.FamilyBinding{},
		&models.Subscription{},
		&models.AuditLog{},
		&models.Payment{},
		&models.LicenseKey{},
	)
}
