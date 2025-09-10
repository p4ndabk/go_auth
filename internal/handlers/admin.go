package handlers

import (
	"net/http"
	"strconv"

	db "go_auth/internal/database"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	db *db.Service
}

// Request structs
type CreateApplicationRequest struct {
	Slug        string `json:"slug" binding:"required,min=2,max=50"`
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
}

type UpdateApplicationRequest struct {
	Slug        string `json:"slug" binding:"required,min=2,max=50"`
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

type CreateRoleRequest struct {
	ApplicationID int    `json:"application_id" binding:"required"`
	Slug          string `json:"slug" binding:"required,min=2,max=50"`
	Name          string `json:"name" binding:"required,min=2,max=100"`
	Description   string `json:"description"`
}

type UpdateRoleRequest struct {
	Slug        string `json:"slug" binding:"required,min=2,max=50"`
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

type CreatePermissionRequest struct {
	ApplicationID int    `json:"application_id" binding:"required"`
	Name          string `json:"name" binding:"required,min=2,max=100"`
	Slug          string `json:"slug" binding:"required,min=2,max=50"`
	Description   string `json:"description"`
}

type UpdatePermissionRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Slug        string `json:"slug" binding:"required,min=2,max=50"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

type AssignRolePermissionRequest struct {
	PermissionID int `json:"permission_id" binding:"required"`
}

type AssignUserRoleRequest struct {
	RoleID int `json:"role_id" binding:"required"`
}

type AssignUserApplicationRequest struct {
	ApplicationID int `json:"application_id" binding:"required"`
}

type AssignUserApplicationRoleRequest struct {
	ApplicationID int  `json:"application_id" binding:"required"`
	RoleID        int  `json:"role_id" binding:"required"`
	ProfileID     *int `json:"profile_id,omitempty"`
}

func NewAdminHandler(dbService *db.Service) *AdminHandler {
	return &AdminHandler{
		db: dbService,
	}
}

// User handlers

// GetUsers lista todos os usuários cadastrados
// @Summary Lista todos os usuários
// @Description Retorna uma lista com todos os usuários cadastrados no sistema
// @Tags Admin - Users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} db.User
// @Failure 500 {object} map[string]interface{}
// @Router /admin/users [get]
func (h *AdminHandler) GetUsers(c *gin.Context) {
	users, err := h.db.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get users",
			"details": err.Error(),
		})
		return
	}

	// Remove senhas do retorno por segurança
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, users)
}

// Application handlers

// CreateApplication creates a new application
// @Summary Create a new application
// @Description Create a new application in the system
// @Tags Admin - Applications
// @Accept json
// @Produce json
// @Security Bearer
// @Param application body CreateApplicationRequest true "Application data"
// @Success 201 {object} db.Application "Application created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - validation error or slug already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/applications [post]
func (h *AdminHandler) CreateApplication(c *gin.Context) {
	var req CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if slug already exists
	exists, err := h.db.ApplicationSlugExists(req.Slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Application slug already exists"})
		return
	}

	app, err := h.db.CreateApplication(req.Slug, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create application"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Application created successfully",
		"application": app,
	})
}

// GetApplications lists all applications
// @Summary Get all applications
// @Description Retrieve a list of all applications in the system
// @Tags Admin - Applications
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} db.Application "List of applications"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/applications [get]
func (h *AdminHandler) GetApplications(c *gin.Context) {
	applications, err := h.db.GetApplications()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get applications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"applications": applications})
}

// GetApplication retrieves a specific application by ID
// @Summary Get application by ID
// @Description Retrieve a specific application by its ID
// @Tags Admin - Applications
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Application ID"
// @Success 200 {object} db.Application "Application details"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid ID"
// @Failure 404 {object} map[string]interface{} "Application not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/applications/{id} [get]
func (h *AdminHandler) GetApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID"})
		return
	}

	app, err := h.db.GetApplicationByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"application": app})
}

func (h *AdminHandler) UpdateApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID"})
		return
	}

	var req UpdateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.db.UpdateApplication(id, req.Slug, req.Name, req.Description, req.Active)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application updated successfully"})
}

func (h *AdminHandler) DeleteApplication(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID"})
		return
	}

	err = h.db.DeleteApplication(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application deleted successfully"})
}

// Role handlers

// CreateRole creates a new role for an application
// @Summary Create a new role
// @Description Create a new role associated with an application
// @Tags Admin - Roles
// @Accept json
// @Produce json
// @Security Bearer
// @Param role body CreateRoleRequest true "Role data"
// @Success 201 {object} db.Role "Role created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - validation error, application not found, or slug already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/roles [post]
func (h *AdminHandler) CreateRole(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify application exists
	_, err := h.db.GetApplicationByID(req.ApplicationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application not found"})
		return
	}

	role, err := h.db.CreateRole(req.ApplicationID, req.Slug, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Role created successfully",
		"role":    role,
	})
}

// GetRoles lists all roles, optionally filtered by application
// @Summary Get all roles
// @Description Retrieve a list of all roles, optionally filtered by application_id
// @Tags Admin - Roles
// @Accept json
// @Produce json
// @Security Bearer
// @Param application_id query int false "Filter by application ID"
// @Success 200 {array} db.Role "List of roles"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/roles [get]
func (h *AdminHandler) GetRoles(c *gin.Context) {
	applicationIDStr := c.Query("application_id")

	var roles []db.Role
	var err error

	if applicationIDStr != "" {
		applicationID, parseErr := strconv.Atoi(applicationIDStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application_id parameter"})
			return
		}
		roles, err = h.db.GetRolesByApplication(applicationID)
	} else {
		roles, err = h.db.GetRoles()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get roles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"roles": roles})
}

func (h *AdminHandler) GetRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	role, err := h.db.GetRoleByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"role": role})
}

func (h *AdminHandler) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.db.UpdateRole(id, req.Slug, req.Name, req.Description, req.Active)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role updated successfully"})
}

func (h *AdminHandler) DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	err = h.db.DeleteRole(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}

// Permission handlers

// CreatePermission creates a new permission for an application
// @Summary Create a new permission
// @Description Create a new permission associated with an application
// @Tags Admin - Permissions
// @Accept json
// @Produce json
// @Security Bearer
// @Param permission body CreatePermissionRequest true "Permission data"
// @Success 201 {object} db.Permission "Permission created successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - validation error, application not found, or slug already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/permissions [post]
func (h *AdminHandler) CreatePermission(c *gin.Context) {
	var req CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify application exists
	_, err := h.db.GetApplicationByID(req.ApplicationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application not found"})
		return
	}

	permission, err := h.db.CreatePermission(req.ApplicationID, req.Name, req.Slug, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create permission"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Permission created successfully",
		"permission": permission,
	})
}

func (h *AdminHandler) GetPermissions(c *gin.Context) {
	applicationIDStr := c.Query("application_id")

	var permissions []db.Permission
	var err error

	if applicationIDStr != "" {
		applicationID, parseErr := strconv.Atoi(applicationIDStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application_id parameter"})
			return
		}
		permissions, err = h.db.GetPermissionsByApplication(applicationID)
	} else {
		permissions, err = h.db.GetPermissions()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get permissions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"permissions": permissions})
}

func (h *AdminHandler) GetPermission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission ID"})
		return
	}

	permission, err := h.db.GetPermissionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Permission not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"permission": permission})
}

func (h *AdminHandler) UpdatePermission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission ID"})
		return
	}

	var req UpdatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.db.UpdatePermission(id, req.Name, req.Slug, req.Description, req.Active)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update permission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission updated successfully"})
}

func (h *AdminHandler) DeletePermission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission ID"})
		return
	}

	err = h.db.DeletePermission(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete permission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission deleted successfully"})
}

// Role-Permission relationship handlers

// AssignPermissionToRole assigns a permission to a role
// @Summary Assign permission to role
// @Description Assign a permission to a specific role
// @Tags Admin - Role-Permission Relationships
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Role ID"
// @Param assignment body AssignRolePermissionRequest true "Permission assignment data"
// @Success 200 {object} map[string]interface{} "Permission assigned successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid ID or permission not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/roles/{id}/permissions [post]
func (h *AdminHandler) AssignPermissionToRole(c *gin.Context) {
	roleIDStr := c.Param("id")
	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	var req AssignRolePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.db.AssignPermissionToRole(roleID, req.PermissionID)
	if err != nil {
		if err.Error() == "permission already assigned to role" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign permission to role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission assigned to role successfully"})
}

func (h *AdminHandler) RemovePermissionFromRole(c *gin.Context) {
	roleIDStr := c.Param("id")
	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	permissionIDStr := c.Param("permissionId")
	permissionID, err := strconv.Atoi(permissionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission ID"})
		return
	}

	err = h.db.RemovePermissionFromRole(roleID, permissionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove permission from role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission removed from role successfully"})
}

func (h *AdminHandler) GetRolePermissions(c *gin.Context) {
	roleIDStr := c.Param("id")
	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	permissions, err := h.db.GetRolePermissions(roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get role permissions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"permissions": permissions})
}

// User-Role relationship handlers
func (h *AdminHandler) AssignRoleToUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req AssignUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.db.AssignRoleToUser(userID, req.RoleID)
	if err != nil {
		if err.Error() == "role already assigned to user" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign role to user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role assigned to user successfully"})
}

func (h *AdminHandler) RemoveRoleFromUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	roleIDStr := c.Param("roleId")
	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	err = h.db.RemoveRoleFromUser(userID, roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove role from user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role removed from user successfully"})
}

// User-Application relationship handlers

// AssignApplicationToUser assigns an application to a user
// @Summary Assign application to user
// @Description Assign an application to a specific user
// @Tags Admin - User-Application Relationships
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "User ID"
// @Param assignment body AssignUserApplicationRequest true "Application assignment data"
// @Success 200 {object} map[string]interface{} "Application assigned successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid ID or application not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/users/{id}/applications [post]
func (h *AdminHandler) AssignApplicationToUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req AssignUserApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify user exists
	_, err = h.db.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// Verify application exists
	_, err = h.db.GetApplicationByID(req.ApplicationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application not found"})
		return
	}

	err = h.db.AssignUserToApplication(userID, req.ApplicationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign application to user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application assigned to user successfully"})
}

// RemoveApplicationFromUser removes an application from a user
// @Summary Remove application from user
// @Description Remove an application assignment from a specific user
// @Tags Admin - User-Application Relationships
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "User ID"
// @Param applicationId path int true "Application ID"
// @Success 200 {object} map[string]interface{} "Application removed successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/users/{id}/applications/{applicationId} [delete]
func (h *AdminHandler) RemoveApplicationFromUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	applicationIDStr := c.Param("applicationId")
	applicationID, err := strconv.Atoi(applicationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID"})
		return
	}

	err = h.db.RemoveUserFromApplication(userID, applicationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove application from user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application removed from user successfully"})
}

// GetUserApplications gets all applications for a user
// @Summary Get user applications
// @Description Retrieve all applications assigned to a specific user
// @Tags Admin - User-Application Relationships
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "User ID"
// @Success 200 {array} db.Application "List of user applications"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/users/{id}/applications [get]
func (h *AdminHandler) GetUserApplications(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	applications, err := h.db.GetUserApplications(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user applications"})
		return
	}

	c.JSON(http.StatusOK, applications)
}

// GetApplicationUsers gets all users for an application
// @Summary Get application users
// @Description Retrieve all users assigned to a specific application
// @Tags Admin - User-Application Relationships
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Application ID"
// @Success 200 {array} db.User "List of application users"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/applications/{id}/users [get]
func (h *AdminHandler) GetApplicationUsers(c *gin.Context) {
	applicationIDStr := c.Param("id")
	applicationID, err := strconv.Atoi(applicationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID"})
		return
	}

	users, err := h.db.GetApplicationUsers(applicationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get application users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

// User-Application-Role relationship handlers (New Model)

// AssignUserApplicationRole assigns a role to a user in a specific application
// @Summary Assign role to user in application
// @Description Assign a specific role to a user within an application context
// @Tags Admin - User-Application-Role Relationships
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "User ID"
// @Param assignment body AssignUserApplicationRoleRequest true "User application role assignment data"
// @Success 200 {object} map[string]interface{} "Role assigned successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid ID or role/application not found"
// @Failure 409 {object} map[string]interface{} "Conflict - user already has this role in this application"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/users/{id}/application-roles [post]
func (h *AdminHandler) AssignUserApplicationRole(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req AssignUserApplicationRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify user exists
	_, err = h.db.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// Verify application exists
	_, err = h.db.GetApplicationByID(req.ApplicationID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application not found"})
		return
	}

	// Verify role exists and belongs to the application
	role, err := h.db.GetRoleByID(req.RoleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role not found"})
		return
	}

	if role.ApplicationID != req.ApplicationID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role does not belong to the specified application"})
		return
	}

	err = h.db.AssignUserApplicationRole(userID, req.ApplicationID, req.RoleID, req.ProfileID)
	if err != nil {
		if err.Error() == "user already has this role in this application" {
			c.JSON(http.StatusConflict, gin.H{"error": "User already has this role in this application"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign role to user in application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role assigned to user in application successfully"})
}

// RemoveUserApplicationRole removes a role from a user in a specific application
// @Summary Remove role from user in application
// @Description Remove a specific role from a user within an application context
// @Tags Admin - User-Application-Role Relationships
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "User ID"
// @Param applicationId path int true "Application ID"
// @Param roleId path int true "Role ID"
// @Success 200 {object} map[string]interface{} "Role removed successfully"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/users/{id}/application-roles/{applicationId}/{roleId} [delete]
func (h *AdminHandler) RemoveUserApplicationRole(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	applicationIDStr := c.Param("applicationId")
	applicationID, err := strconv.Atoi(applicationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID"})
		return
	}

	roleIDStr := c.Param("roleId")
	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	err = h.db.RemoveUserApplicationRole(userID, applicationID, roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove role from user in application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role removed from user in application successfully"})
}

// GetUserApplicationRoles gets all application-role assignments for a user
// @Summary Get user application roles
// @Description Retrieve all application-role assignments for a specific user
// @Tags Admin - User-Application-Role Relationships
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "User ID"
// @Success 200 {array} db.UserApplicationRole "List of user application role assignments"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/users/{id}/application-roles [get]
func (h *AdminHandler) GetUserApplicationRoles(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	assignments, err := h.db.GetUserApplicationRoles(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user application roles"})
		return
	}

	c.JSON(http.StatusOK, assignments)
}

// GetUserRolesInApplication gets all roles for a user in a specific application
// @Summary Get user roles in application
// @Description Retrieve all roles assigned to a user within a specific application
// @Tags Admin - User-Application-Role Relationships
// @Accept json
// @Produce json
// @Security Bearer
// @Param userId path int true "User ID"
// @Param applicationId path int true "Application ID"
// @Success 200 {array} db.Role "List of user roles in application"
// @Failure 400 {object} map[string]interface{} "Bad request - invalid ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /admin/application-roles/users/{userId}/applications/{applicationId}/roles [get]
func (h *AdminHandler) GetUserRolesInApplication(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	applicationIDStr := c.Param("applicationId")
	applicationID, err := strconv.Atoi(applicationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID"})
		return
	}

	roles, err := h.db.GetUserRolesInApplication(userID, applicationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user roles in application"})
		return
	}

	c.JSON(http.StatusOK, roles)
}
