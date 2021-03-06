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

// Creds is a struct contains necessary aws secret tokens
type Creds struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

// Encode is to transfer creds to json byte present
func (c *Creds) Encode() ([]byte, error) {
	data, err := json.Marshal(c)
	return data, err
}

// Decode is to restore creds from json bytes
func (c *Creds) Decode(data []byte) error {
	err := json.Unmarshal(data, c)
	return err
}

// NewCred is to create a new creds object from json bytes
func NewCred(data []byte) (*Creds, error) {
	c := &Creds{}
	err := c.Decode(data)
	return c, err
}

// SetEnv is to set aws tokens to environment variables, to be used by other aws commands
func (c *Creds) SetEnv() {
	os.Setenv("AWS_ACCESS_KEY_ID", c.AccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", c.SecretAccessKey)
	os.Setenv("AWS_SESSION_TOKEN", c.SessionToken)
}

// FetchCreds is to call aws apis to get credentials
func FetchCreds(user string, pass string, roleName string, conf *parser.Config) *Creds {
	client := &http.Client{Timeout: c.ReqTimeout * time.Second, Transport: tr}
	okta := provider.NewOkta(user, pass, conf.Provider.AuthAPI, conf.Provider.IdpLoginURL, client)
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

	creds, err := sts.GetAWSCredentials(roleARN, pARN, samlAssertion, conf.Provider.SessionDuration)
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
