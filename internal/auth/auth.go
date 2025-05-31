package auth

import (
	"database/sql"
	"net/http"
	"time"

	"go_auth/config"
	"go_auth/internal/user"
	"go_auth/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	
)

var jwtSecret = []byte(config.GetEnv("JWT_SECRET", "default"))

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Success(c.Writer, http.StatusBadRequest, gin.H{
				"error": "no authorized",
			})
			return
		}

		user, err := user.GetUserByEmail(db, req.Email)
		if err != nil {
			response.Fail(c.Writer, http.StatusUnauthorized, "no authorized")
			return
		}

		if !CheckPasswordHash(req.Password, user.Password) {
			response.Fail(c.Writer, http.StatusUnauthorized, "no authorized")
			return
		}
		permissions, err := user.GetUserPermissions(db, user.ID)

		if err != nil {
			response.Fail(c.Writer, http.StatusInternalServerError, "Erro ao buscar permissões do usuário")
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": user.ID,
			"email":   user.Email,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
			"permissions": permissions,
		})

		tokenString, err := token.SignedString(jwtSecret)

		if err != nil {
			response.Fail(c.Writer, http.StatusInternalServerError, "Erro ao gerar token")
			return
		}

		response.Success(c.Writer, http.StatusOK, gin.H{"token": string(tokenString)})
	}
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
