package db

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func (pgCfg *PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", pgCfg.Host, pgCfg.Port, pgCfg.User, pgCfg.Password, pgCfg.DBName, pgCfg.SSLMode)
}

func (pgCfg *PostgresConfig) Open() (*sql.DB, error) {
	db, err := sql.Open("pgx", pgCfg.String())
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	return db, nil
}
