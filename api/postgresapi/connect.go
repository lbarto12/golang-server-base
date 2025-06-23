package postgresapi

import (
	"fmt"
	"golang-server-base/api/postgresapi/models"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func EnvGetOptions() models.PostgresOptions {
	host, ok := os.LookupEnv("POSTGRES_HOST")
	if !ok {
		panic("POSTGRES_HOST environment variable not set")
	}

	port, ok := os.LookupEnv("POSTGRES_PORT")
	if !ok {
		panic("POSTGRES_PORT environment variable not set")
	}

	user, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		panic("POSTGRES_USER environment variable not set")
	}

	password, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		panic("POSTGRES_PASSWORD environment variable not set")
	}

	dbname, ok := os.LookupEnv("POSTGRES_DATABASE")
	if !ok {
		panic("POSTGRES_DBNAME environment variable not set")
	}

	maxOpenConnections, ok := os.LookupEnv("POSTGRES_MAX_OPEN_CONNECTIONS")
	if !ok {
		panic("POSTGRES_MAX_OPEN_CONNECTIONS environment variable not set")
	}

	maxOpenConnectionsInt, err := strconv.ParseInt(maxOpenConnections, 10, 64)
	if err != nil {
		panic("POSTGRES_MAX_OPEN_CONNECTIONS variable not a valid integer")
	}

	return models.PostgresOptions{
		Host:               host,
		Port:               port,
		User:               user,
		Pass:               password,
		Database:           dbname,
		MaxOpenConnections: maxOpenConnectionsInt,
	}
}

func Connect(options models.PostgresOptions) (*sqlx.DB, error) {

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", // TODO: at some point swap this out for an SSL environment variable?
		options.Host,
		options.Port,
		options.User,
		options.Pass,
		options.Database,
	)

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return db, nil
}
