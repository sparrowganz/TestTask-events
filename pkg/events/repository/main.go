package repository

import (
	"TestTask-events/pkg/app"
	"TestTask-events/pkg/db/cache"
	"TestTask-events/pkg/events"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type Repository struct {
	db    *sqlx.DB
	cache *cache.Cache
}

func New(core app.Core, db *sqlx.DB, c *cache.Cache) (*Repository, error) {

	r := &Repository{
		db:    db,
		cache: c,
	}

	err := r.initMigrations()
	if err != nil {
		return nil, err
	}

	r.startCollector(core)

	return r, nil
}

func (r Repository) startCollector(core app.Core) {
	//Start collector
	core.Group().Go(func() error {
		ticker := time.NewTicker(time.Second / 10)

		for {
			select {
			case <-ticker.C:
				err := r.save(r.cache.GetAll())
				if err != nil {
					log.Println("[ERROR] failed save to db: " + err.Error())
				}
			case <-core.Context().Done():
				return r.save(r.cache.GetAll())
			}
		}
	})
}

func (r Repository) initMigrations() error {
	_, err := r.db.Exec(`
CREATE TABLE IF NOT EXISTS events
(
    client_time String,
    device_id String,
    device_os String,
    session String,
    sequence Int, 
	event String,
    param_int Int , 
    param_str String ,
	ip String, 
	server_time String
) engine=Memory ;`)
	return err
}

func (r Repository) Collect(event *events.Event) {
	r.cache.Set(event)
}

func (r Repository) save(events []*events.Event) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
		INSERT INTO events 
		(client_time,device_id,device_os,session,sequence,event,param_int,param_str,ip,server_time)
		VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)`,
	)
	if err != nil {
		return err
	}

	for _, event := range events {
		_, err = stmt.Exec(
			event.ClientTime, event.DeviceID, event.DeviceOS,
			event.Session, event.Sequence, event.Event,
			event.ParamInt, event.ParamStr, event.IP, event.ServerTime,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
