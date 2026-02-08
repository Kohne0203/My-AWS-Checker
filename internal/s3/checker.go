package s3

import "context"

type Checker struct {
	client S3ClientInterface
}

func NewChecker(client S3ClientInterface) *Checker {
	return &Checker{
		client: client,
	}
}

func (c *Checker) AuditBucket(ctx context.Context, bucketName string) BucketAuditResult {
	region, err := c.client.GetBucketRegion(ctx, bucketName)
	if err != nil {
		return BucketAuditResult{
			BucketName: bucketName,
			Status:     StatusError,
			Error:      err,
		}
	}

	config, err := c.client.GetBucketPublicAccessConfig(ctx, bucketName)
	if err != nil {
		return BucketAuditResult{
			BucketName: bucketName,
			Status:     StatusError,
			Error:      err,
		}
	}
	status := evaluateStatus(config)

	return BucketAuditResult{
		BucketName: bucketName,
		Status:     status,
		Region:     region,
		Error:      nil,
	}
}

func (c *Checker) AuditAllBuckets(ctx context.Context) ([]BucketAuditResult, error) {
	buckets, err := c.client.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}

	var results []BucketAuditResult
	for _, bucket := range buckets {
		result := c.AuditBucket(ctx, *bucket.Name)
		results = append(results, result)
	}
	return results, nil
}

func evaluateStatus(config *PublicAccessConfig) SecurityStatus {
	if !config.Exists {
		return StatusWarningNoConfig
	}

	if config.IsFullyProtected() {
		return StatusSafe
	}

	return StatusWarningPartial
}
