package middleware

import (
	"net/http"

	"fullstack-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LicenseKeyMiddleware 验证 License Key 的中间件
func LicenseKeyMiddleware(db *gorm.DB, requiredFeature string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			c.Abort()
			return
		}

		// 从请求头或请求体获取 Key
		keyCode := c.GetHeader("X-License-Key")
		if keyCode == "" {
			// 尝试从查询参数获取
			keyCode = c.Query("key")
		}

		if keyCode == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "需要有效的 License Key",
				"message": "请购买并激活 License Key 以使用此功能",
			})
			c.Abort()
			return
		}

		// 查询 Key
		var key models.LicenseKey
		if err := db.Where("key_code = ? AND user_id = ?", keyCode, userID).First(&key).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "无效的 License Key",
					"message": "Key 不存在或未激活",
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "验证 Key 失败"})
			}
			c.Abort()
			return
		}

		// 检查 Key 状态
		if key.Status == "revoked" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "License Key 已被撤销",
			})
			c.Abort()
			return
		}

		if key.Status == "exhausted" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "License Key 额度已用尽",
				"message": "请购买新的 Key 以继续使用",
			})
			c.Abort()
			return
		}

		// 检查额度
		if key.QuotaUsed >= key.QuotaTotal {
			// 更新状态为 exhausted
			key.Status = "exhausted"
			db.Save(&key)

			c.JSON(http.StatusForbidden, gin.H{
				"error": "License Key 额度已用尽",
				"message": "请购买新的 Key 以继续使用",
			})
			c.Abort()
			return
		}

		// 检查功能权限
		if !checkFeaturePermission(key.ProductType, requiredFeature) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "当前 License Key 不支持此功能",
				"message": "请升级到更高级别的 Key",
			})
			c.Abort()
			return
		}

		// 将 Key 信息存入上下文
		c.Set("license_key", key)
		c.Next()
	}
}

// ConsumeQuota 消耗额度的中间件（在请求成功后调用）
func ConsumeQuota(db *gorm.DB, amount int) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // 先执行业务逻辑

		// 只有在请求成功时才消耗额度
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			keyInterface, exists := c.Get("license_key")
			if !exists {
				return
			}

			key := keyInterface.(models.LicenseKey)

			// 更新已使用额度
			if err := db.Model(&models.LicenseKey{}).
				Where("id = ?", key.ID).
				UpdateColumn("quota_used", gorm.Expr("quota_used + ?", amount)).
				Error; err != nil {
				// 记录错误但不影响响应
				c.Set("quota_consume_error", err)
			}
		}
	}
}

// checkFeaturePermission 检查功能权限
func checkFeaturePermission(productType, feature string) bool {
	// 功能权限映射
	permissions := map[string][]string{
		"basic": {
			"email_verify",
		},
		"pro": {
			"email_verify",
			"email_import",
			"task_management",
		},
		"enterprise": {
			"email_verify",
			"email_import",
			"task_management",
			"api_access",
			"priority_support",
		},
	}

	allowedFeatures, exists := permissions[productType]
	if !exists {
		return false
	}

	for _, f := range allowedFeatures {
		if f == feature {
			return true
		}
	}

	return false
}

// OptionalLicenseKey 可选的 License Key 中间件（不强制要求）
func OptionalLicenseKey(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.Next()
			return
		}

		// 尝试获取 Key
		keyCode := c.GetHeader("X-License-Key")
		if keyCode == "" {
			keyCode = c.Query("key")
		}

		if keyCode != "" {
			var key models.LicenseKey
			if err := db.Where("key_code = ? AND user_id = ?", keyCode, userID).First(&key).Error; err == nil {
				if key.Status == "active" && key.QuotaUsed < key.QuotaTotal {
					c.Set("license_key", key)
					c.Set("has_license", true)
				}
			}
		}

		c.Next()
	}
}
