package clickhouse

import (
	"fmt"
	"github.com/roistat/go-clickhouse"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func New(cfg Config) (*clickhouse.Conn, error) {

	transport := clickhouse.NewHttpTransport()
	conn := clickhouse.NewConn(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), transport)
	err := conn.Ping()
	if err != nil {
		return nil, err
	}

	err = initDB(conn)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func initDB(conn *clickhouse.Conn) error {
	q := clickhouse.NewQuery(`CREATE DATABASE IF NOT EXISTS test`)
	return q.Exec(conn)
}
