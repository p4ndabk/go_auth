package application

import (
	"database/sql"
	"go_auth/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterApplicationHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var app Application
		if err := c.ShouldBindJSON(&app); err != nil {
			response.Fail(c.Writer, http.StatusBadRequest, "JSON inválido")
			return
		}

		if app.Slug == "" || app.Name == "" {
			response.Fail(c.Writer, http.StatusBadRequest, "Campos obrigatórios ausentes")
			return
		}

		if err := store(db, &app); err != nil {
			response.Fail(c.Writer, http.StatusInternalServerError, "Erro ao criar aplicação")
			return
		}

		c.Status(http.StatusCreated)
		response.Success(c.Writer, http.StatusCreated, app)
	}
}

func ShowApplicationHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		if idStr == "" {
			response.Fail(c.Writer, http.StatusBadRequest, "ID da aplicação ausente")
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.Fail(c.Writer, http.StatusBadRequest, "ID inválido")
			return
		}

		app, err := show(db, id)
		if err != nil {
			response.Fail(c.Writer, http.StatusInternalServerError, "Erro ao buscar aplicação")
			return
		}

		if app == nil {
			response.Fail(c.Writer, http.StatusNotFound, "Aplicação não encontrada")
			return
		}

		response.Success(c.Writer, http.StatusOK, app)
	}
}
