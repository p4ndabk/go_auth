package user

import (
	"database/sql"
	"go_auth/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterUserHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			response.Fail(c.Writer, http.StatusBadRequest, "JSON inválido")
			return
		}

		if user.Name == "" || user.Email == "" || user.Password == "" {
			response.Fail(c.Writer, http.StatusBadRequest, "Campos obrigatórios ausentes")
			return
		}

		exists, err := EmailExists(db, user.Email)
		if err != nil {
			response.Fail(c.Writer, http.StatusInternalServerError, "Erro ao verificar email")
			return
		}

		if exists {
			response.Fail(c.Writer, http.StatusConflict, "Email já cadastrado")
			return
		}

		if err := store(db, &user); err != nil {
			response.Fail(c.Writer, http.StatusInternalServerError, "Erro ao criar usuário")
			return
		}

		user.Password = ""
		response.Success(c.Writer, http.StatusCreated, user)
	}
}

func ShowUserHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		if idStr == "" {
			response.Fail(c.Writer, http.StatusBadRequest, "ID do usuário ausente")
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.Fail(c.Writer, http.StatusBadRequest, "ID inválido")
			return
		}

		user, err := show(db, id)
		if err != nil {
			response.Fail(c.Writer, http.StatusInternalServerError, "Erro ao buscar usuário")
			return
		}

		if user == nil {
			response.Fail(c.Writer, http.StatusNotFound, "Usuário não encontrado")
			return
		}

		user.Password = ""
		response.Success(c.Writer, http.StatusOK, user)
	}
}
