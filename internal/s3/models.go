package s3

import "github.com/aws/aws-sdk-go-v2/service/s3/types"

type BucketAuditResult struct {
	BucketName string
	Region     string
	Status     SecurityStatus
	Error      error
}

type SecurityStatus int

const (
	StatusUnknown SecurityStatus = iota
	StatusSafe
	StatusWarningNoConfig
	StatusWarningPartial
	StatusError
)

func (s SecurityStatus) String() string {
	switch s {
	case StatusUnknown:
		return "UNKNOWN"
	case StatusSafe:
		return "SAFE"
	case StatusWarningNoConfig:
		return "WARNING: NO CONFIGURATION"
	case StatusWarningPartial:
		return "WARNING: PARTIAL CONFIGURATION"
	case StatusError:
		return "ERROR"
	default:
		return "Unknown"
	}
}

type PublicAccessConfig struct {
	Config *types.PublicAccessBlockConfiguration
	Exists bool
}

func (p *PublicAccessConfig) IsFullyProtected() bool {
	if !p.Exists || p.Config == nil {
		return false
	}
	return *p.Config.BlockPublicAcls &&
		*p.Config.BlockPublicPolicy &&
		*p.Config.IgnorePublicAcls &&
		*p.Config.RestrictPublicBuckets
}
