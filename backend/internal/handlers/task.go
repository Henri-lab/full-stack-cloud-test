package handlers

import (
	"net/http"
	"strconv"

	"fullstack-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskHandler struct {
	db *gorm.DB
}

func NewTaskHandler(db *gorm.DB) *TaskHandler {
	return &TaskHandler{db: db}
}

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=200"`
	Description string `json:"description"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	var tasks []models.Task

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get all tasks for the user
	if err := h.db.Where("creator_id = ?", userID).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userID, _ := c.Get("user_id")

	var task models.Task
	if err := h.db.Where("id = ? AND creator_id = ?", id, userID).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	task := models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      "open",
		CreatorID:   userID.(uint),
	}

	if err := h.db.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")

	var task models.Task
	if err := h.db.Where("id = ? AND creator_id = ?", id, userID).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task"})
		return
	}

	// Update fields
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Status != "" {
		task.Status = req.Status
	}

	if err := h.db.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userID, _ := c.Get("user_id")

	var task models.Task
	if err := h.db.Where("id = ? AND creator_id = ?", id, userID).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task"})
		return
	}

	if err := h.db.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
