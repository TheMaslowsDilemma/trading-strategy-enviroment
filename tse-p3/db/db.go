package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"tse-p3/db/migrate"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Init() {
	var (
		pg_user	string
		pg_host	string
		pg_pswd	string
		connstr	string
		config	*pgxpool.Config
		err		error
	)

	pg_host = getEnv("DATABASE_HOST", "localhost")
	pg_user = getEnv("DATABASE_USER", "admin")
	pg_pswd = getEnv("DATABASE_PSWD", "postgres")
	
	connstr = fmt.Sprintf("postgres://%s:%s@%s:5432/trading-sim?sslmode=disable", pg_user, pg_pswd, pg_host)

	if err != nil {
		log.Fatalf("Failed to resolve migrations path: %v", err)
	}

	if err = migrate.Run(connstr); err != nil {
		log.Fatalf("Miration Run failed: %v", err)
	}

	config, err = pgxpool.ParseConfig(connstr)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v", err)
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute

	Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	log.Println("Database pool initialized")
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}