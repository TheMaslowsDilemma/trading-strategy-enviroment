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
		pg_port string
		pg_user	string
		pg_host	string
		pg_pswd	string
		pg_url	string
		pg_ssl	string

		config	*pgxpool.Config
		err		error
	)

	pg_url = getEnv("DATABASE_URL", "")
	pg_port = getEnv("DATABASE_PORT", "5432")
	pg_host = getEnv("DATABASE_HOST", "localhost")
	pg_user = getEnv("DATABASE_USER", "admin")
	pg_pswd = getEnv("DATABASE_PSWD", "postgres")
	if pg_url == "" {

		pg_ssl = "disable"
	} else {
		pg_ssl = "require"

	}
	pg_url = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%v",
		pg_user, pg_pswd, pg_host, pg_port, "trading-sim", pg_ssl,
	)
	fmt.Println(pg_url)


	if err = migrate.Run(pg_url); err != nil {
		log.Fatalf("Miration Run failed: %v", err)
	}

	config, err = pgxpool.ParseConfig(pg_url)
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