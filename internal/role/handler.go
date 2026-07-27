package role

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
	Slug          string `json:"slug" binding:"required,min=2,max=50"`
	Name          string `json:"name" binding:"required,min=2,max=100"`
	Description   string `json:"description"`
}

type UpdateRequest struct {
	Slug        string `json:"slug" binding:"required,min=2,max=50"`
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

type AssignPermissionRequest struct {
	PermissionID uint `json:"permission_id" binding:"required"`
}

// Create creates a new role for an application
// @Summary Create a new role
// @Description Create a new role associated with an application
// @Tags Admin - Roles
// @Accept json
// @Produce json
// @Security Bearer
// @Param role body CreateRequest true "Role data"
// @Success 201 {object} Role
// @Failure 400 {object} apierror.Body
// @Failure 409 {object} apierror.Body
// @Router /api/admin/roles [post]
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

	r, err := h.Service.Create(req.ApplicationID, req.Slug, req.Name, req.Description)
	if err != nil {
		if errors.Is(err, ErrSlugTaken) {
			apierror.Respond(c, apierror.Conflict("slug_taken", err.Error()))
			return
		}
		apierror.Respond(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "role created successfully", "role": r})
}

// List returns all roles, optionally filtered by application
// @Summary Get all roles
// @Description Retrieve a list of all roles, optionally filtered by application_id
// @Tags Admin - Roles
// @Accept json
// @Produce json
// @Security Bearer
// @Param application_id query int false "Filter by application ID"
// @Success 200 {array} Role
// @Router /api/admin/roles [get]
func (h *Handler) List(c *gin.Context) {
	applicationIDStr := c.Query("application_id")

	var roles []Role
	var err error

	if applicationIDStr != "" {
		applicationID, parseErr := strconv.ParseUint(applicationIDStr, 10, 64)
		if parseErr != nil {
			apierror.Respond(c, apierror.BadRequest("invalid_application_id", "invalid application_id parameter"))
			return
		}
		roles, err = h.Service.ListByApplication(uint(applicationID))
	} else {
		roles, err = h.Service.List()
	}

	if err != nil {
		apierror.Respond(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"roles": roles})
}

// Get returns a single role by ID
// @Summary Get role by ID
// @Description Retrieve a specific role by its ID
// @Tags Admin - Roles
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Role ID"
// @Success 200 {object} Role
// @Failure 400 {object} apierror.Body
// @Failure 404 {object} apierror.Body
// @Router /api/admin/roles/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_id", "invalid role ID"))
		return
	}

	r, err := h.Service.GetByID(uint(id))
	if err != nil {
		apierror.Respond(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"role": r})
}

// Update updates a role
// @Summary Update a role
// @Description Update an existing role's fields
// @Tags Admin - Roles
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Role ID"
// @Param role body UpdateRequest true "Role data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} apierror.Body
// @Router /api/admin/roles/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_id", "invalid role ID"))
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

	c.JSON(http.StatusOK, gin.H{"message": "role updated successfully"})
}

// Delete removes a role
// @Summary Delete a role
// @Description Delete a role by ID
// @Tags Admin - Roles
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Role ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} apierror.Body
// @Router /api/admin/roles/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_id", "invalid role ID"))
		return
	}

	if err := h.Service.Delete(uint(id)); err != nil {
		apierror.Respond(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "role deleted successfully"})
}

// AssignPermission grants a permission to a role
// @Summary Assign permission to role
// @Description Assign a permission to a specific role (both must belong to the same application)
// @Tags Admin - Role-Permission Relationships
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Role ID"
// @Param assignment body AssignPermissionRequest true "Permission assignment data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} apierror.Body
// @Failure 409 {object} apierror.Body
// @Router /api/admin/roles/{id}/permissions [post]
func (h *Handler) AssignPermission(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_id", "invalid role ID"))
		return
	}

	var req AssignPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_request", err.Error()))
		return
	}

	if err := h.Service.AssignPermission(uint(roleID), req.PermissionID); err != nil {
		switch {
		case errors.Is(err, ErrAlreadyAssigned):
			apierror.Respond(c, apierror.Conflict("already_assigned", err.Error()))
		case errors.Is(err, ErrApplicationMismatch):
			apierror.Respond(c, apierror.BadRequest("application_mismatch", err.Error()))
		default:
			apierror.Respond(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "permission assigned to role successfully"})
}

// RemovePermission revokes a permission from a role
// @Summary Remove permission from role
// @Description Remove a permission from a specific role
// @Tags Admin - Role-Permission Relationships
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Role ID"
// @Param permissionId path int true "Permission ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} apierror.Body
// @Router /api/admin/roles/{id}/permissions/{permissionId} [delete]
func (h *Handler) RemovePermission(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_id", "invalid role ID"))
		return
	}

	permissionID, err := strconv.ParseUint(c.Param("permissionId"), 10, 64)
	if err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_permission_id", "invalid permission ID"))
		return
	}

	if err := h.Service.RemovePermission(uint(roleID), uint(permissionID)); err != nil {
		apierror.Respond(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "permission removed from role successfully"})
}

// ListPermissions lists the permissions granted to a role
// @Summary List role permissions
// @Description Retrieve all permissions assigned to a specific role
// @Tags Admin - Role-Permission Relationships
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Role ID"
// @Success 200 {array} go_auth_internal_permission.Permission
// @Failure 400 {object} apierror.Body
// @Router /api/admin/roles/{id}/permissions [get]
func (h *Handler) ListPermissions(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_id", "invalid role ID"))
		return
	}

	perms, err := h.Service.ListPermissions(uint(roleID))
	if err != nil {
		apierror.Respond(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"permissions": perms})
}
