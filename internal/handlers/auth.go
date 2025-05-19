package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"go_auth/config"
	"go_auth/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
			return
		}

		user, err := models.GetUserByEmail(db, req.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no authorized"})
			return
		}

		if user.Password != req.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no authorized"})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": user.ID,
			"email":   user.Email,
			"exp":     time.Now().Add(24 * time.Hour).Unix(), // expira em 1 dia
		})

		tokenString, err := token.SignedString(jwtSecret)

		fmt.Println("Jwt key:", jwtSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": string(tokenString)})
	}
}
