package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitSQLite(t *testing.T) {
	db := InitSQLite(":memory:")

	assert.NotNil(t, db, "O banco de dados não deve ser nil")

	err := db.Ping()
	assert.NoError(t, err, "Deve conectar ao banco de dados sem erros")
}

func TestGetDB(t *testing.T) {
	_ = InitSQLite(":memory:")

	dbInstance := GetDB()
	assert.NotNil(t, dbInstance, "GetDB deve retornar uma instância válida")
}
