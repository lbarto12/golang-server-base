package minioapi

import "github.com/minio/minio-go/v7"

var minioClient *minio.Client

func Init(options MinioOptions) error {
	var err error
	minioClient, err = Connect(options)
	return err
}

func Client() *minio.Client {
	return minioClient
}
