package db

import (
	"database/sql"
	"fmt"

	"github.com/DenisOzindzheDev/chatix/auth/internal/config"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitPostgres(cfg *config.Config) error {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
	)
	var err error

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	if err = db.Ping(); err != nil {
		return fmt.Errorf("database ping: %w", err)
	}

	return nil
}

func GetDB() *sql.DB {
	return db
}
