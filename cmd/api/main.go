package main

import (
	"go_auth/config"
	"go_auth/internal/database"
	"go_auth/pkg/logs"
	"go_auth/routes"
)

func main() {
	config.LoadConfig()

	logs.Init(config.GetEnv("APP_ENV", "development"))

	env := config.GetEnv("APP_ENV", "development")
	port := config.GetEnv("PORT", "8080")

	db := database.InitDB()

	router := routes.SetupRouter(env, db)

	logs.Info("Iniciando o servidor na porta " + port)
	if err := router.Run(":" + port); err != nil {
		logs.Error("Erro ao iniciar o servidor: %v", err)
	}
}
