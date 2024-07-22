package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/joho/godotenv"
)

type Store struct {
	DB *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{
		DB: db,
	}
}

func Open(config SQLServerConfig) (*sql.DB, error) {
	db, err := sql.Open("sqlserver", config.String())
	if err != nil {
		return nil, err
	}
	return db, nil
}

func DefaultSQLServerConfig() SQLServerConfig {
	godotenv.Load(".env")

	return SQLServerConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("SSLMode"),
	}
}

type SQLServerConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (cfg SQLServerConfig) String() string {
	return fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}
