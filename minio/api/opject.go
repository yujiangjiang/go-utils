package api

import (
	"context"
	"fmt"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

var minioClient *minio.Client

func init() {
	endpoint := "10.37.20.122:35389"
	accessKeyID := "7LHWAAYWDNL268278OUa"
	secretAccessKey := "Y3JqcdYgCK2jX7Bg6ZCxKv61j3VWYvZc7DxAdmJd"

	accessKeyID = "minio"
	secretAccessKey = "minio123"
	useSSL := false
	client := NewClient(endpoint, accessKeyID, secretAccessKey, useSSL)
	if client == nil {
		panic("NewClient failed")
	}
	minioClient = client
}

func NewClient(endpoint, accessKeyID, secretAccessKey string, useSSL bool) *minio.Client {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	return minioClient
}

func CreateBucket(bucketName string, region string) error {
	exists, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	if !exists {
		if err := minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{
			Region: region,
		}); err != nil {
			return err
		}
	}
	return nil
}

func PutObject(bucketName, objectName string, filePath string) error {
	info, err := minioClient.FPutObject(context.Background(), bucketName, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	info.ETag = fmt.Sprintf("%d", info.Size)
	return nil
}
