package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// STS is a wrapper struct for sts.STS type
type STS struct {
	C *sts.STS
}

// NewSTS is to create a nwe sts instance to perform aws sts api
func NewSTS() (*STS, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}
	return &STS{
		C: sts.New(cfg),
	}, nil
}

// GetCredentials is to call assume api to get aws secret tokens for futher commands
func (s *STS) GetCredentials(roleARN string, principalArn string, samlAssertion string, duration int64) (*sts.Credentials, error) {
	input := &sts.AssumeRoleWithSAMLInput{
		DurationSeconds: aws.Int64(duration),
		PrincipalArn:    aws.String(principalArn),
		RoleArn:         aws.String(roleARN),
		SAMLAssertion:   aws.String(samlAssertion),
	}
	req := s.C.AssumeRoleWithSAMLRequest(input)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}
	return resp.Credentials, nil
}
