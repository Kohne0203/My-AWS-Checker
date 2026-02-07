package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3ClientInterface defines the interface for interacting with AWS S3 service.
type S3ClientInterface interface {
	ListBuckets(ctx context.Context) ([]types.Bucket, error)
	GetBucketRegion(ctx context.Context, bucketName string) (string, error)
	GetBucketPublicAccessConfig(ctx context.Context, bucketName string) (*PublicAccessConfig, error)
}

// CheckerInterface defines the interface for Security Check for AWS S3.
type CheckerInterface interface {
	AuditAllBuckets(ctx context.Context) ([]BucketAuditResult, error)
	AuditBucket(ctx context.Context, bucketName string) BucketAuditResult
}
