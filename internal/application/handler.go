package application

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go_auth/internal/apierror"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

type CreateRequest struct {
	Slug        string `json:"slug" binding:"required,min=2,max=50"`
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
}

type UpdateRequest struct {
	Slug        string `json:"slug" binding:"required,min=2,max=50"`
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

// Create creates a new application
// @Summary Create a new application
// @Description Create a new application in the system
// @Tags Admin - Applications
// @Accept json
// @Produce json
// @Security Bearer
// @Param application body CreateRequest true "Application data"
// @Success 201 {object} Application
// @Failure 400 {object} apierror.Body
// @Failure 409 {object} apierror.Body
// @Failure 500 {object} apierror.Body
// @Router /api/admin/applications [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_request", err.Error()))
		return
	}

	app, err := h.Service.Create(req.Slug, req.Name, req.Description)
	if err != nil {
		if errors.Is(err, ErrSlugTaken) {
			apierror.Respond(c, apierror.Conflict("slug_taken", err.Error()))
			return
		}
		apierror.Respond(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "application created successfully", "application": app})
}

// List returns all applications
// @Summary Get all applications
// @Description Retrieve a list of all applications in the system
// @Tags Admin - Applications
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} Application
// @Failure 500 {object} apierror.Body
// @Router /api/admin/applications [get]
func (h *Handler) List(c *gin.Context) {
	apps, err := h.Service.List()
	if err != nil {
		apierror.Respond(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"applications": apps})
}

// Get returns a single application by ID
// @Summary Get application by ID
// @Description Retrieve a specific application by its ID
// @Tags Admin - Applications
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Application ID"
// @Success 200 {object} Application
// @Failure 400 {object} apierror.Body
// @Failure 404 {object} apierror.Body
// @Router /api/admin/applications/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_id", "invalid application ID"))
		return
	}

	app, err := h.Service.GetByID(uint(id))
	if err != nil {
		apierror.Respond(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"application": app})
}

// Update updates an application
// @Summary Update an application
// @Description Update an existing application's fields
// @Tags Admin - Applications
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Application ID"
// @Param application body UpdateRequest true "Application data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} apierror.Body
// @Failure 500 {object} apierror.Body
// @Router /api/admin/applications/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_id", "invalid application ID"))
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_request", err.Error()))
		return
	}

	if err := h.Service.Update(uint(id), req.Slug, req.Name, req.Description, req.Active); err != nil {
		apierror.Respond(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "application updated successfully"})
}

// Delete removes an application
// @Summary Delete an application
// @Description Delete an application by ID
// @Tags Admin - Applications
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Application ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} apierror.Body
// @Failure 500 {object} apierror.Body
// @Router /api/admin/applications/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_id", "invalid application ID"))
		return
	}

	if err := h.Service.Delete(uint(id)); err != nil {
		apierror.Respond(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "application deleted successfully"})
}
