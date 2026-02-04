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
	Short: "Audit S3 bucket public access settings",
	Long: `Check all S3 buckets for public access configurations.
Reports buckets as "Safe" or "Warning" based on their
PublicAccessBlock settings.`,
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
		fmt.Printf("Bucket=%s\n", aws.ToString(bucket.Name))
		region, err := getBucketRegion(clientBasic.S3Client, aws.ToString(bucket.Name))
		if err != nil {
			fmt.Printf("Could not get region for bucket %s: %v", aws.ToString(bucket.Name), err)
			continue
		} else {
			fmt.Printf("Region=%s\n", region)
		}
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

func getBucketRegion(client *s3.Client, bucketName string) (string, error) {
	location, err := client.GetBucketLocation(context.TODO(), &s3.GetBucketLocationInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return "", err
	}
	region := string(location.LocationConstraint)
	if region == "" {
		region = "us-east-1"
	}
	return region, nil
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
