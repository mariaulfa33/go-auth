package usecase

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadEnv(t *testing.T) {
	key := "envKey"
	value := "envValue"

	t.Setenv(key, value)
	LoadEnv("../.env")

	assert.Equal(t, os.Getenv(key), value)
}

func TestLoadEnvError(t *testing.T) {
	err := LoadEnv("none")

	assert.NotNil(t, err)

}
