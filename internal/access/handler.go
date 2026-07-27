package access

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go_auth/internal/apierror"
	"go_auth/internal/application"
	"go_auth/internal/auth"
	"go_auth/internal/role"
	"go_auth/internal/user"
)

// Handler owns the endpoints that need more than one domain's data: the
// user-application-role relationship itself, and login/me — which need
// user credentials, the access aggregation, and JWT issuance together.
// Register (user-only) stays in the user package.
type Handler struct {
	Service      *Service
	Users        *user.Service
	Applications *application.Service
	Roles        *role.Service
	Auth         *auth.Service
}

func NewHandler(service *Service, users *user.Service, applications *application.Service, roles *role.Service, authService *auth.Service) *Handler {
	return &Handler{Service: service, Users: users, Applications: applications, Roles: roles, Auth: authService}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserInfo struct {
	ID           uint                    `json:"id"`
	UUID         string                  `json:"uuid"`
	Username     string                  `json:"username"`
	Email        string                  `json:"email"`
	Applications []UserApplicationAccess `json:"applications"`
}

type LoginResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

func toTokenApplications(access []UserApplicationAccess) []auth.ApplicationAccess {
	apps := make([]auth.ApplicationAccess, len(access))
	for i, a := range access {
		apps[i] = auth.ApplicationAccess{
			Application: auth.ApplicationInfo{
				ID:          a.Application.ID,
				UUID:        a.Application.UUID,
				Slug:        a.Application.Slug,
				Name:        a.Application.Name,
				Description: a.Application.Description,
				Active:      a.Application.Active,
			},
			Roles:       a.Roles,
			Permissions: a.Permissions,
		}
	}
	return apps
}

// Login authenticates a user and returns a JWT token
// @Summary User login
// @Description Authenticate user with email and password, returns JWT token with user info
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "User login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} apierror.Body
// @Failure 401 {object} apierror.Body
// @Router /api/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_request", err.Error()))
		return
	}

	u, err := h.Users.GetByEmail(req.Email)
	if err != nil {
		apierror.Respond(c, apierror.Unauthorized("invalid_credentials", "invalid credentials"))
		return
	}

	if !h.Auth.CheckPassword(req.Password, u.Password) {
		apierror.Respond(c, apierror.Unauthorized("invalid_credentials", "invalid credentials"))
		return
	}

	userAccess, err := h.Service.GetUserAccess(u.ID)
	if err != nil {
		apierror.Respond(c, err)
		return
	}

	token, err := h.Auth.GenerateToken(u.ID, u.Email, toTokenApplications(userAccess))
	if err != nil {
		apierror.Respond(c, err)
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User: UserInfo{
			ID:           u.ID,
			UUID:         u.UUID,
			Username:     u.Username,
			Email:        u.Email,
			Applications: userAccess,
		},
	})
}

// Me returns current user information
// @Summary Get current user
// @Description Get current authenticated user information including roles and permissions
// @Tags Authentication
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} UserInfo
// @Failure 401 {object} apierror.Body
// @Failure 404 {object} apierror.Body
// @Router /api/me [get]
func (h *Handler) Me(c *gin.Context) {
	claims, exists := auth.GetUserFromContext(c)
	if !exists {
		apierror.Respond(c, apierror.Unauthorized("missing_user", "user not found in context"))
		return
	}

	u, err := h.Users.GetByID(claims.UserID)
	if err != nil {
		apierror.Respond(c, err)
		return
	}

	userAccess, err := h.Service.GetUserAccess(u.ID)
	if err != nil {
		apierror.Respond(c, err)
		return
	}

	c.JSON(http.StatusOK, UserInfo{
		ID:           u.ID,
		UUID:         u.UUID,
		Username:     u.Username,
		Email:        u.Email,
		Applications: userAccess,
	})
}

type AssignRequest struct {
	ApplicationID uint  `json:"application_id" binding:"required"`
	RoleID        uint  `json:"role_id" binding:"required"`
	ProfileID     *uint `json:"profile_id,omitempty"`
}

// Assign grants a role to a user within an application
// @Summary Assign role to user in application
// @Description Assign a specific role to a user within an application context
// @Tags Admin - User-Application-Role Relationships
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "User ID"
// @Param assignment body AssignRequest true "Assignment data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} apierror.Body
// @Failure 409 {object} apierror.Body
// @Router /api/admin/users/{id}/application-roles [post]
func (h *Handler) Assign(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_id", "invalid user ID"))
		return
	}

	var req AssignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_request", err.Error()))
		return
	}

	if _, err := h.Users.GetByID(uint(userID)); err != nil {
		apierror.Respond(c, apierror.BadRequest("user_not_found", "user not found"))
		return
	}

	if _, err := h.Applications.GetByID(req.ApplicationID); err != nil {
		apierror.Respond(c, apierror.BadRequest("application_not_found", "application not found"))
		return
	}

	r, err := h.Roles.GetByID(req.RoleID)
	if err != nil {
		apierror.Respond(c, apierror.BadRequest("role_not_found", "role not found"))
		return
	}
	if r.ApplicationID != req.ApplicationID {
		apierror.Respond(c, apierror.BadRequest("application_mismatch", "role does not belong to the specified application"))
		return
	}

	if err := h.Service.Assign(uint(userID), req.ApplicationID, req.RoleID, req.ProfileID); err != nil {
		if errors.Is(err, ErrAlreadyAssigned) {
			apierror.Respond(c, apierror.Conflict("already_assigned", err.Error()))
			return
		}
		apierror.Respond(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "role assigned to user in application successfully"})
}

