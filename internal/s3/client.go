package s3

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

type S3APIClient interface {
	GetPublicAccessBlock(ctx context.Context, params *s3.GetPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error)
	ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
	GetBucketLocation(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error)
}

type Client struct {
	s3Client S3APIClient
}

func NewClient(s3Client S3APIClient) *Client {
	return &Client{
		s3Client: s3Client,
	}
}

func (c *Client) ListBuckets(ctx context.Context) ([]types.Bucket, error) {
	var err error
	var output *s3.ListBucketsOutput
	var buckets []types.Bucket

	bucketPaginator := s3.NewListBucketsPaginator(c.s3Client, &s3.ListBucketsInput{})
	for bucketPaginator.HasMorePages() {
		output, err = bucketPaginator.NextPage(ctx)
		if err != nil {
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) && apiErr.ErrorCode() == "AccessDenied" {
				fmt.Println("You don't have permission to access this bucket for this account")
				err = apiErr
			} else {
				fmt.Printf("Could not list buckets: %v", err)
			}
			break
		} else {
			buckets = append(buckets, output.Buckets...)
		}
	}
	return buckets, err
}

func (c *Client) GetBucketRegion(ctx context.Context, bucketName string) (string, error) {
	location, err := c.s3Client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return "", err
	}
	region := string(location.LocationConstraint)
	// us-east-1は空文字で返ってくるため、変数に詰める必要がある
	if region == "" {
		region = "us-east-1"
	}
	return region, nil
}

func (c *Client) GetBucketPublicAccessConfig(ctx context.Context, bucketName string) (*PublicAccessConfig, error) {
	response, err := c.s3Client.GetPublicAccessBlock(ctx, &s3.GetPublicAccessBlockInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "NoSuchPublicAccessBlockConfiguration" {
			return &PublicAccessConfig{
				Config: nil,
				Exists: false,
			}, nil
		}
		return nil, err
	}

	return &PublicAccessConfig{
		Config: response.PublicAccessBlockConfiguration,
		Exists: true,
	}, nil
}
