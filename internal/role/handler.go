package role

import (
	"database/sql"
	"net/http"
	"strconv"

	"go_auth/pkg/response"

	"github.com/gin-gonic/gin"
)

func IndexRolesHandler(db *sql.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        roles, err := GetAllRoles(db)
        if err != nil {
            response.Fail(c.Writer, http.StatusInternalServerError, "Erro ao listar roles")
            return
        }

        response.Success(c.Writer, http.StatusOK, roles)
    }
}

func RegisterRoleHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var r Role
		if err := c.ShouldBindJSON(&r); err != nil {
			response.Fail(c.Writer, http.StatusBadRequest, "JSON inválido")
			return
		}

		if r.Slug == "" || r.Name == "" || r.ApplicationID == 0 {
			response.Fail(c.Writer, http.StatusBadRequest, "Campos obrigatórios ausentes")
			return
		}

		if err := CreateRole(db, &r); err != nil {
			response.Fail(c.Writer, http.StatusInternalServerError, "Erro ao criar role")
			return
		}

		response.Success(c.Writer, http.StatusCreated, r)
	}
}

func ShowRoleHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.Fail(c.Writer, http.StatusBadRequest, "ID inválido")
			return
		}

		role, err := GetRoleByID(db, id)
		if err != nil {
			response.Fail(c.Writer, http.StatusInternalServerError, "Erro ao buscar role")
			return
		}

		if role == nil {
			response.Fail(c.Writer, http.StatusNotFound, "Role não encontrada")
			return
		}

		response.Success(c.Writer, http.StatusOK, role)
	}
}

func UpdateRoleHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.Fail(c.Writer, http.StatusBadRequest, "ID inválido")
			return
		}

		var r Role
		if err := c.ShouldBindJSON(&r); err != nil {
			response.Fail(c.Writer, http.StatusBadRequest, "JSON inválido")
			return
		}

		r.ID = id

		if r.Slug == "" || r.Name == "" || r.ApplicationID == 0 {
			response.Fail(c.Writer, http.StatusBadRequest, "Campos obrigatórios ausentes")
			return
		}

		err = UpdateRole(db, &r)
		if err != nil {
			response.Fail(c.Writer, http.StatusInternalServerError, "Erro ao atualizar role")
			return
		}

		response.Success(c.Writer, http.StatusOK, r)
	}
}

func DeleteRoleHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.Fail(c.Writer, http.StatusBadRequest, "ID inválido")
			return
		}

		err = DeleteRole(db, id)
		if err != nil {
			response.Fail(c.Writer, http.StatusInternalServerError, "Erro ao deletar role")
			return
		}

		c.Status(http.StatusNoContent)
	}
}


func AttachPermissionsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleIDStr := c.Param("id")
		roleID, err := strconv.Atoi(roleIDStr)
		if err != nil {
			response.Fail(c.Writer, http.StatusBadRequest, "ID inválido")
			return
		}

		var payload struct {
			PermissionIDs []int `json:"permission_ids"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			response.Fail(c.Writer, http.StatusBadRequest, "JSON inválido")
			return
		}

		if len(payload.PermissionIDs) == 0 {
			response.Fail(c.Writer, http.StatusBadRequest, "permission_ids é obrigatório")
			return
		}

		err = AttachPermissions(db, roleID, payload.PermissionIDs)
		if err != nil {
			response.Fail(c.Writer, http.StatusInternalServerError, "Erro ao anexar permissões")
			return
		}

		response.Success(c.Writer, http.StatusOK, gin.H{"message": "Permissões anexadas com sucesso"})
	}
}
