package handlers

import (
	"net/http"
	"regexp"
	"time"

	"go_auth/auth"
	"go_auth/db"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	db   *db.Service
	auth *auth.Service
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

type UserInfo struct {
	ID           int                        `json:"id"`
	UUID         string                     `json:"uuid"`
	Username     string                     `json:"username"`
	Email        string                     `json:"email"`
	Applications []db.UserApplicationAccess `json:"applications"`
}

func New(dbService *db.Service, authService *auth.Service) *Handler {
	return &Handler{
		db:   dbService,
		auth: authService,
	}
}

// Register creates a new user account
// @Summary Register a new user
// @Description Create a new user account with username, email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User registration data"
// @Success 201 {object} map[string]interface{} "User created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - validation error or user already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate email format
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Check if email already exists
	emailExists, err := h.db.EmailExists(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if emailExists {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Check if username already exists
	usernameExists, err := h.db.UsernameExists(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if usernameExists {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
		return
	}

	// Hash password
	hashedPassword, err := h.auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// Create user
	user, err := h.db.CreateUser(req.Username, req.Email, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"id":       user.ID,
			"uuid":     user.UUID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// Login authenticates a user and returns a JWT token
// @Summary User login
// @Description Authenticate user with email and password, returns JWT token with user info
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "User login credentials"
// @Success 200 {object} LoginResponse "Login successful"
// @Failure 400 {object} map[string]interface{} "Bad request - validation error"
// @Failure 401 {object} map[string]interface{} "Unauthorized - invalid credentials"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user by email
	user, err := h.db.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password
	if !h.auth.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Get user roles and permissions organized by application
	userAccess, err := h.db.GetUserAccessByApplications(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user access"})
		return
	}

	// Generate JWT token with application structure
	token, err := h.auth.GenerateToken(user.ID, user.Email, userAccess)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	response := LoginResponse{
		Token: token,
		User: UserInfo{
			ID:           user.ID,
			UUID:         user.UUID,
			Username:     user.Username,
			Email:        user.Email,
			Applications: userAccess,
		},
	}

	c.JSON(http.StatusOK, response)
}

// Me returns current user information
// @Summary Get current user
// @Description Get current authenticated user information including roles and permissions
// @Tags Authentication
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} UserInfo "Current user information"
// @Failure 401 {object} map[string]interface{} "Unauthorized - invalid or missing token"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /me [get]
func (h *Handler) Me(c *gin.Context) {
	// Get user from context (set by auth middleware)
	userClaims, exists := auth.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	// Get fresh user data from database
	user, err := h.db.GetUserByID(userClaims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Get user access organized by applications
	userAccess, err := h.db.GetUserAccessByApplications(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user access"})
		return
	}

	userInfo := UserInfo{
		ID:           user.ID,
		UUID:         user.UUID,
		Username:     user.Username,
		Email:        user.Email,
		Applications: userAccess,
	}

	c.JSON(http.StatusOK, userInfo)
}

// Health returns the health status of the API
// @Summary Health check
// @Description Returns the health status of the API and database connection
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "API is healthy"
// @Failure 500 {object} map[string]interface{} "API is unhealthy"
// @Router /health [get]
func (h *Handler) Health(c *gin.Context) {
	// Check database connection
	if err := h.db.Ping(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":    "unhealthy",
			"message":   "Database connection failed",
			"timestamp": time.Now().Format(time.RFC3339),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"message":   "API is running normally",
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0.0",
	})
}
