package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"fullstack-backend/internal/database"
	"fullstack-backend/internal/models"
)

func escapeSQL(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}

func formatTimeLiteral(t time.Time) string {
	if t.IsZero() {
		return "NULL"
	}
	return fmt.Sprintf("'%s'", t.UTC().Format("2006-01-02 15:04:05"))
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

	var emails []models.Email
	if err := db.Order("id asc").Find(&emails).Error; err != nil {
		fmt.Fprintf(os.Stderr, "failed to query emails: %v\n", err)
		os.Exit(1)
	}

	if len(emails) == 0 {
		fmt.Println("-- no emails found")
		return
	}

	fmt.Println("INSERT INTO emails (main, password, deputy, key_2fa, banned, price, sold, need_repair, source, created_at, updated_at) VALUES")
	for i, email := range emails {
		line := fmt.Sprintf(
			"('%s','%s','%s','%s',%t,%d,%t,%t,'%s',%s,%s)",
			escapeSQL(email.Main),
			escapeSQL(email.Password),
			escapeSQL(email.Deputy),
			escapeSQL(email.Key2FA),
			email.Banned,
			email.Price,
			email.Sold,
			email.NeedRepair,
			escapeSQL(email.Source),
			formatTimeLiteral(email.CreatedAt),
			formatTimeLiteral(email.UpdatedAt),
		)
		if i < len(emails)-1 {
			line += ","
		} else {
			line += ";"
		}
		fmt.Println(line)
	}
}
