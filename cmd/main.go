package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Definir ambiente (desenvolvimento, produção, etc.)
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	// Criar instância do roteador
	router := gin.Default()

	// Rotas básicas (placeholder)
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"env":    env,
		})
	})

	// Porta padrão
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor rodando na porta %s (env: %s)", port, env)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
