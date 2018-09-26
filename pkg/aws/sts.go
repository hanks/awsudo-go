package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type STS struct {
	C *sts.STS
}

func NewSTS() (*STS, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}
	return &STS{
		C: sts.New(cfg),
	}, nil
}

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
