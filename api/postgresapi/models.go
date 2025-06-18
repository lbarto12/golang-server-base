package postgresapi

type PostgresOptions struct {
	Host               string
	Port               string
	User               string
	Pass               string
	Database           string
	MaxOpenConnections int64
}
