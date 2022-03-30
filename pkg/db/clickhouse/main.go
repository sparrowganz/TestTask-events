package clickhouse

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mailru/go-clickhouse/v2"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func New(cfg Config) (*sqlx.DB, error) {
	c, err := sqlx.Open("chhttp", fmt.Sprintf("http://%s:%d/%s", cfg.Host, cfg.Port, cfg.Database))
	if err != nil {
		return nil, err
	}

	if err := c.Ping(); err != nil {
		return nil, err
	}

	return c, initDB(c)
}

func initDB(conn *sqlx.DB) error {
	_, err := conn.Exec(`CREATE DATABASE IF NOT EXISTS test`)
	return err
}
