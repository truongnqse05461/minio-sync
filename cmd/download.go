/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/truongnqse05461/minio-sync/internal"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download objects from source MinIO Server to local directory",
	Long:  `Download objects from source MinIO Server to local directory`,
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
				objectCount := internal.CountBucketObject(ctx, srcClient, b)
				if objectCount > 0 {
					log.Printf("%d/%d|Downloading to local directory [%s] bucket: [%s]\n", i+1, idLen, localDir, b)
					err = internal.SyncBucketToLocal(srcClient, b, localDir)
					if err != nil {
						log.Printf("Download failed: %s - %v", b, err)
						continue
					}
					log.Printf("Download success: %s\n", b)
				}
			}
		}

		fmt.Println("Download completed.")
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().StringP("file", "f", "test.csv", "List bucket file name (.csv)")
	downloadCmd.Flags().StringP("directory", "d", "tmp", "Local directory to save objects")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
