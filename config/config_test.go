package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvWithExistingKey(t *testing.T) {
	os.Setenv("TEST_KEY", "test_value")
	defer os.Unsetenv("TEST_KEY")

	result := GetEnv("TEST_KEY", "default")
	assert.Equal(t, "test_value", result, "Deve retornar o valor da variável de ambiente")
}

func TestGetEnvWithMissingKey(t *testing.T) {
	os.Unsetenv("MISSING_KEY")

	result := GetEnv("MISSING_KEY", "default")
	assert.Equal(t, "default", result, "Deve retornar o valor padrão quando a variável não existe")
}
