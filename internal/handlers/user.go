package handlers

import (
	"database/sql"
	"encoding/json"
	"go_auth/internal/models"
	"net/http"
)

func RegisterUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		if user.Name == "" || user.Email == "" || user.Password == "" {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Campos obrigatórios ausentes",
			})
			return
		}

		exists, err := models.EmailExists(db, user.Email)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Erro ao verificar email",
			})
			return
		}

		if exists {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Email já cadastrado",
			})
			return
		}

		if err := models.CreateUser(db, &user); err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Erro ao criar usuário",
			})
			return
		}

		w.WriteHeader(http.StatusCreated)
		user.Password = ""
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Usuário criado com sucesso",
			"user":    user,
		})
	}
}
