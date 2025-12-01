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
		connstr	string
		run_mgs string

		config	*pgxpool.Config
		err		error
	)

	pg_url =  getEnv("DATABASE_URL", "")
	pg_port = getEnv("DATABASE_PORT", "5432")
	pg_host = getEnv("DATABASE_HOST", "localhost")
	pg_user = getEnv("DATABASE_USER", "admin")
	pg_pswd = getEnv("DATABASE_PSWD", "postgres")
	run_mgs = getEnv("RUN_MIGRATIONS", "")

	if pg_url == "" {
		// THIS is for LOCAL Dev
		connstr = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			pg_user, pg_pswd, pg_host, pg_port, "trading-sim",
		)
	} else {
		connstr = buildMigrationURL(pg_url)
	}

	config, err = pgxpool.ParseConfig(connstr)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v", err)
	}

	if err = migrate.Run(connstr); err != nil {
		log.Fatalf("Miration Run failed: %v", err)
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

func buildMigrationURL(baseURL string) string {
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Fatalf("Failed to parse DATABASE_URL: %v", err)
	}

	adminUser := getEnv("DATABASE_ADMIN_USER", "")
	adminPass := getEnv("DATABASE_ADMIN_PASS", "")

	// Only override if both admin user and pass are set
	if adminUser != "" && adminPass != "" {
		u.User = url.UserPassword(adminUser, adminPass)
		log.Printf("Using admin user '%s' for migrations", adminUser)
	} else if adminUser != "" || adminPass != "" {
		log.Println("Warning: Both DATABASE_ADMIN_USER and DATABASE_ADMIN_PASS must be set to override migration credentials")
	}

	// Reconstruct the URL
	return u.String()
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}