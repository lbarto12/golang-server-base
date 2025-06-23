package postgresapi

import (
	"golang-server-base/api/postgresapi/models"

	"github.com/jmoiron/sqlx"
)

// database Package private, but globally accessible reference to postgres
var database *sqlx.DB

func Init(options models.PostgresOptions) error {
	var err error
	database, err = Connect(options)
	if err != nil {
		return err
	}
	database.SetMaxOpenConns(int(options.MaxOpenConnections))
	return nil
}

func Database() (*sqlx.DB, error) {
	return database, database.Ping()
}
