package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"fullstack-backend/internal/database"
	"fullstack-backend/internal/models"

	"gorm.io/gorm/clause"
)

type emailMeta struct {
	Banned     bool   `json:"banned"`
	Price      int    `json:"price"`
	Sold       bool   `json:"sold"`
	NeedRepair bool   `json:"need_repair"`
	From       string `json:"from"`
}

type emailEntry struct {
	Main     string    `json:"main"`
	Password string    `json:"password"`
	Deputy   string    `json:"deputy"`
	Key2FA   string    `json:"key_2FA"`
	Meta     emailMeta `json:"meta"`
}

type emailPayload struct {
	Emails []emailEntry `json:"emails"`
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/fullstack?sslmode=disable"
	}

	db, err := database.Connect(dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect database: %v\n", err)
		os.Exit(1)
	}

	if err := database.Migrate(db); err != nil {
		fmt.Fprintf(os.Stderr, "failed to migrate: %v\n", err)
		os.Exit(1)
	}

	root, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get working dir: %v\n", err)
		os.Exit(1)
	}

	jsonPath := filepath.Join(root, "..", "..", "frontend", "src", "resource", "emails.json")
	payloadData, err := os.ReadFile(jsonPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read %s: %v\n", jsonPath, err)
		os.Exit(1)
	}

	var payload emailPayload
	if err := json.Unmarshal(payloadData, &payload); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse emails.json: %v\n", err)
		os.Exit(1)
	}

	if len(payload.Emails) == 0 {
		fmt.Println("no emails found in emails.json; nothing to seed")
		return
	}

	userIDRaw := os.Getenv("SEED_USER_ID")
	if userIDRaw == "" {
		fmt.Fprintln(os.Stderr, "SEED_USER_ID is required to seed emails")
		os.Exit(1)
	}
	userIDValue, err := strconv.ParseUint(userIDRaw, 10, 32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid SEED_USER_ID: %v\n", err)
		os.Exit(1)
	}
	userID := uint(userIDValue)

	importName := os.Getenv("SEED_IMPORT_NAME")
	if importName == "" {
		importName = "seed-emails"
	}

	tx := db.Begin()
	importRecord := models.EmailImport{
		UserID: userID,
		Name:   importName,
	}
	if err := tx.Create(&importRecord).Error; err != nil {
		tx.Rollback()
		fmt.Fprintf(os.Stderr, "failed to create import record: %v\n", err)
		os.Exit(1)
	}

	for _, entry := range payload.Emails {
		email := models.Email{
			UserID:     userID,
			ImportID:   importRecord.ID,
			Main:       entry.Main,
			Password:   entry.Password,
			Deputy:     entry.Deputy,
			Key2FA:     entry.Key2FA,
			Banned:     entry.Meta.Banned,
			Price:      entry.Meta.Price,
			Sold:       entry.Meta.Sold,
			NeedRepair: entry.Meta.NeedRepair,
			Source:     entry.Meta.From,
		}

		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}, {Name: "main"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"import_id",
				"password",
				"deputy",
				"key_2fa",
				"banned",
				"price",
				"sold",
				"need_repair",
				"source",
				"updated_at",
			}),
		}).Create(&email).Error; err != nil {
			tx.Rollback()
			fmt.Fprintf(os.Stderr, "failed to upsert %s: %v\n", entry.Main, err)
			os.Exit(1)
		}
	}

	if err := tx.Commit().Error; err != nil {
		fmt.Fprintf(os.Stderr, "failed to commit seed transaction: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("seeded %d emails\n", len(payload.Emails))
}
