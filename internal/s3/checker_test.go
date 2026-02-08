package s3

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/assert"
)

func TestEvaluateStatus(t *testing.T) {
	tests := []struct {
		name       string
		config     *PublicAccessConfig
		wantStatus SecurityStatus
	}{
		{
			name: "正常系：完全な保護",
			config: &PublicAccessConfig{
				Config: &types.PublicAccessBlockConfiguration{
					BlockPublicAcls:       aws.Bool(true),
					BlockPublicPolicy:     aws.Bool(true),
					IgnorePublicAcls:      aws.Bool(true),
					RestrictPublicBuckets: aws.Bool(true),
				},
				Exists: true,
			},
			wantStatus: StatusSafe,
		},
		{
			name: "正常系：部分的な保護",
			config: &PublicAccessConfig{
				Config: &types.PublicAccessBlockConfiguration{
					BlockPublicAcls:       aws.Bool(false),
					BlockPublicPolicy:     aws.Bool(true),
					IgnorePublicAcls:      aws.Bool(true),
					RestrictPublicBuckets: aws.Bool(true),
				},
				Exists: true,
			},
			wantStatus: StatusWarningPartial,
		},
		{
			name: "正常系：設定が存在しない",
			config: &PublicAccessConfig{
				Config: nil,
				Exists: false,
			},
			wantStatus: StatusWarningNoConfig,
		},
		{
			name: "正常系：すべて無効",
			config: &PublicAccessConfig{
				Config: &types.PublicAccessBlockConfiguration{
					BlockPublicAcls:       aws.Bool(false),
					BlockPublicPolicy:     aws.Bool(false),
					IgnorePublicAcls:      aws.Bool(false),
					RestrictPublicBuckets: aws.Bool(false),
				},
				Exists: true,
			},
			wantStatus: StatusWarningPartial,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := evaluateStatus(tt.config)
			assert.Equal(t, tt.wantStatus, got)
		})
	}
}
