// Package auth is a cross-cutting concern (JWT + password hashing), same
// infra exception as apierror/health/docs: no model/service, no dependency
// on any domain package. ApplicationAccess/ApplicationInfo are auth's own
// DTOs for what goes inside the token — deliberately decoupled from the
// access/application domain types so this package stays a leaf.
package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"go_auth/internal/apierror"
)

type ApplicationInfo struct {
	ID          uint   `json:"id"`
	UUID        string `json:"uuid"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

type ApplicationAccess struct {
	Application ApplicationInfo `json:"application"`
	Roles       []string        `json:"roles"`
	Permissions []string        `json:"permissions"`
}

type Claims struct {
	UserID       uint                `json:"user_id"`
	Email        string              `json:"email"`
	Applications []ApplicationAccess `json:"applications"`
	jwt.RegisteredClaims
}

type UserClaims struct {
	UserID       uint
	Email        string
	Applications []ApplicationAccess
}

type Service struct {
	jwtSecret []byte
}

func New(jwtSecret string) *Service {
	return &Service{jwtSecret: []byte(jwtSecret)}
}

func (s *Service) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *Service) CheckPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (s *Service) GenerateToken(userID uint, email string, applications []ApplicationAccess) (string, error) {
	claims := &Claims{
		UserID:       userID,
		Email:        email,
		Applications: applications,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *Service) ValidateToken(tokenString string) (*UserClaims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return &UserClaims{
		UserID:       claims.UserID,
		Email:        claims.Email,
		Applications: claims.Applications,
	}, nil
}

func (s *Service) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			apierror.Respond(c, apierror.Unauthorized("missing_authorization", "authorization header required"))
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			apierror.Respond(c, apierror.Unauthorized("missing_bearer_token", "bearer token required"))
			c.Abort()
			return
		}

		claims, err := s.ValidateToken(tokenString)
		if err != nil {
			apierror.Respond(c, apierror.Unauthorized("invalid_token", "invalid token"))
			c.Abort()
			return
		}

		c.Set("user", claims)
		c.Next()
	}
}

func GetUserFromContext(c *gin.Context) (*UserClaims, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}

	userClaims, ok := user.(*UserClaims)
	return userClaims, ok
}
