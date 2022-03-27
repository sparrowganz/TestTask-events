package repository

import (
	"github.com/roistat/go-clickhouse"
)

type Repository struct {
	db *clickhouse.Conn
}

func New(db *clickhouse.Conn) *Repository {
	return &Repository{db: db}
}

func (r Repository) InitMigrations() error {
	query := clickhouse.NewQuery(`
CREATE TABLE IF NOT EXISTS events
(
    client_time Datetime,
    device_id String,
    device_os String,
    session String,
    sequence Int, 
	event String,
    param_int Int , 
    param_str String ,
	ip String, 
	server_time Datetime
) engine=Memory ;`)
	return query.Exec(r.db)
}

func (r Repository) Save(event interface{}) error {

	data := event.(map[string]interface{})

	q := clickhouse.NewQuery(
		"INSERT INTO events VALUES (?,?,?,?,?,?,?,?,?,?)",
		data["client_time"], data["device_id"], data["device_os"], data["session"], data["sequence"],
		data["event"], data["param_int"], data["param_str"], data["ip"], data["server_time"],
	)

	return q.Exec(r.db)
}
