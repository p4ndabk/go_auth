package user

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func RegisterUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
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

		exists, err := EmailExists(db, user.Email)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Erro ao verificar email",
			})
			return
		}

		if exists {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Email já cadastrado",
				"err": err,
			})
			return
		}

		if err := CreateUser(db, &user); err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "Erro ao criar usuário",
				"err": err,
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
