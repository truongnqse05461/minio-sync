package internal

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Initialize MinIO client
func InitializeMinIOClient(endpoint, accessKey, secretKey string, useSSL bool) (*minio.Client, error) {
	log.Printf("Initializing MinIO Client to %s using access key %s\n", endpoint, accessKey)
	return minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
}

func CountBucketObject(ctx context.Context, client *minio.Client, bucket string) int {
	objectCh := client.ListObjects(ctx, bucket, minio.ListObjectsOptions{Recursive: true})
	count := 0
	for object := range objectCh {
		if object.Err != nil {
			fmt.Printf("error listing object: %v", object.Err)
			continue
		}
		count++
	}
	return count
}

// Synchronize buckets
func SynchronizeBuckets(srcClient, destClient *minio.Client, srcBucket, destBucket string) error {
	ctx := context.Background()

	// Ensure the destination bucket exists
	err := destClient.MakeBucket(ctx, destBucket, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := destClient.BucketExists(ctx, destBucket)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s already exists\n", destBucket)
		} else {
			return fmt.Errorf("could not create bucket: %v", err)
		}
	}

	// List all objects in the source bucket
	objectCh := srcClient.ListObjects(ctx, srcBucket, minio.ListObjectsOptions{Recursive: true})
	count := 0
	for object := range objectCh {
		if object.Err != nil {
			return fmt.Errorf("error listing object: %v", object.Err)
		}

		// Get object from the source bucket
		objectReader, err := srcClient.GetObject(ctx, srcBucket, object.Key, minio.GetObjectOptions{})
		if err != nil {
			return fmt.Errorf("could not get object: %v", err)
		}

		// Upload object to the destination bucket
		_, err = destClient.PutObject(ctx, destBucket, object.Key, objectReader, object.Size, minio.PutObjectOptions{})
		if err != nil {
			return fmt.Errorf("could not put object: %v", err)
		}
		count++
		// Close the object reader
		objectReader.Close()
	}

	return nil
}

func SyncBucketToLocal(srcMinioClient *minio.Client, srcBucketName, localDir string) error {
	ctx := context.Background()

	// Create the local directory for the bucket if it doesn't exist
	bucketDir := filepath.Join(localDir, srcBucketName)
	if err := os.MkdirAll(bucketDir, os.ModePerm); err != nil {
		return err
	}

	// List all objects from the source bucket
	objectCh := srcMinioClient.ListObjects(ctx, srcBucketName, minio.ListObjectsOptions{})

	// Loop through the objects and save them to local files
	for object := range objectCh {
		if object.Err != nil {
			return object.Err
		}

		// Get object from the source bucket
		srcObject, err := srcMinioClient.GetObject(ctx, srcBucketName, object.Key, minio.GetObjectOptions{})
		if err != nil {
			return err
		}
		defer srcObject.Close()

		// Create local file
		localFilePath := filepath.Join(bucketDir, object.Key)
		localFile, err := os.Create(localFilePath)
		if err != nil {
			return err
		}
		defer localFile.Close()

		// Copy object data to local file
		if _, err = io.Copy(localFile, srcObject); err != nil {
			return err
		}

		fmt.Printf("Successfully saved object: %s to local file: %s\n", object.Key, localFilePath)
	}

	return nil
}

func UploadFilesToBucket(destMinioClient *minio.Client, destBucketName, bucketDir string) error {
	ctx := context.Background()
	// Ensure the destination bucket exists
	err := destMinioClient.MakeBucket(ctx, destBucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := destMinioClient.BucketExists(ctx, destBucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s already exists\n", destBucketName)
		} else {
			return fmt.Errorf("could not create bucket: %v", err)
		}
	}
	// Walk through the bucket directory and upload each file to the destination bucket
	err = filepath.Walk(bucketDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// Open local file
			localFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer localFile.Close()

			// Get the object key (relative path within the bucket directory)
			objectKey, err := filepath.Rel(bucketDir, path)
			if err != nil {
				return err
			}

			// Detect MIME type
			ext := strings.ToLower(filepath.Ext(path))
			mimeType := mime.TypeByExtension(ext)
			if mimeType == "" {
				mimeType = "application/octet-stream"
			}

			// Upload file to the destination bucket
			_, err = destMinioClient.PutObject(ctx, destBucketName, objectKey, localFile, info.Size(), minio.PutObjectOptions{
				ContentType: mimeType,
			})
			if err != nil {
				return err
			}

			fmt.Printf("Successfully uploaded file: %s to object: %s\n", path, objectKey)
		}

		return nil
	})

	return err
}
