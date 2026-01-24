package handlers

import (
	"net/http"
	"os"
	"sync"
	"time"

	"fullstack-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Rate limiter for login attempts
type RateLimiter struct {
	attempts map[string]*LoginAttempt
	mu       sync.RWMutex
}

type LoginAttempt struct {
	Count     int
	FirstTry  time.Time
	BlockedAt time.Time
}

var rateLimiter = &RateLimiter{
	attempts: make(map[string]*LoginAttempt),
}

const (
	maxLoginAttempts = 5
	blockDuration    = 15 * time.Minute
	attemptWindow    = 10 * time.Minute
)

func (rl *RateLimiter) IsBlocked(ip string) bool {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	attempt, exists := rl.attempts[ip]
	if !exists {
		return false
	}

	// Check if block has expired
	if !attempt.BlockedAt.IsZero() && time.Since(attempt.BlockedAt) < blockDuration {
		return true
	}

	return false
}

func (rl *RateLimiter) RecordAttempt(ip string, success bool) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if success {
		delete(rl.attempts, ip)
		return
	}

	attempt, exists := rl.attempts[ip]
	if !exists || time.Since(attempt.FirstTry) > attemptWindow {
		rl.attempts[ip] = &LoginAttempt{
			Count:    1,
			FirstTry: time.Now(),
		}
		return
	}

	attempt.Count++
	if attempt.Count >= maxLoginAttempts {
		attempt.BlockedAt = time.Now()
	}
}

func (rl *RateLimiter) GetRemainingAttempts(ip string) int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	attempt, exists := rl.attempts[ip]
	if !exists {
		return maxLoginAttempts
	}

	return maxLoginAttempts - attempt.Count
}

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// validatePassword checks password strength
func validatePassword(password string) (bool, string) {
	if len(password) < 12 {
		return false, "Password must be at least 12 characters"
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false
	specialChars := "!@#$%^&*()_+-=[]{}|;':\",./<>?"

	for _, c := range password {
		switch {
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= '0' && c <= '9':
			hasNumber = true
		default:
			for _, sc := range specialChars {
				if c == sc {
					hasSpecial = true
					break
				}
			}
		}
	}

	if !hasUpper {
		return false, "Password must contain at least one uppercase letter"
	}
	if !hasLower {
		return false, "Password must contain at least one lowercase letter"
	}
	if !hasNumber {
		return false, "Password must contain at least one number"
	}
	if !hasSpecial {
		return false, "Password must contain at least one special character (!@#$%^&*()_+-=[]{}|;':\",./<>?)"
	}

	return true, ""
}

func (h *AuthHandler) Register(c *gin.Context) {
	clientIP := c.ClientIP()

	// Rate limit registration attempts to prevent abuse
	if rateLimiter.IsBlocked(clientIP) {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Too many requests. Please try again later.",
		})
		return
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate password strength
	if valid, msg := validatePassword(req.Password); !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := h.db.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hash password with higher cost for production
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	clientIP := c.ClientIP()

	// Check if IP is blocked due to too many failed attempts
	if rateLimiter.IsBlocked(clientIP) {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Too many failed login attempts. Please try again later.",
		})
		return
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		rateLimiter.RecordAttempt(clientIP, false)
		remaining := rateLimiter.GetRemainingAttempts(clientIP)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":              "Invalid credentials",
			"remaining_attempts": remaining,
		})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		rateLimiter.RecordAttempt(clientIP, false)
		remaining := rateLimiter.GetRemainingAttempts(clientIP)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":              "Invalid credentials",
			"remaining_attempts": remaining,
		})
		return
	}

	// Successful login - clear failed attempts
	rateLimiter.RecordAttempt(clientIP, true)

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	})

	jwtSecret := getJWTSecret()
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable is required")
	}
	return secret
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// In a stateless JWT system, logout is handled client-side
	// by removing the token from storage
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
