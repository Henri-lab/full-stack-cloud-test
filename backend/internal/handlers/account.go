package handlers

import (
	"net/http"
	"strconv"
	"time"

	"fullstack-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccountHandler struct {
	db *gorm.DB
}

func NewAccountHandler(db *gorm.DB) *AccountHandler {
	return &AccountHandler{db: db}
}

type AccountResponse struct {
	ID     uint   `json:"id"`
	Type   string `json:"type"`
	Main   string `json:"main"`
	Status string `json:"status"`
	Source string `json:"source"`
}

type FamilyInfo struct {
	Capacity int `json:"capacity"`
	Used     int `json:"used"`
}

type AccountListItem struct {
	AccountResponse
	Family *FamilyInfo `json:"family,omitempty"`
}

type AccountCredentialsResponse struct {
	ID       uint   `json:"id"`
	Main     string `json:"main"`
	Password string `json:"password"`
	Key2FA   string `json:"key_2FA"`
}

type AccountClaimRequest struct {
	AccountID uint `json:"account_id" binding:"required"`
}

type AccountBindRequest struct {
	AccountID      uint   `json:"account_id" binding:"required"`
	MemberEmail    string `json:"member_email" binding:"required,email"`
	MemberPassword string `json:"member_password" binding:"required"`
}

type AccountPurchaseRequest struct {
	AccountID uint  `json:"account_id" binding:"required"`
	PaymentID *uint `json:"payment_id"`
}

func (h *AccountHandler) ListAccounts(c *gin.Context) {
	var accounts []models.Account
	query := h.db.Model(&models.Account{}).Where("status <> ?", "retired").Order("id asc")
	if accountType := c.Query("type"); accountType != "" {
		query = query.Where("type = ?", accountType)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&accounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch accounts"})
		return
	}

	results := make([]AccountListItem, 0, len(accounts))
	for _, account := range accounts {
		item := AccountListItem{
			AccountResponse: AccountResponse{
				ID:     account.ID,
				Type:   account.Type,
				Main:   account.Main,
				Status: account.Status,
				Source: account.Source,
			},
		}
		if account.Type == "family" {
			var group models.FamilyGroup
			if err := h.db.Where("account_id = ?", account.ID).First(&group).Error; err == nil {
				var count int64
				h.db.Model(&models.FamilyBinding{}).Where("family_group_id = ?", group.ID).Count(&count)
				item.Family = &FamilyInfo{Capacity: group.Capacity, Used: int(count)}
			}
		}
		results = append(results, item)
	}

	c.JSON(http.StatusOK, results)
}

func (h *AccountHandler) ClaimTemporary(c *gin.Context) {
	var req AccountClaimRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var account models.Account
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND type = ?", req.AccountID, "temporary").
		First(&account).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to lock account"})
		}
		return
	}

	if account.Status != "available" {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{"error": "Account is not available"})
		return
	}

	now := time.Now()
	account.Status = "locked"
	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account"})
		return
	}

	usage := models.TemporaryUsage{
		AccountID: account.ID,
		UserID:    userID.(uint),
		StartedAt: now,
		ExpiresAt: now.Add(24 * time.Hour),
	}
	if err := tx.Create(&usage).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create usage"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account claimed", "expires_at": usage.ExpiresAt})
}

func (h *AccountHandler) ReleaseTemporary(c *gin.Context) {
	var req AccountClaimRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var usage models.TemporaryUsage
	if err := tx.Where("account_id = ? AND user_id = ? AND returned_at IS NULL", req.AccountID, userID).
		Order("started_at desc").
		First(&usage).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "No active usage found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch usage"})
		}
		return
	}

	now := time.Now()
	usage.ReturnedAt = &now
	if err := tx.Save(&usage).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update usage"})
		return
	}

	if err := tx.Model(&models.Account{}).
		Where("id = ? AND type = ?", req.AccountID, "temporary").
		Update("status", "available").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to release account"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account released"})
}

func (h *AccountHandler) PurchaseExclusive(c *gin.Context) {
	var req AccountPurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var account models.Account
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND type = ?", req.AccountID, "exclusive").
		First(&account).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to lock account"})
		}
		return
	}

	if account.Status != "available" {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{"error": "Account is not available"})
		return
	}

	account.Status = "sold"
	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account"})
		return
	}

	purchase := models.ExclusivePurchase{
		AccountID:   account.ID,
		UserID:      userID.(uint),
		PaymentID:   req.PaymentID,
		PurchasedAt: time.Now(),
	}
	if err := tx.Create(&purchase).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create purchase"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Purchase successful"})
}

func (h *AccountHandler) BindFamily(c *gin.Context) {
	var req AccountBindRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var account models.Account
	if err := h.db.Where("id = ? AND type = ?", req.AccountID, "family").First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch account"})
		}
		return
	}

	var group models.FamilyGroup
	if err := h.db.Where("account_id = ?", account.ID).First(&group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			group = models.FamilyGroup{AccountID: account.ID, Capacity: 5}
			if err := h.db.Create(&group).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create group"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch group"})
			return
		}
	}

	var count int64
	h.db.Model(&models.FamilyBinding{}).Where("family_group_id = ?", group.ID).Count(&count)
	if int(count) >= group.Capacity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Family group is full"})
		return
	}

	var existing models.FamilyBinding
	if err := h.db.Where("family_group_id = ? AND user_id = ?", group.ID, userID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Already bound to this family group"})
		return
	}

	binding := models.FamilyBinding{
		FamilyGroupID:     group.ID,
		UserID:            userID.(uint),
		MemberEmail:       req.MemberEmail,
		MemberPasswordEnc: req.MemberPassword,
	}
	if err := h.db.Create(&binding).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to bind family account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Family binding created"})
}

func (h *AccountHandler) UnbindFamily(c *gin.Context) {
	var req AccountClaimRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var group models.FamilyGroup
	if err := h.db.Where("account_id = ?", req.AccountID).First(&group).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Family group not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch family group"})
		}
		return
	}

	if err := h.db.Where("family_group_id = ? AND user_id = ?", group.ID, userID).Delete(&models.FamilyBinding{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unbind family account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Family binding removed"})
}

func (h *AccountHandler) GetExclusiveCredentials(c *gin.Context) {
	accountID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var purchase models.ExclusivePurchase
	if err := h.db.Where("account_id = ? AND user_id = ?", accountID, userID).First(&purchase).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusForbidden, gin.H{"error": "No access to credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify purchase"})
		}
		return
	}

	var account models.Account
	if err := h.db.Where("id = ? AND type = ?", accountID, "exclusive").First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch account"})
		}
		return
	}

	c.JSON(http.StatusOK, AccountCredentialsResponse{
		ID:       account.ID,
		Main:     account.Main,
		Password: account.Password,
		Key2FA:   account.Key2FA,
	})
}
