package user

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"

	"go_auth/internal/apierror"
	"go_auth/internal/auth"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

type Handler struct {
	Service *Service
	Auth    *auth.Service
}

func NewHandler(service *Service, authService *auth.Service) *Handler {
	return &Handler{Service: service, Auth: authService}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Register creates a new user account
// @Summary Register a new user
// @Description Create a new user account with username, email and password
// @Tags Users
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User registration data"
// @Success 201 {object} map[string]interface{} "User created successfully"
// @Failure 400 {object} apierror.Body "Bad request - validation error"
// @Failure 409 {object} apierror.Body "Email or username already taken"
// @Failure 500 {object} apierror.Body "Internal server error"
// @Router /api/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_request", err.Error()))
		return
	}

	if !emailRegex.MatchString(req.Email) {
		apierror.Respond(c, apierror.BadRequest("invalid_email", "invalid email format"))
		return
	}

	hashedPassword, err := h.Auth.HashPassword(req.Password)
	if err != nil {
		apierror.Respond(c, err)
		return
	}

	u, err := h.Service.Create(req.Username, req.Email, hashedPassword)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmailTaken):
			apierror.Respond(c, apierror.Conflict("email_taken", "email already registered"))
		case errors.Is(err, ErrUsernameTaken):
			apierror.Respond(c, apierror.Conflict("username_taken", "username already taken"))
		default:
			apierror.Respond(c, err)
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
		"user": gin.H{
			"id":       u.ID,
			"uuid":     u.UUID,
			"username": u.Username,
			"email":    u.Email,
		},
	})
}

// List returns all registered users
// @Summary List all users
// @Description Retrieve a list of all users registered in the system
// @Tags Admin - Users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} User
// @Failure 500 {object} apierror.Body
// @Router /api/admin/users [get]
func (h *Handler) List(c *gin.Context) {
	users, err := h.Service.List()
	if err != nil {
		apierror.Respond(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// Get returns a single user by ID
// @Summary Get user by ID
// @Description Retrieve a specific user by their ID
// @Tags Admin - Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Failure 400 {object} apierror.Body
// @Failure 404 {object} apierror.Body
// @Router /api/admin/users/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_id", "invalid user ID"))
		return
	}

	u, err := h.Service.GetByID(uint(id))
	if err != nil {
		apierror.Respond(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": u})
}
