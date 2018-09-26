package creds

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	c "github.com/hanks/awsudo-go/configs"
	"github.com/hanks/awsudo-go/pkg/aws"
	"github.com/hanks/awsudo-go/pkg/parser"
	"github.com/hanks/awsudo-go/pkg/provider"
)

var tr = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

type Creds struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

func (c *Creds) Encode() ([]byte, error) {
	data, err := json.Marshal(c)
	return data, err
}

func (c *Creds) Decode(data []byte) error {
	err := json.Unmarshal(data, c)
	return err
}

func (c *Creds) SetEnv() {
	os.Setenv("AWS_ACCESS_KEY_ID", c.AccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", c.SecretAccessKey)
	os.Setenv("AWS_SESSION_TOKEN", c.SessionToken)
}

func FetchCreds(user string, pass string, roleName string, conf *parser.Config) *Creds {
	client := &http.Client{Timeout: c.ReqTimeout * time.Second, Transport: tr}
	okta := provider.NewOkta(user, pass, conf.Provider.AUTH_API, conf.Provider.IDP_LOGIN_URL, client)
	sessionToken, err := okta.GetSessionToken()
	if err != nil {
		log.Fatalf("GetSessionToken Error: %v\n", err)
	}
	samlAssertion, err := okta.GetSAMLAssertion(sessionToken)
	if err != nil {
		log.Fatalf("GetSAMLAssertion Error: %v\n", err)
	}

	_, roleARN := conf.GetARN(roleName)
	if roleARN == "" {
		log.Fatalf("Can not find role name: %v", roleName)
	}
	pARN := conf.GetPrincipalArn(roleARN)
	if pARN == "" {
		log.Fatalf("Role ARN format is error: %v", roleARN)
	}
	sts, _ := aws.NewSTS()

	creds, err := sts.GetCredentials(roleARN, pARN, samlAssertion, conf.Provider.SESSION_DURATION)
	if err != nil {
		log.Fatalf("GetCredentials Error: %v", err)
	}

	if *creds.AccessKeyId == "" || *creds.SecretAccessKey == "" || *creds.SessionToken == "" {
		log.Fatalf("Credentials Error: %v", *creds)
	}

	result := &Creds{
		AccessKeyID:     *creds.AccessKeyId,
		SecretAccessKey: *creds.SecretAccessKey,
		SessionToken:    *creds.SessionToken,
	}
	return result
}
