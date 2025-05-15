package integration

import (
	"context"
	"domashka-backend/config"
	"domashka-backend/pkg/postgres"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func setupTestDatabase() (*postgres.Postgres, error) {
	cfg := config.GetConfig()
	pg, err := postgres.New(cfg.DB.GetDSN(), postgres.MaxPoolSize(cfg.DB.PoolCapacity))
	if err != nil {
		return nil, err
	}
	_, err = pg.Pool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS integration_test_table(id BIGINT, name TEXT)")
	if err != nil {
		return nil, err
	}
	return pg, nil
}

func teardownTestDatabase(db *postgres.Postgres) {
	defer db.Close()
	_, err := db.Pool.Exec(context.Background(), "DROP TABLE IF EXISTS integration_test_table")
	if err != nil {
		log.Fatal(err)
	}
}

func TestIntegration_DatabaseInteraction(t *testing.T) {
	if testing.Short() {
		t.Skip("пропускаем интеграционный тест")
	}

	db, err := setupTestDatabase()
	if err != nil {
		t.Fatalf("не удалось инициализировать БД: %v", err)
	}
	defer teardownTestDatabase(db)

	id, name := int64(1), "name"
	_, err = db.Pool.Exec(context.Background(), "INSERT INTO integration_test_table (id, name) VALUES ($1, $2)", id, name)
	if err != nil {
		t.Fatal(err)
	}
	var (
		actualID   int64
		actualName string
	)
	err = db.Pool.QueryRow(context.Background(), "SELECT id, name FROM integration_test_table WHERE id = $1", id).Scan(&actualID, &actualName)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, id, actualID)
	require.Equal(t, name, actualName)
}
