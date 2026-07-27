package permission

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go_auth/internal/apierror"
	"go_auth/internal/application"
)

type Handler struct {
	Service      *Service
	Applications *application.Service
}

func NewHandler(service *Service, applications *application.Service) *Handler {
	return &Handler{Service: service, Applications: applications}
}

type CreateRequest struct {
	ApplicationID uint   `json:"application_id" binding:"required"`
	Name          string `json:"name" binding:"required,min=2,max=100"`
	Slug          string `json:"slug" binding:"required,min=2,max=50"`
	Description   string `json:"description"`
}

type UpdateRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Slug        string `json:"slug" binding:"required,min=2,max=50"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

// Create creates a new permission for an application
// @Summary Create a new permission
// @Description Create a new permission associated with an application
// @Tags Admin - Permissions
// @Accept json
// @Produce json
// @Security Bearer
// @Param permission body CreateRequest true "Permission data"
// @Success 201 {object} Permission
// @Failure 400 {object} apierror.Body
// @Failure 409 {object} apierror.Body
// @Router /api/admin/permissions [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_request", err.Error()))
		return
	}

	if _, err := h.Applications.GetByID(req.ApplicationID); err != nil {
		apierror.Respond(c, apierror.BadRequest("application_not_found", "application not found"))
		return
	}

	p, err := h.Service.Create(req.ApplicationID, req.Name, req.Slug, req.Description)
	if err != nil {
		if errors.Is(err, ErrSlugTaken) {
			apierror.Respond(c, apierror.Conflict("slug_taken", err.Error()))
			return
		}
		apierror.Respond(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "permission created successfully", "permission": p})
}

// List returns all permissions, optionally filtered by application
// @Summary Get all permissions
// @Description Retrieve a list of all permissions, optionally filtered by application_id
// @Tags Admin - Permissions
// @Accept json
// @Produce json
// @Security Bearer
// @Param application_id query int false "Filter by application ID"
// @Success 200 {array} Permission
// @Failure 400 {object} apierror.Body
// @Router /api/admin/permissions [get]
func (h *Handler) List(c *gin.Context) {
	applicationIDStr := c.Query("application_id")

	var perms []Permission
	var err error

	if applicationIDStr != "" {
		applicationID, parseErr := strconv.ParseUint(applicationIDStr, 10, 64)
		if parseErr != nil {
			apierror.Respond(c, apierror.BadRequest("invalid_application_id", "invalid application_id parameter"))
			return
		}
		perms, err = h.Service.ListByApplication(uint(applicationID))
	} else {
		perms, err = h.Service.List()
	}

	if err != nil {
		apierror.Respond(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"permissions": perms})
}

// Get returns a single permission by ID
// @Summary Get permission by ID
// @Description Retrieve a specific permission by its ID
// @Tags Admin - Permissions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Permission ID"
// @Success 200 {object} Permission
// @Failure 400 {object} apierror.Body
// @Failure 404 {object} apierror.Body
// @Router /api/admin/permissions/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_id", "invalid permission ID"))
		return
	}

	p, err := h.Service.GetByID(uint(id))
	if err != nil {
		apierror.Respond(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"permission": p})
}

// Update updates a permission
// @Summary Update a permission
// @Description Update an existing permission's fields
// @Tags Admin - Permissions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Permission ID"
// @Param permission body UpdateRequest true "Permission data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} apierror.Body
// @Router /api/admin/permissions/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_id", "invalid permission ID"))
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_request", err.Error()))
		return
	}

	if err := h.Service.Update(uint(id), req.Name, req.Slug, req.Description, req.Active); err != nil {
		apierror.Respond(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "permission updated successfully"})
}

// Delete removes a permission
// @Summary Delete a permission
// @Description Delete a permission by ID
// @Tags Admin - Permissions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Permission ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} apierror.Body
// @Router /api/admin/permissions/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_id", "invalid permission ID"))
		return
	}

	if err := h.Service.Delete(uint(id)); err != nil {
		apierror.Respond(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "permission deleted successfully"})
}
