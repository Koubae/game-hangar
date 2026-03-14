package postgres

import (
	"os"
	"testing"
)

func TestDatabasePostgresConfig_GetConnectionString(t *testing.T) {
	tests := []struct {
		name     string
		config   *DatabasePostgresConfig
		expected string
	}{
		{
			name: "all fields provided",
			config: &DatabasePostgresConfig{
				user:     "user",
				password: "pass",
				host:     "localhost",
				port:     5432,
				Database: "testdb",
				sslMode:  "disable",
			},
			expected: "postgres://user:pass@localhost:5432/testdb?sslmode=disable",
		},
		{
			name: "connectionString overrides",
			config: &DatabasePostgresConfig{
				connectionString: "postgres://overridden:5432/db",
				user:             "user",
				password:         "pass",
				host:             "localhost",
				port:             5432,
				Database:         "testdb",
				sslMode:          "disable",
			},
			expected: "postgres://overridden:5432/db",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := tt.config.GetConnectionString(); got != tt.expected {
					t.Errorf("GetConnectionString() = %v, want %v", got, tt.expected)
				}
			},
		)
	}
}

func TestDatabasePostgresConfig_String(t *testing.T) {
	config := &DatabasePostgresConfig{
		Driver:   "postgres",
		Database: "testdb",
		host:     "localhost",
		port:     5432,
	}
	expected := "DB[postgres] database:testdb connected @ localhost:5432"
	if got := config.String(); got != expected {
		t.Errorf("String() = %v, want %v", got, expected)
	}
}

func TestLoadConfig(t *testing.T) {
	os.Setenv("POSTGRES_DB", "env_db")
	os.Setenv("POSTGRES_HOST", "env_host")
	os.Setenv("POSTGRES_PORT", "9999")
	os.Setenv("POSTGRES_USER", "env_user")
	os.Setenv("POSTGRES_PASS", "env_pass")

	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.Database != "env_db" {
		t.Errorf("expected Database env_db, got %s", cfg.Database)
	}
	if cfg.host != "env_host" {
		t.Errorf("expected host env_host, got %s", cfg.host)
	}
	if cfg.port != 9999 {
		t.Errorf("expected port 9999, got %d", cfg.port)
	}
}
