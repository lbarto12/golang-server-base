Example .ENV

```.env
# postgres
POSTGRES_HOST=localhost
POSTGRES_PORT=5050
POSTGRES_DATABASE=postgres-database-name
POSTGRES_USER=postgres_user
POSTGRES_PASSWORD=postgres_password
POSTGRES_MAX_OPEN_CONNECTIONS=75

# minio
MINIO_ENDPOINT=localhost:9000
MINIO_USER=minio_user
MINIO_PASSWORD=minio_password
MINIO_DEFAULT_BUCKET=default_bucket
MINIO_USE_SSL=false
```

# How to run

In the root directory, follow these steps:

Launch Postrgres and Minio
```bash
docker compose up
```

Launch the server
```bash
go run .
```


