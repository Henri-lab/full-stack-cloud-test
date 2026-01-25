package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Username     string         `gorm:"uniqueIndex;not null" json:"username"`
	Email        string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"not null" json:"-"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type Task struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	Status      string         `gorm:"default:'open'" json:"status"`
	CreatorID   uint           `gorm:"not null" json:"creator_id"`
	Creator     User           `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Email struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	UserID     uint           `gorm:"not null;index:idx_email_user_main,unique" json:"-"`
	ImportID   uint           `gorm:"index" json:"-"`
	Main       string         `gorm:"not null;index:idx_email_user_main,unique" json:"main"`
	Password   string         `gorm:"not null" json:"password"`
	Deputy     string         `json:"deputy"`
	Key2FA     string         `gorm:"column:key_2fa" json:"key_2FA"`
	Status     string         `gorm:"default:'unknown'" json:"-"` // unknown, live, verify, dead
	Banned     bool           `gorm:"default:false" json:"-"`
	Price      int            `gorm:"default:0" json:"-"`
	Sold       bool           `gorm:"default:false" json:"-"`
	NeedRepair bool           `gorm:"default:false" json:"-"`
	Source     string         `gorm:"column:source" json:"-"`
	Import     EmailImport    `gorm:"foreignKey:ImportID" json:"-"`
	Familys    []EmailFamily  `gorm:"foreignKey:EmailID" json:"familys,omitempty"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type EmailImport struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	Name       string    `gorm:"not null" json:"name"`
	SourceFile string    `json:"-"`
	CreatedAt  time.Time `json:"created_at"`
}

type EmailFamily struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	EmailID   uint           `gorm:"not null;index" json:"email_id"`
	Email     string         `gorm:"not null" json:"email"`
	Password  string         `gorm:"not null" json:"password"`
	Code      string         `json:"code"`
	Contact   string         `json:"contact"`
	Issue     string         `json:"issue"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Payment 支付订单
type Payment struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	UserID        uint           `gorm:"not null;index" json:"user_id"`
	OrderNo       string         `gorm:"uniqueIndex;not null" json:"order_no"`
	Amount        int            `gorm:"not null" json:"amount"`                                    // 金额（分）
	ProductType   string         `gorm:"not null" json:"product_type"`                             // basic, pro, enterprise
	QuotaAmount   int            `gorm:"not null" json:"quota_amount"`                             // 购买的次数额度
	Status        string         `gorm:"default:'pending'" json:"status"`                          // pending, paid, expired, refunded
	PaymentMethod string         `json:"payment_method"`                                           // alipay, wechat
	TransactionID string         `gorm:"index" json:"transaction_id"`                              // 第三方支付流水号
	PaidAt        *time.Time     `json:"paid_at"`
	ExpiredAt     time.Time      `json:"expired_at"`                                               // 订单过期时间（15分钟）
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// LicenseKey 授权密钥
type LicenseKey struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	PaymentID   uint           `gorm:"not null;index" json:"payment_id"`
	KeyCode     string         `gorm:"uniqueIndex;not null" json:"key_code"`
	ProductType string         `gorm:"not null" json:"product_type"`                  // basic, pro, enterprise
	QuotaTotal  int            `gorm:"not null" json:"quota_total"`                   // 总次数
	QuotaUsed   int            `gorm:"default:0" json:"quota_used"`                   // 已使用次数
	Status      string         `gorm:"default:'active'" json:"status"`                // active, exhausted, revoked
	ActivatedAt *time.Time     `json:"activated_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
