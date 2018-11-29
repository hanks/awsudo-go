package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// ExtendedSTS is a extended struct for sts.STS type
type ExtendedSTS struct {
	aSTS *sts.STS
}

// NewSTS is to create a nwe sts instance to perform aws sts api
func NewSTS() (*ExtendedSTS, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}
	return &ExtendedSTS{
		aSTS: sts.New(cfg),
	}, nil
}

// GetAWSCredentials is to call assume api to get aws secret tokens for futher commands
func (s *ExtendedSTS) GetAWSCredentials(roleARN string, principalArn string, samlAssertion string, duration int64) (*sts.Credentials, error) {
	input := &sts.AssumeRoleWithSAMLInput{
		DurationSeconds: aws.Int64(duration),
		PrincipalArn:    aws.String(principalArn),
		RoleArn:         aws.String(roleARN),
		SAMLAssertion:   aws.String(samlAssertion),
	}
	req := s.aSTS.AssumeRoleWithSAMLRequest(input)
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}
	return resp.Credentials, nil
}
