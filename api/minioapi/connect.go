package minioapi

import (
	"os"
	"strconv"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func EnvGetOptions() MinioOptions {
	endpoint, ok := os.LookupEnv("MINIO_ENDPOINT")
	if !ok {
		panic("MINIO_ENDPOINT environment variable not set")
	}
	user, ok := os.LookupEnv("MINIO_USER")
	if !ok {
		panic("MINIO_USER environment variable not set")
	}

	password, ok := os.LookupEnv("MINIO_PASSWORD")
	if !ok {
		panic("MINIO_PASSWORD environment variable not set")
	}

	useSSL, ok := os.LookupEnv("MINIO_USE_SSL")
	if !ok {
		panic("MINIO_USE_SSL environment variable not set")
	}

	useSSLBool, err := strconv.ParseBool(useSSL)
	if err != nil {
		panic(err)
	}

	return MinioOptions{
		Endpoint:        endpoint,
		AccessKey:       user,
		SecretAccessKey: password,
		UseSSL:          useSSLBool,
	}
}

func Connect(options MinioOptions) (*minio.Client, error) {
	return minio.New(options.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(options.AccessKey, options.SecretAccessKey, ""),
		Secure: options.UseSSL,
	})
}
