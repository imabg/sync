package database

import (
	"database/sql"

	"github.com/imabg/sync/pkg/logger"
	_ "github.com/lib/pq"
)

// CreatePostgresConnection accept POSTGRES_URI check the connection & retur the db instance
func (dbCtx DatabaseCtx) CreatePostgresConnection() (*sql.DB, error) {
	db, err := sql.Open("postgres", dbCtx.config.Env.PostgresURI)
	if e := db.Ping(); e != nil {
		return nil, e
	}
	logger.Log.InfoLog.Infoln("Postgres database is connected")
	return db, err
}

// DisconnectPostgresConnection closes the active conection
func (dbCtx DatabaseCtx) DisconnectPostgresConnection(db *sql.DB) error {
	return db.Close()
}
