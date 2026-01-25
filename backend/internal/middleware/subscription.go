package middleware

import (
	"net/http"
	"time"

	"fullstack-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SubscriptionMiddleware enforces active subscriptions for protected routes.
func SubscriptionMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			c.Abort()
			return
		}

		var sub models.Subscription
		if err := db.Where("user_id = ? AND status = ? AND expires_at > ?", userID, "active", time.Now()).
			Order("expires_at desc").
			First(&sub).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusForbidden, gin.H{"error": "订阅已过期或不存在"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "订阅校验失败"})
			}
			c.Abort()
			return
		}

		c.Set("subscription", sub)
		c.Next()
	}
}