// Remove revokes a specific role from a user within an application
// @Summary Remove role from user in application
// @Description Remove a specific role from a user within an application context
// @Tags Admin - User-Application-Role Relationships
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "User ID"
// @Param applicationId path int true "Application ID"
// @Param roleId path int true "Role ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} apierror.Body
// @Router /api/admin/users/{id}/application-roles/{applicationId}/{roleId} [delete]
func (h *Handler) Remove(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_id", "invalid user ID"))
		return
	}
	applicationID, err := strconv.ParseUint(c.Param("applicationId"), 10, 64)
	if err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_application_id", "invalid application ID"))
		return
	}
	roleID, err := strconv.ParseUint(c.Param("roleId"), 10, 64)
	if err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_role_id", "invalid role ID"))
		return
	}

	if err := h.Service.Remove(uint(userID), uint(applicationID), uint(roleID)); err != nil {
		apierror.Respond(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "role removed from user in application successfully"})
}

// ListForUser lists a user's application-role assignments
// @Summary Get user application roles
// @Description Retrieve all application-role assignments for a specific user
// @Tags Admin - User-Application-Role Relationships
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "User ID"
// @Success 200 {array} UserApplicationRole
// @Failure 400 {object} apierror.Body
// @Router /api/admin/users/{id}/application-roles [get]
func (h *Handler) ListForUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_id", "invalid user ID"))
		return
	}

	assignments, err := h.Service.ListByUser(uint(userID))
	if err != nil {
		apierror.Respond(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"application_roles": assignments})
}

// UserAccess returns the consolidated access of a user
// @Summary Get consolidated user access
// @Description Retrieve, for every application the user holds a role in, the role slugs and the permission slugs inherited from them — the same aggregation embedded in the JWT at login
// @Tags Admin - User-Application-Role Relationships
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "User ID"
// @Success 200 {array} UserApplicationAccess
// @Failure 400 {object} apierror.Body
// @Router /api/admin/users/{id}/access [get]
func (h *Handler) UserAccess(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_id", "invalid user ID"))
		return
	}

	access, err := h.Service.GetUserAccess(uint(userID))
	if err != nil {
		apierror.Respond(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"access": access})
}

// ListApplicationsForUser lists the applications a user has access to
// @Summary Get user applications
// @Description Retrieve all applications assigned to a specific user
// @Tags Admin - User-Application-Role Relationships
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "User ID"
// @Success 200 {array} go_auth_internal_application.Application
// @Failure 400 {object} apierror.Body
// @Router /api/admin/users/{id}/applications [get]
func (h *Handler) ListApplicationsForUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_id", "invalid user ID"))
		return
	}

	apps, err := h.Service.ListApplicationsByUser(uint(userID))
	if err != nil {
		apierror.Respond(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"applications": apps})
}

// RemoveAllForUserInApplication revokes every role a user holds in an application
// @Summary Remove all roles from a user in an application
// @Description Revoke every role assignment a user holds within an application
// @Tags Admin - User-Application-Role Relationships
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "User ID"
// @Param applicationId path int true "Application ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} apierror.Body
// @Router /api/admin/users/{id}/applications/{applicationId} [delete]
func (h *Handler) RemoveAllForUserInApplication(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_id", "invalid user ID"))
		return
	}
	applicationID, err := strconv.ParseUint(c.Param("applicationId"), 10, 64)
	if err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_application_id", "invalid application ID"))
		return
	}

	if err := h.Service.RemoveAllForUserInApplication(uint(userID), uint(applicationID)); err != nil {
		apierror.Respond(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "application access removed from user successfully"})
}

// ListUsersForApplication lists the users assigned to an application
// @Summary Get application users
// @Description Retrieve all users assigned to a specific application
// @Tags Admin - User-Application-Role Relationships
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Application ID"
// @Success 200 {array} go_auth_internal_user.User
// @Failure 400 {object} apierror.Body
// @Router /api/admin/applications/{id}/users [get]
func (h *Handler) ListUsersForApplication(c *gin.Context) {
	applicationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_id", "invalid application ID"))
		return
	}

	users, err := h.Service.ListUsersByApplication(uint(applicationID))
	if err != nil {
		apierror.Respond(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// ListRolesForUserInApplication lists the roles a user holds within an application
// @Summary Get user roles in application
// @Description Retrieve all roles assigned to a user within a specific application
// @Tags Admin - User-Application-Role Relationships
// @Accept json
// @Produce json
// @Security Bearer
// @Param userId path int true "User ID"
// @Param applicationId path int true "Application ID"
// @Success 200 {array} go_auth_internal_role.Role
// @Failure 400 {object} apierror.Body
// @Router /api/admin/application-roles/users/{userId}/applications/{applicationId}/roles [get]
func (h *Handler) ListRolesForUserInApplication(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_id", "invalid user ID"))
		return
	}
	applicationID, err := strconv.ParseUint(c.Param("applicationId"), 10, 64)
	if err != nil {
		apierror.Respond(c, apierror.BadRequest("invalid_application_id", "invalid application ID"))
		return
	}

	roles, err := h.Service.ListRolesForUserInApplication(uint(userID), uint(applicationID))
	if err != nil {
		apierror.Respond(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"roles": roles})
}
