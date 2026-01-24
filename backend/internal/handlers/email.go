package handlers

import (
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

type EmailResponse struct {
	ID       uint      `json:"id"`
	Main     string    `json:"main"`
	Password string    `json:"password"`
	Deputy   string    `json:"deputy"`
	Key2FA   string    `json:"key_2FA"`
	Meta     EmailMeta `json:"meta"`
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
	}
}

func (h *EmailHandler) GetEmails(c *gin.Context) {
	var emails []models.Email
	if err := h.db.Order("id asc").Find(&emails).Error; err != nil {
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
