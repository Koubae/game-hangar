package postgres

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
)

func TestConnectorPostgres_Ping(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	connector := &ConnectorPostgres{
		Pool:   mock,
		config: &DatabasePostgresConfig{Database: "testdb", host: "localhost", port: 5432},
	}

	mock.ExpectPing()

	if err := connector.Ping(context.Background()); err != nil {
		t.Errorf("Ping() error = %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestConnectorPostgres_Shutdown(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}

	connector := &ConnectorPostgres{
		Pool: mock,
	}

	// pgxmock Close doesn't need to be expected as it doesn't return anything or talk to DB
	if err := connector.Shutdown(); err != nil {
		t.Errorf("Shutdown() error = %v", err)
	}
}

func TestConnectorPostgres_String(t *testing.T) {
	config := &DatabasePostgresConfig{
		Driver:   "postgres",
		Database: "testdb",
		host:     "localhost",
		port:     5432,
	}
	connector := &ConnectorPostgres{
		config: config,
	}
	expected := config.String()
	if got := connector.String(); got != expected {
		t.Errorf("String() = %v, want %v", got, expected)
	}
}
