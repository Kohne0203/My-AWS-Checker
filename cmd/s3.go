/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/cobra"

	s3pkg "my-aws-checker/internal/s3"
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

	// Create a new S3 client using the loaded configuration
	s3Client := s3.NewFromConfig(cfg)

	client := s3pkg.NewClient(s3Client)
	checker := s3pkg.NewChecker(client)

	results, err := checker.AuditAllBuckets(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	defer writer.Flush()

	fmt.Fprintln(writer, "BUCKET NAME\tREGION\tSTATUS")
	fmt.Fprintln(writer, "-----------\t------\t------")

	for _, result := range results {
		fmt.Fprintf(writer, "%s\t%s\t%s\n", result.BucketName, result.Region, result.Status)
	}
}
