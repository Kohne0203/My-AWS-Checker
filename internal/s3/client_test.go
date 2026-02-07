package s3

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/stretchr/testify/assert"
)

type mockS3APIClient struct {
	GetPublicAccessBlockFunc func(ctx context.Context, params *s3.GetPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error)
	ListBucketsFunc          func(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
	GetBucketLocationFunc    func(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error)
}

func (m *mockS3APIClient) GetPublicAccessBlock(ctx context.Context, params *s3.GetPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {
	return m.GetPublicAccessBlockFunc(ctx, params, optFns...)
}

func (m *mockS3APIClient) ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	return m.ListBucketsFunc(ctx, params, optFns...)
}

func (m *mockS3APIClient) GetBucketLocation(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error) {
	return m.GetBucketLocationFunc(ctx, params, optFns...)
}

func TestGetBucketPublicAccessConfig(t *testing.T) {
	tests := []struct {
		name          string
		bucketName    string
		mockFunc      func(ctx context.Context, params *s3.GetPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error)
		wantExists    bool
		wantConfigNil bool
		wantErr       bool
	}{
		{
			name:       "正常系：設定が存在する",
			bucketName: "test-bucket",
			mockFunc: func(ctx context.Context, params *s3.GetPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {
				blockPublic := true
				return &s3.GetPublicAccessBlockOutput{
					PublicAccessBlockConfiguration: &types.PublicAccessBlockConfiguration{
						BlockPublicAcls:       &blockPublic,
						BlockPublicPolicy:     &blockPublic,
						IgnorePublicAcls:      &blockPublic,
						RestrictPublicBuckets: &blockPublic,
					},
				}, nil
			},
			wantExists:    true,
			wantConfigNil: false,
			wantErr:       false,
		},
		{
			name:       "異常系：設定が存在しない",
			bucketName: "bucket-no-config",
			mockFunc: func(ctx context.Context, params *s3.GetPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {
				return nil, &smithy.GenericAPIError{
					Code:    "NoSuchPublicAccessBlockConfiguration",
					Message: "The public access block configuration was not found",
				}
			},
			wantExists:    false,
			wantConfigNil: true,
			wantErr:       false,
		},
		{
			name:       "異常系：エラーが発生する",
			bucketName: "error-bucket",
			mockFunc: func(ctx context.Context, params *s3.GetPublicAccessBlockInput, optFns ...func(*s3.Options)) (*s3.GetPublicAccessBlockOutput, error) {
				return nil, &smithy.GenericAPIError{
					Code:    "AccessDenied",
					Message: "AccessDenied",
				}
			},
			wantExists:    false,
			wantConfigNil: true,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockS3APIClient{
				GetPublicAccessBlockFunc: tt.mockFunc,
			}
			client := NewClient(mockClient)
			result, err := client.GetBucketPublicAccessConfig(context.Background(), tt.bucketName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantExists, result.Exists)
				assert.Equal(t, tt.wantConfigNil, result.Config == nil)
			}
		})
	}
}
