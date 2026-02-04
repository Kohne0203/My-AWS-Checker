/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/spf13/cobra"
)

// s3Cmd represents the s3 command
var s3Cmd = &cobra.Command{
	Use:   "s3",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("S3 check start")
		checkBuckets()
	},
}

type BucketBasics struct {
	S3Client *s3.Client
}

func init() {
	rootCmd.AddCommand(s3Cmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// s3Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// s3Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func checkBuckets() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// create an S3 service client
	clientBasic := BucketBasics{
		S3Client: s3.NewFromConfig(cfg),
	}

	// get the list of objects in the bucket
	output, err := clientBasic.ListBuckets(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	for _, bucket := range output {
		log.Printf("bucket=%s", aws.ToString(bucket.Name))
		// 各バケットのパブリックアクセス設定を取得する
		publicAccess, err := clientBasic.S3Client.GetPublicAccessBlock(context.TODO(), &s3.GetPublicAccessBlockInput{
			Bucket: bucket.Name,
		})
		if err != nil {
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) && apiErr.ErrorCode() == "NoSuchPublicAccessBlockConfiguration" {
				fmt.Printf("No public access block configuration found for bucket %s\n", aws.ToString(bucket.Name))
				continue
			} else {
				fmt.Printf("Could not get public access block for bucket %s: %v", aws.ToString(bucket.Name), err)
			}
		}
		status := checkBucketStatus(publicAccess.PublicAccessBlockConfiguration)
		fmt.Printf("Status: %s\n", status)
	}
}

func (basics BucketBasics) ListBuckets(ctx context.Context) ([]types.Bucket, error) {
	var err error
	var output *s3.ListBucketsOutput
	var buckets []types.Bucket
	bucketPaginator := s3.NewListBucketsPaginator(basics.S3Client, &s3.ListBucketsInput{})
	for bucketPaginator.HasMorePages() {
		output, err = bucketPaginator.NextPage(ctx)
		if err != nil {
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) && apiErr.ErrorCode() == "AccessDenied" {
				fmt.Println("You don't have permission to access this bucket for this account")
				err = apiErr
			} else {
				log.Printf("Could not list buckets: %v", err)
			}
			break
		} else {
			buckets = append(buckets, output.Buckets...)
		}
	}
	return buckets, err
}

func checkBucketStatus(config *types.PublicAccessBlockConfiguration) string {
	if config == nil {
		return "Warning - No Configuration"
	}
	if *config.BlockPublicAcls && *config.BlockPublicPolicy && *config.IgnorePublicAcls && *config.RestrictPublicBuckets {
		return "Safe"
	}
	return "Warning - Partial Configuration"
}
