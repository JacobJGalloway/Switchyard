package postgres_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/JacobJGalloway/switchyard-go/internal/repository/postgres"
)

// testPool is the shared pgxpool.Pool for the test suite.
var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "switchyard_test",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(fmt.Sprintf("start postgres container: %v", err))
	}
	defer container.Terminate(ctx)

	host, err := container.Host(ctx)
	if err != nil {
		panic(fmt.Sprintf("container host: %v", err))
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		panic(fmt.Sprintf("container port: %v", err))
	}

	dbURL := fmt.Sprintf("postgres://test:test@%s:%s/switchyard_test", host, port.Port())

	// Run migrations. Go tests run from the package directory, so ../../migrations
	// resolves to Switchyard-Go/internal/migrations.
	migrateURL := fmt.Sprintf("pgx5://test:test@%s:%s/switchyard_test", host, port.Port())
	mig, err := migrate.New("file://../../migrations", migrateURL)
	if err != nil {
		panic(fmt.Sprintf("migrate init: %v", err))
	}
	if err := mig.Up(); err != nil && err != migrate.ErrNoChange {
		panic(fmt.Sprintf("migrate up: %v", err))
	}

	testPool, err = pgxpool.New(ctx, dbURL)
	if err != nil {
		panic(fmt.Sprintf("pgxpool: %v", err))
	}
	defer testPool.Close()

	m.Run()
}

func TestAnalyticsRepo_BOLsByStatus_Empty(t *testing.T) {
	repo := postgres.NewAnalyticsRepo(testPool)
	counts, err := repo.BOLsByStatus(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, counts)
}

func TestAnalyticsRepo_StopCompletionRate_NoStops(t *testing.T) {
	repo := postgres.NewAnalyticsRepo(testPool)
	rate, err := repo.StopCompletionRate(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 0.0, rate)
}

func TestAnalyticsRepo_FulfilledInWindow_Empty(t *testing.T) {
	repo := postgres.NewAnalyticsRepo(testPool)
	count, err := repo.FulfilledInWindow(context.Background(), time.Now().Add(-7*24*time.Hour))
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestAnalyticsRepo_OperatingCostByBOL_Empty(t *testing.T) {
	repo := postgres.NewAnalyticsRepo(testPool)
	costs, err := repo.OperatingCostByBOL(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, costs)
}

func TestAnalyticsRepo_OperatingCostByDriver_Empty(t *testing.T) {
	repo := postgres.NewAnalyticsRepo(testPool)
	costs, err := repo.OperatingCostByDriver(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, costs)
}

func TestAnalyticsRepo_OperatingCostByWarehouse_Empty(t *testing.T) {
	repo := postgres.NewAnalyticsRepo(testPool)
	costs, err := repo.OperatingCostByWarehouse(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, costs)
}
