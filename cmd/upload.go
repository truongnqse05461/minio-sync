/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/truongnqse05461/minio-sync/internal"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload files from local directory to destination MinIO Server",
	Long:  `Upload files from local directory to destination MinIO Server.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		inputFile, _ := cmd.Flags().GetString("file")
		localDir, _ := cmd.Flags().GetString("directory")
		srcClient, err := internal.InitializeMinIOClient(cfg.Src.Endpoint, cfg.Src.AccessKey, cfg.Src.SecretKey, cfg.Src.UseSSL)
		if err != nil {
			log.Fatalf("Could not initialize source MinIO client: %v", err)
		}

		destClient, err := internal.InitializeMinIOClient(cfg.Dest.Endpoint, cfg.Dest.AccessKey, cfg.Dest.SecretKey, cfg.Dest.UseSSL)
		if err != nil {
			log.Fatalf("Could not initialize destination MinIO client: %v", err)
		}

		log.Printf("SRC: %s online: %v\n", srcClient.EndpointURL(), srcClient.IsOnline())
		log.Printf("DEST: %s online: %v\n", destClient.EndpointURL(), destClient.IsOnline())
		fmt.Println("---------------------------------------------------------------------------")
		buckets := internal.ReadBucketName(inputFile)
		idLen := len(buckets)
		log.Printf("Found %d buckets from %s\n", idLen, inputFile)
		for i, b := range buckets {
			exist, err := srcClient.BucketExists(ctx, b)
			if err != nil {
				log.Fatalf("Failed to call BucketExists: %v", err)
			}
			if exist {
				log.Printf("%d/%d|Uploading from local directory [%s] bucket: [%s]\n", i+1, idLen, localDir, b)
				bucketDir := filepath.Join(localDir, b)
				err = internal.UploadFilesToBucket(destClient, b, bucketDir)
				if err != nil {
					log.Printf("Upload failed: %s - %v", b, err)
					continue
				}
				log.Printf("Upload success: %s\n", b)
			}
		}

		fmt.Println("Upload completed.")
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().StringP("file", "f", "test.csv", "List bucket file name (.csv)")
	uploadCmd.Flags().StringP("directory", "d", "tmp", "Local directory to get files")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
