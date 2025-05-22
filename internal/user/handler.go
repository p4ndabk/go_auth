package user

import (
	"database/sql"
	"encoding/json"
	"go_auth/pkg/response"
	"net/http"
)

func RegisterUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			response.Fail(w, http.StatusBadRequest, "JSON inválido")
			return
		}

		if user.Name == "" || user.Email == "" || user.Password == "" {
			response.Fail(w, http.StatusBadRequest, "Campos obrigatórios ausentes")
			return
		}

		exists, err := EmailExists(db, user.Email)
		if err != nil {
			response.Fail(w, http.StatusInternalServerError, "Erro ao verificar email")
			return
		}

		if exists {
			response.Fail(w, http.StatusConflict, "Email já cadastrado")
			return
		}

		if err := store(db, &user); err != nil {
			response.Fail(w, http.StatusInternalServerError, "Erro ao criar usuário")
			return
		}

		w.WriteHeader(http.StatusCreated)
		user.Password = ""
		response.Success(w, http.StatusCreated, user)
	}
}

func ShowUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			response.Fail(w, http.StatusBadRequest, "ID do usuário ausente")
			return
		}

		user, err := show(db, id)
		if err != nil {
			response.Fail(w, http.StatusInternalServerError, "Erro ao buscar usuário")
			return
		}

		if user == nil {
			response.Fail(w, http.StatusNotFound, "Usuário não encontrado")
			return
		}

		user.Password = ""
		response.Success(w, http.StatusOK, user)
	}
}
