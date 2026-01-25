package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"fullstack-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EmailHandler struct {
	db *gorm.DB
}

func NewEmailHandler(db *gorm.DB) *EmailHandler {
	return &EmailHandler{db: db}
}

type EmailMeta struct {
	Banned     bool   `json:"banned"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	Price      int    `json:"price"`
	Sold       bool   `json:"sold"`
	NeedRepair bool   `json:"need_repair"`
	From       string `json:"from"`
}

type EmailFamilyResponse struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code"`
	Contact  string `json:"contact"`
	Issue    string `json:"issue"`
}

type EmailResponse struct {
	ID       uint                  `json:"id"`
	Main     string                `json:"main"`
	Password string                `json:"password"`
	Deputy   string                `json:"deputy"`
	Key2FA   string                `json:"key_2FA"`
	Meta     EmailMeta             `json:"meta"`
	Familys  []EmailFamilyResponse `json:"familys"`
}

type EmailMetaInput struct {
	Banned     *bool   `json:"banned"`
	Price      *int    `json:"price"`
	Sold       *bool   `json:"sold"`
	NeedRepair *bool   `json:"need_repair"`
	From       *string `json:"from"`
}

type CreateEmailRequest struct {
	Main     string          `json:"main" binding:"required,email"`
	Password string          `json:"password" binding:"required"`
	Deputy   string          `json:"deputy" binding:"omitempty,email"`
	Key2FA   string          `json:"key_2FA"`
	Meta     *EmailMetaInput `json:"meta"`
}

type UpdateEmailRequest struct {
	Main     string          `json:"main" binding:"omitempty,email"`
	Password string          `json:"password"`
	Deputy   string          `json:"deputy" binding:"omitempty,email"`
	Key2FA   string          `json:"key_2FA"`
	Meta     *EmailMetaInput `json:"meta"`
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

func emailToResponse(email models.Email) EmailResponse {
	familys := make([]EmailFamilyResponse, 0, len(email.Familys))
	for _, f := range email.Familys {
		familys = append(familys, EmailFamilyResponse{
			ID:       f.ID,
			Email:    f.Email,
			Password: f.Password,
			Code:     f.Code,
			Contact:  f.Contact,
			Issue:    f.Issue,
		})
	}
	return EmailResponse{
		ID:       email.ID,
		Main:     email.Main,
		Password: email.Password,
		Deputy:   email.Deputy,
		Key2FA:   email.Key2FA,
		Meta: EmailMeta{
			Banned:     email.Banned,
			CreatedAt:  formatTime(email.CreatedAt),
			UpdatedAt:  formatTime(email.UpdatedAt),
			Price:      email.Price,
			Sold:       email.Sold,
			NeedRepair: email.NeedRepair,
			From:       email.Source,
		},
		Familys: familys,
	}
}

func (h *EmailHandler) GetEmails(c *gin.Context) {
	var emails []models.Email
	if err := h.db.Preload("Familys").Order("id asc").Find(&emails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch emails"})
		return
	}

	responses := make([]EmailResponse, 0, len(emails))
	for _, email := range emails {
		responses = append(responses, emailToResponse(email))
	}

	c.JSON(http.StatusOK, responses)
}

func (h *EmailHandler) GetEmail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email ID"})
		return
	}

	var email models.Email
	if err := h.db.First(&email, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Email not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch email"})
		return
	}

	c.JSON(http.StatusOK, emailToResponse(email))
}

func (h *EmailHandler) CreateEmail(c *gin.Context) {
	var req CreateEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email := models.Email{
		Main:     req.Main,
		Password: req.Password,
		Deputy:   req.Deputy,
		Key2FA:   req.Key2FA,
	}

	if req.Meta != nil {
		if req.Meta.Banned != nil {
			email.Banned = *req.Meta.Banned
		}
		if req.Meta.Price != nil {
			email.Price = *req.Meta.Price
		}
		if req.Meta.Sold != nil {
			email.Sold = *req.Meta.Sold
		}
		if req.Meta.NeedRepair != nil {
			email.NeedRepair = *req.Meta.NeedRepair
		}
		if req.Meta.From != nil {
			email.Source = *req.Meta.From
		}
	}

	if err := h.db.Create(&email).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create email"})
		return
	}

	c.JSON(http.StatusCreated, emailToResponse(email))
}

