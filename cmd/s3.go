/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

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

	// get the list of S3 buckets
	output, err := clientBasic.ListBuckets(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	defer writer.Flush()

	fmt.Fprintln(writer, "BUCKET NAME\tREGION\tSTATUS")
	fmt.Fprintln(writer, "-----------\t------\t------")

	for _, bucket := range output {
		bucketName := aws.ToString(bucket.Name)
		region, err := getBucketRegion(clientBasic.S3Client, aws.ToString(bucket.Name))
		if err != nil {
			region = "ERROR"
		}
		// get the public access block configuration for the bucket
		publicAccess, err := clientBasic.S3Client.GetPublicAccessBlock(context.TODO(), &s3.GetPublicAccessBlockInput{
			Bucket: bucket.Name,
		})
		var status string
		if err != nil {
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) && apiErr.ErrorCode() == "NoSuchPublicAccessBlockConfiguration" {
				status = "WARNING - No Public Access"
			} else {
				status = "ERROR"
			}
		} else {
			status = checkBucketStatus(publicAccess.PublicAccessBlockConfiguration)
		}
		fmt.Fprintf(writer, "%s\t%s\t%s\n", bucketName, region, status)
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
