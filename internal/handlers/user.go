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
			http.Error(w, "Campos obrigatórios ausentes", http.StatusUnprocessableEntity)
			return
		}

		exists, err := models.EmailExists(db, user.Email)
		if err != nil {
			http.Error(w, "Erro ao verificar email", http.StatusInternalServerError)
			return
		}

		if exists {
			http.Error(w, "Email já cadastrado", http.StatusUnprocessableEntity)
			return
		}

		if err := models.CreateUser(db, &user); err != nil {
			http.Error(w, "Erro ao criar usuário", http.StatusInternalServerError)
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