func (h *EmailHandler) UpdateEmail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email ID"})
		return
	}

	var req UpdateEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var email models.Email
	if err := h.db.First(&email, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Email not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch email"})
		return
	}

	if req.Main != "" {
		email.Main = req.Main
	}
	if req.Password != "" {
		email.Password = req.Password
	}
	if req.Deputy != "" {
		email.Deputy = req.Deputy
	}
	if req.Key2FA != "" {
		email.Key2FA = req.Key2FA
	}
	if req.Meta != nil {
		if req.Meta.Banned != nil {
			email.Banned = *req.Meta.Banned
		}
		if req.Meta.Price != nil {
			email.Price = *req.Meta.Price
		}
		if req.Meta.Sold != nil {
			email.Sold = *req.Meta.Sold
		}
		if req.Meta.NeedRepair != nil {
			email.NeedRepair = *req.Meta.NeedRepair
		}
		if req.Meta.From != nil {
			email.Source = *req.Meta.From
		}
	}

	if err := h.db.Save(&email).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update email"})
		return
	}

	c.JSON(http.StatusOK, emailToResponse(email))
}

func (h *EmailHandler) DeleteEmail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email ID"})
		return
	}

	var email models.Email
	if err := h.db.First(&email, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Email not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch email"})
		return
	}

	if err := h.db.Delete(&email).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email deleted successfully"})
}

// Import request types
type ImportFamilyInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code"`
	Contact  string `json:"contact"`
	Issue    string `json:"issue"`
}

type ImportEmailInput struct {
	Main     string              `json:"main"`
	Password string              `json:"password"`
	Deputy   string              `json:"deputy"`
	Key2FA   string              `json:"key_2FA"`
	Meta     *EmailMetaInput     `json:"meta"`
	Familys  []ImportFamilyInput `json:"familys"`
}

type ImportRequest struct {
	Emails []ImportEmailInput `json:"emails"`
}

func (h *EmailHandler) ImportEmails(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer f.Close()

	var req ImportRequest
	if err := json.NewDecoder(f).Decode(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format: " + err.Error()})
		return
	}

	if len(req.Emails) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No emails to import"})
		return
	}

	// Check for duplicates in database
	for _, emailInput := range req.Emails {
		var existing models.Email
		if err := h.db.Where("main = ?", emailInput.Main).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": fmt.Sprintf("Email already exists: %s", emailInput.Main),
			})
			return
		}
	}

	// Use transaction to insert all emails
	tx := h.db.Begin()
	imported := 0

	for _, emailInput := range req.Emails {
		email := models.Email{
			Main:     emailInput.Main,
			Password: emailInput.Password,
			Deputy:   emailInput.Deputy,
			Key2FA:   emailInput.Key2FA,
		}

		if emailInput.Meta != nil {
			if emailInput.Meta.Banned != nil {
				email.Banned = *emailInput.Meta.Banned
			}
			if emailInput.Meta.Price != nil {
				email.Price = *emailInput.Meta.Price
			}
			if emailInput.Meta.Sold != nil {
				email.Sold = *emailInput.Meta.Sold
			}
			if emailInput.Meta.NeedRepair != nil {
				email.NeedRepair = *emailInput.Meta.NeedRepair
			}
			if emailInput.Meta.From != nil {
				email.Source = *emailInput.Meta.From
			}
		}

		if err := tx.Create(&email).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("Failed to import email %s: %s", emailInput.Main, err.Error()),
			})
			return
		}

		// Insert family emails
		for _, familyInput := range emailInput.Familys {
			family := models.EmailFamily{
				EmailID:  email.ID,
				Email:    familyInput.Email,
				Password: familyInput.Password,
				Code:     familyInput.Code,
				Contact:  familyInput.Contact,
				Issue:    familyInput.Issue,
			}
			if err := tx.Create(&family).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("Failed to import family email for %s: %s", emailInput.Main, err.Error()),
				})
				return
			}
		}

		imported++
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Import successful",
		"imported": imported,
	})
}
