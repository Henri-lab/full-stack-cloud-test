package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"fullstack-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PaymentHandler struct {
	db *gorm.DB
}

func NewPaymentHandler(db *gorm.DB) *PaymentHandler {
	return &PaymentHandler{db: db}
}

// ProductConfig 产品配置
type ProductConfig struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Price       int    `json:"price"`        // 价格（分）
	QuotaAmount int    `json:"quota_amount"` // 次数额度
	Features    []string `json:"features"`
}

var ProductConfigs = []ProductConfig{
	{
		Type:        "basic",
		Name:        "基础版",
		Price:       1000, // 10元
		QuotaAmount: 100,
		Features:    []string{"邮箱验证 100 次", "基础功能"},
	},
	{
		Type:        "pro",
		Name:        "专业版",
		Price:       3000, // 30元
		QuotaAmount: 500,
		Features:    []string{"邮箱验证 500 次", "所有功能", "优先支持"},
	},
	{
		Type:        "enterprise",
		Name:        "企业版",
		Price:       5000, // 50元
		QuotaAmount: 1000,
		Features:    []string{"邮箱验证 1000 次", "所有功能", "专属客服", "API 访问"},
	},
}

// GetProducts 获取产品列表
func (h *PaymentHandler) GetProducts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"products": ProductConfigs,
	})
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	ProductType string `json:"product_type" binding:"required"`
}

// CreateOrder 创建支付订单
func (h *PaymentHandler) CreateOrder(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 查找产品配置
	var productConfig *ProductConfig
	for _, p := range ProductConfigs {
		if p.Type == req.ProductType {
			productConfig = &p
			break
		}
	}
	if productConfig == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的产品类型"})
		return
	}

	// 生成订单号
	orderNo := generateOrderNo()

	// 创建订单
	payment := models.Payment{
		UserID:      userID.(uint),
		OrderNo:     orderNo,
		Amount:      productConfig.Price,
		ProductType: productConfig.Type,
		QuotaAmount: productConfig.QuotaAmount,
		Status:      "pending",
		ExpiredAt:   time.Now().Add(15 * time.Minute), // 15分钟过期
	}

	if err := h.db.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建订单失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order": payment,
		"product": productConfig,
	})
}

// GetOrder 获取订单详情
func (h *PaymentHandler) GetOrder(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	orderNo := c.Param("order_no")

	var payment models.Payment
	if err := h.db.Where("order_no = ? AND user_id = ?", orderNo, userID).First(&payment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询订单失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": payment})
}

// PaymentNotifyRequest 支付回调请求（模拟）
type PaymentNotifyRequest struct {
	OrderNo       string `json:"order_no" binding:"required"`
	TransactionID string `json:"transaction_id" binding:"required"`
	PaymentMethod string `json:"payment_method" binding:"required"` // alipay, wechat
}

// PaymentNotify 支付回调（实际应该是支付平台回调）
func (h *PaymentHandler) PaymentNotify(c *gin.Context) {
	var req PaymentNotifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 开启事务
	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 查找订单
	var payment models.Payment
	if err := tx.Where("order_no = ?", req.OrderNo).First(&payment).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询订单失败"})
		}
		return
	}

	// 检查订单状态
	if payment.Status != "pending" {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "订单状态异常"})
		return
	}

	// 检查订单是否过期
	if time.Now().After(payment.ExpiredAt) {
		payment.Status = "expired"
		tx.Save(&payment)
		tx.Commit()
		c.JSON(http.StatusBadRequest, gin.H{"error": "订单已过期"})
		return
	}

	// 更新订单状态
	now := time.Now()
	payment.Status = "paid"
	payment.PaymentMethod = req.PaymentMethod
	payment.TransactionID = req.TransactionID
	payment.PaidAt = &now

	if err := tx.Save(&payment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新订单失败"})
		return
	}

	// 生成 License Key
	keyCode := generateLicenseKey()
	licenseKey := models.LicenseKey{
		UserID:      payment.UserID,
		PaymentID:   payment.ID,
		KeyCode:     keyCode,
		ProductType: payment.ProductType,
		QuotaTotal:  payment.QuotaAmount,
		QuotaUsed:   0,
		Status:      "active",
		ActivatedAt: &now,
	}

	if err := tx.Create(&licenseKey).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成密钥失败"})
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "提交事务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "支付成功",
		"license_key": licenseKey,
	})
}

// GetMyKeys 获取我的密钥列表
func (h *PaymentHandler) GetMyKeys(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var keys []models.LicenseKey
	if err := h.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&keys).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询密钥失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"keys": keys})
}

// ActivateKeyRequest 激活密钥请求
type ActivateKeyRequest struct {
	KeyCode string `json:"key_code" binding:"required"`
}

// ActivateKey 激活密钥
func (h *PaymentHandler) ActivateKey(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var req ActivateKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	var key models.LicenseKey
	if err := h.db.Where("key_code = ?", req.KeyCode).First(&key).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "密钥不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询密钥失败"})
		}
		return
	}

	// 检查密钥是否已被其他用户激活
	if key.ActivatedAt != nil && key.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "密钥已被其他用户激活"})
		return
	}

	// 检查密钥状态
	if key.Status == "revoked" {
		c.JSON(http.StatusForbidden, gin.H{"error": "密钥已被撤销"})
		return
	}

	if key.Status == "exhausted" {
		c.JSON(http.StatusForbidden, gin.H{"error": "密钥额度已用尽"})
		return
	}

	// 如果是首次激活，更新用户ID
	if key.ActivatedAt == nil {
		now := time.Now()
		key.UserID = userID.(uint)
		key.ActivatedAt = &now
		if err := h.db.Save(&key).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "激活密钥失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "激活成功",
		"key": key,
	})
}

// CheckKeyRequest 检查密钥请求
type CheckKeyRequest struct {
	KeyCode string `json:"key_code" binding:"required"`
}

// CheckKey 检查密钥状态
func (h *PaymentHandler) CheckKey(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	var req CheckKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	var key models.LicenseKey
	if err := h.db.Where("key_code = ? AND user_id = ?", req.KeyCode, userID).First(&key).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "密钥不存在或未激活"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询密钥失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"key": key,
		"quota_remaining": key.QuotaTotal - key.QuotaUsed,
	})
}

// 生成订单号
func generateOrderNo() string {
	timestamp := time.Now().Format("20060102150405")
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	randomStr := hex.EncodeToString(randomBytes)
	return fmt.Sprintf("ORD%s%s", timestamp, randomStr)
}

// 生成 License Key
func generateLicenseKey() string {
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	keyStr := hex.EncodeToString(randomBytes)
	// 格式化为 XXXX-XXXX-XXXX-XXXX
	return fmt.Sprintf("%s-%s-%s-%s",
		keyStr[0:4],
		keyStr[4:8],
		keyStr[8:12],
		keyStr[12:16],
	)
}
