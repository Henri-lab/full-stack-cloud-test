package handlers

import (
	"net/http"
	"time"

	"fullstack-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SubscriptionHandler struct {
	db *gorm.DB
}

func NewSubscriptionHandler(db *gorm.DB) *SubscriptionHandler {
	return &SubscriptionHandler{db: db}
}

type ActivateSubscriptionRequest struct {
	Plan         string `json:"plan"`
	DurationDays int    `json:"duration_days"`
}

func (h *SubscriptionHandler) GetMySubscription(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var sub models.Subscription
	if err := h.db.Where("user_id = ?", userID).Order("expires_at desc").First(&sub).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"subscription": nil})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询订阅失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subscription": sub})
}

func (h *SubscriptionHandler) ActivateSubscription(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var req ActivateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.DurationDays <= 0 {
		req.DurationDays = 30
	}
	if req.Plan == "" {
		req.Plan = "monthly"
	}

	now := time.Now()
	var sub models.Subscription
	if err := h.db.Where("user_id = ? AND status = ?", userID, "active").Order("expires_at desc").First(&sub).Error; err == nil {
		if sub.ExpiresAt.After(now) {
			sub.ExpiresAt = sub.ExpiresAt.Add(time.Duration(req.DurationDays) * 24 * time.Hour)
		} else {
			sub.ExpiresAt = now.Add(time.Duration(req.DurationDays) * 24 * time.Hour)
		}
		sub.Plan = req.Plan
		sub.Status = "active"
		if err := h.db.Save(&sub).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新订阅失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"subscription": sub})
		return
	}

	sub = models.Subscription{
		UserID:    userID.(uint),
		Plan:      req.Plan,
		ExpiresAt: now.Add(time.Duration(req.DurationDays) * 24 * time.Hour),
		Status:    "active",
	}
	if err := h.db.Create(&sub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建订阅失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subscription": sub})
}
