package main

import (
	"log"

	"go_auth/config"
	"go_auth/routes"
)

func main() {
	config.LoadConfig()
	env := config.GetEnv("APP_ENV", "development")
	port := config.GetEnv("PORT", "8080")

	router := routes.SetupRouter(env)

	log.Printf("Servidor rodando na porta %s (env: %s)", port, env)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
