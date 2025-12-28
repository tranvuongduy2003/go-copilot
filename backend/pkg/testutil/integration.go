package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type IntegrationTestSuite struct {
	DatabasePool     *pgxpool.Pool
	DatabaseURL      string
	CleanupFunctions []func()
}

func NewIntegrationTestSuite(t *testing.T) *IntegrationTestSuite {
	t.Helper()

	if os.Getenv("INTEGRATION_TEST") != "true" {
		t.Skip("Skipping integration test. Set INTEGRATION_TEST=true to run.")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = buildDatabaseURLFromEnv()
	}

	if databaseURL == "" {
		t.Skip("Skipping integration test. DATABASE_URL not set.")
	}

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	return &IntegrationTestSuite{
		DatabasePool:     pool,
		DatabaseURL:      databaseURL,
		CleanupFunctions: make([]func(), 0),
	}
}

func buildDatabaseURLFromEnv() string {
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	database := os.Getenv("DATABASE_NAME")
	sslMode := os.Getenv("DATABASE_SSL_MODE")

	if host == "" || user == "" || database == "" {
		return ""
	}

	if port == "" {
		port = "5432"
	}
	if sslMode == "" {
		sslMode = "disable"
	}

	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, database, sslMode,
	)
}

func (suite *IntegrationTestSuite) Cleanup(t *testing.T) {
	t.Helper()

	for i := len(suite.CleanupFunctions) - 1; i >= 0; i-- {
		suite.CleanupFunctions[i]()
	}

	if suite.DatabasePool != nil {
		suite.DatabasePool.Close()
	}
}

func (suite *IntegrationTestSuite) AddCleanup(cleanupFunction func()) {
	suite.CleanupFunctions = append(suite.CleanupFunctions, cleanupFunction)
}

func (suite *IntegrationTestSuite) TruncateTables(t *testing.T, tables ...string) {
	t.Helper()

	ctx := context.Background()

	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)
		_, err := suite.DatabasePool.Exec(ctx, query)
		if err != nil {
			t.Fatalf("Failed to truncate table %s: %v", table, err)
		}
	}
}

func (suite *IntegrationTestSuite) CleanAllTables(t *testing.T) {
	t.Helper()

	tables := []string{
		"audit_logs",
		"password_reset_tokens",
		"refresh_tokens",
		"user_roles",
		"role_permissions",
		"users",
		"roles",
		"permissions",
	}

	suite.TruncateTables(t, tables...)
}

func (suite *IntegrationTestSuite) ExecuteSQL(t *testing.T, query string, args ...interface{}) sql.Result {
	t.Helper()

	ctx := context.Background()
	commandTag, err := suite.DatabasePool.Exec(ctx, query, args...)
	if err != nil {
		t.Fatalf("Failed to execute SQL: %v", err)
	}

	return &pgxCommandTagResult{commandTag: commandTag}
}

type pgxCommandTagResult struct {
	commandTag interface{ RowsAffected() int64 }
}

func (result *pgxCommandTagResult) LastInsertId() (int64, error) {
	return 0, fmt.Errorf("LastInsertId not supported")
}

func (result *pgxCommandTagResult) RowsAffected() (int64, error) {
	return result.commandTag.RowsAffected(), nil
}

func (suite *IntegrationTestSuite) WithTransaction(t *testing.T, testFunction func(ctx context.Context)) {
	t.Helper()

	ctx := context.Background()
	transaction, err := suite.DatabasePool.Begin(ctx)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}

	defer func() {
		if err := transaction.Rollback(ctx); err != nil && err.Error() != "tx is closed" {
			t.Logf("Failed to rollback transaction: %v", err)
		}
	}()

	transactionContext := context.WithValue(ctx, transactionContextKey{}, transaction)
	testFunction(transactionContext)
}

type transactionContextKey struct{}

func (suite *IntegrationTestSuite) SeedTestData(t *testing.T, seedFunction func(pool *pgxpool.Pool) error) {
	t.Helper()

	if err := seedFunction(suite.DatabasePool); err != nil {
		t.Fatalf("Failed to seed test data: %v", err)
	}
}

func SetupIntegrationTest(t *testing.T) (*IntegrationTestSuite, func()) {
	t.Helper()

	suite := NewIntegrationTestSuite(t)

	suite.CleanAllTables(t)

	return suite, func() {
		suite.CleanAllTables(t)
		suite.Cleanup(t)
	}
}
