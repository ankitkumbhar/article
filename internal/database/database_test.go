package database_test

import (
	"article/internal/database"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_InitDB(t *testing.T) {
	tests := []struct {
		name    string
		loadEnv func(t *testing.T)
	}{
		{
			name:    "success - with fallback dsn values",
			loadEnv: func(t *testing.T) {},
		},
		{
			name: "success - with predefined env",
			loadEnv: func(t *testing.T) {
				t.Setenv("DB_USERNAME", "test-username")
				t.Setenv("DB_PASSWORD", "test-password")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.loadEnv(t)
			gotResp, err := database.InitDB()

			assert.Nil(t, err)
			assert.NotNil(t, gotResp)
		})
	}
}
