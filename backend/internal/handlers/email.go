package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
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
	Status   string                `json:"status"` // unknown, live, verify, dead
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
		Status:   email.Status,
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
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDValue.(uint)

	var emails []models.Email
	query := h.db.Preload("Familys").Order("id asc").Where("user_id = ?", userID)
	if importIDRaw := c.Query("import_id"); importIDRaw != "" {
		importID, err := strconv.ParseUint(importIDRaw, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid import ID"})
			return
		}
		query = query.Where("import_id = ?", importID)
	}

	if err := query.Find(&emails).Error; err != nil {
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

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDValue.(uint)

	var email models.Email
	if err := h.db.Preload("Familys").Where("id = ? AND user_id = ?", id, userID).First(&email).Error; err != nil {
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

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDValue.(uint)

	email := models.Email{
		UserID:   userID,
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

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDValue.(uint)

	var email models.Email
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&email).Error; err != nil {
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

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDValue.(uint)

	var email models.Email
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&email).Error; err != nil {
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

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDValue.(uint)

	// Check for duplicates in database
	for _, emailInput := range req.Emails {
		var existing models.Email
		if err := h.db.Where("user_id = ? AND main = ?", userID, emailInput.Main).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": fmt.Sprintf("Email already exists: %s", emailInput.Main),
			})
			return
		}
	}

	// Use transaction to insert all emails
	tx := h.db.Begin()
	imported := 0

	importName := strings.TrimSpace(file.Filename)
	if importName == "" {
		importName = fmt.Sprintf("import-%s", time.Now().Format("20060102-150405"))
	}

	importRecord := models.EmailImport{
		UserID:     userID,
		Name:       importName,
		SourceFile: file.Filename,
	}
	if err := tx.Create(&importRecord).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create import record"})
		return
	}

	for _, emailInput := range req.Emails {
		email := models.Email{
			UserID:   userID,
			ImportID: importRecord.ID,
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
		"message":     "Import successful",
		"imported":    imported,
		"import_id":   importRecord.ID,
		"import_name": importRecord.Name,
	})
}

type EmailImportSummary struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	Count     int    `json:"count"`
}

func (h *EmailHandler) GetEmailImports(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDValue.(uint)

	var imports []models.EmailImport
	if err := h.db.Where("user_id = ?", userID).Order("created_at desc").Find(&imports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch imports"})
		return
	}

	type importCount struct {
		ImportID uint
		Count    int
	}
	var counts []importCount
	if err := h.db.Model(&models.Email{}).
		Select("import_id, COUNT(*) as count").
		Where("user_id = ? AND import_id <> 0", userID).
		Group("import_id").
		Scan(&counts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch import counts"})
		return
	}

	countMap := make(map[uint]int, len(counts))
	for _, c := range counts {
		countMap[c.ImportID] = c.Count
	}

	result := make([]EmailImportSummary, 0, len(imports))
	for _, item := range imports {
		result = append(result, EmailImportSummary{
			ID:        item.ID,
			Name:      item.Name,
			CreatedAt: formatTime(item.CreatedAt),
			Count:     countMap[item.ID],
		})
	}

	c.JSON(http.StatusOK, result)
}

// Verify request types
type VerifyEmailRequest struct {
	Emails []string `json:"mail" binding:"required"`
	Key    string   `json:"key"`                       // 第三方 API 需要
	Method string   `json:"method" binding:"required"` // "api" 或 "smtp"
}

type VerifyEmailResponse struct {
	Email  string `json:"email"`
	Status string `json:"status"` // live, verify, dead, error
	Error  string `json:"error,omitempty"`
}

func (h *EmailHandler) VerifyEmails(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.Emails) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No emails to verify"})
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDValue.(uint)

	var results []VerifyEmailResponse

	// 根据验证方法选择不同的验证逻辑
	switch req.Method {
	case "smtp":
		// 使用 SMTP 验证
		results = h.verifyEmailsSMTP(req.Emails)
	case "api":
		// 使用第三方 API 验证
		if req.Key == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Key is required for API method"})
			return
		}
		var err error
		results, err = h.verifyEmailsAPI(req.Emails, req.Key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid method. Use 'smtp' or 'api'"})
		return
	}

	// 更新数据库中的邮箱状态
	for _, result := range results {
		var dbEmail models.Email
		if err := h.db.Where("user_id = ? AND main = ?", userID, result.Email).First(&dbEmail).Error; err == nil {
			dbEmail.Status = result.Status
			h.db.Save(&dbEmail)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"results": results,
		"total":   len(results),
		"method":  req.Method,
	})
}

// verifyEmailsAPI 使用第三方 API 验证邮箱
func (h *EmailHandler) verifyEmailsAPI(emails []string, key string) ([]VerifyEmailResponse, error) {
	apiURL := "https://gmailver.com/php/check1.php"
	payload := map[string]interface{}{
		"mail":      emails,
		"key":       key,
		"fastCheck": false,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare request")
	}

	// 发送请求到第三方 API
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to verify emails: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response")
	}

	// 解析响应 - 第三方 API 返回格式：{"data": {"email": "status", ...}}
	var apiResponse struct {
		Message      string                 `json:"message"`
		Data         map[string]interface{} `json:"data"`
		ResponseTime string                 `json:"responseTime"`
		Status       string                 `json:"status"`
	}

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse response: %s", string(body))
	}

	// 检查 API 是否返回错误
	if apiResponse.Status == "error" || apiResponse.Data == nil {
		return nil, fmt.Errorf("third-party API returned error: %s", apiResponse.Message)
	}

	// 从 data 字段中提取邮箱状态
	results := make([]VerifyEmailResponse, 0)
	for email, statusInterface := range apiResponse.Data {
		status, ok := statusInterface.(string)
		if !ok {
			status = "error"
		}

		results = append(results, VerifyEmailResponse{
			Email:  email,
			Status: status,
		})
	}

	return results, nil
}

// verifyEmailsSMTP 使用 SMTP 验证邮箱
func (h *EmailHandler) verifyEmailsSMTP(emails []string) []VerifyEmailResponse {
	verifier := NewSMTPVerifier()
	results := make([]VerifyEmailResponse, 0)

	for _, email := range emails {
		status, err := verifier.VerifyEmail(email)
		result := VerifyEmailResponse{
			Email:  email,
			Status: status,
		}
		if err != nil {
			result.Error = err.Error()
		}
		results = append(results, result)

		// 添加延迟避免被限流
		time.Sleep(500 * time.Millisecond)
	}

	return results
}
