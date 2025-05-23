package permission

import (
	"database/sql"
	"go_auth/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func IndexPermissionHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions, err := GetAllPermissions(db)
		if err != nil {
			response.Fail(c.Writer, http.StatusInternalServerError, "Erro ao buscar permissões")
			return
		}

		response.Success(c.Writer, http.StatusOK, permissions)
	}
}

func RegisterPermissionHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var p Permission
		if err := c.ShouldBindJSON(&p); err != nil {
			response.Fail(c.Writer, http.StatusBadRequest, "JSON inválido")
			return
		}

		if p.Slug == "" || p.Name == "" || p.ApplicationID == 0 {
			response.Fail(c.Writer, http.StatusBadRequest, "Campos obrigatórios ausentes")
			return
		}

		if err := CreatePermission(db, &p); err != nil {
			response.Fail(c.Writer, http.StatusInternalServerError, "Erro ao criar permissão")
			return
		}

		response.Success(c.Writer, http.StatusCreated, p)
	}
}

func ShowPermissionHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		if idStr == "" {
			response.Fail(c.Writer, http.StatusBadRequest, "ID da permissão ausente")
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.Fail(c.Writer, http.StatusBadRequest, "ID inválido")
			return
		}

		p, err := GetPermissionByID(db, id)
		if err != nil {
			response.Fail(c.Writer, http.StatusInternalServerError, "Erro ao buscar permissão")
			return
		}

		if p == nil {
			response.Fail(c.Writer, http.StatusNotFound, "Permissão não encontrada")
			return
		}

		response.Success(c.Writer, http.StatusOK, p)
	}
}

func UpdatePermissionHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.Fail(c.Writer, http.StatusBadRequest, "ID inválido")
			return
		}

		var p Permission
		if err := c.ShouldBindJSON(&p); err != nil {
			response.Fail(c.Writer, http.StatusBadRequest, "JSON inválido")
			return
		}

		p.ID = id

		if err := UpdatePermission(db, &p); err != nil {
			response.Fail(c.Writer, http.StatusInternalServerError, "Erro ao atualizar permissão")
			return
		}

		response.Success(c.Writer, http.StatusOK, p)
	}
}

func DeletePermissionHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.Fail(c.Writer, http.StatusBadRequest, "ID inválido")
			return
		}

		if err := DeletePermission(db, id); err != nil {
			response.Fail(c.Writer, http.StatusInternalServerError, "Erro ao deletar permissão")
			return
		}

		response.Success(c.Writer, http.StatusOK, map[string]interface{}{
			"message": "Permissão deletada com sucesso",
		})
	}
}
