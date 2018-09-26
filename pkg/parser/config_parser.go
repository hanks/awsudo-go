package parser

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/BurntSushi/toml"
	"github.com/hanks/awsudo-go/utils"
)

type provider struct {
	IDP                string `toml:"idp"`
	IDP_LOGIN_URL      string `toml:"idp_login_url"`
	SAML_PROVIDER_NAME string `toml:"saml_provide_name"`
	AUTH_API           string `toml:"auth_api"`
	SESSION_DURATION   int64  `toml:"session_duration"`
}

type role struct {
	Name string `toml:"name"`
	ARN  string `toml:"arn"`
}

type agent struct {
	Expiration int64 `toml:"expiration"`
}

type Config struct {
	Provider provider `toml:"provider"`
	Agent    agent    `toml:"agent"`
	Roles    []role   `toml:"roles"`
}

func LoadConfig(path string) (*Config, error) {
	var conf Config
	_, err := toml.DecodeFile(path, &conf)
	if err != nil {
		return nil, err
	}

	conf.Validate()

	return &conf, nil
}

func WriteConfig(path string, config *Config) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Create config (%s) error, please check: %v", path, err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	encoder := toml.NewEncoder(w)
	encoder.Indent = ""
	err = encoder.Encode(config)
	if err != nil {
		log.Fatalf("Write config (%s) error, please check: %v", path, err)
	}
}

func (c *Config) Validate() {
	if c.Agent.Expiration > c.Provider.SESSION_DURATION {
		log.Fatalf("Agent expiration (%d) should be smaller than session duration (%d), please check config file",
			c.Agent.Expiration,
			c.Provider.SESSION_DURATION,
		)
	}
}

func (c *Config) GetARN(name string) (bool, string) {
	existed := false
	arn := ""

	for _, role := range c.Roles {
		if role.Name == name {
			existed = true
			arn = role.ARN

			break
		}
	}

	return existed, arn
}

func (c *Config) GetPrincipalArn(arn string) string {
	r := regexp.MustCompile(`.*::(?P<AccountID>\d+):.*`)
	groups := r.FindStringSubmatch(arn)
	if len(groups) == 0 {
		return ""
	}

	return fmt.Sprintf("arn:aws:iam::%s:saml-provider/%s", groups[1], c.Provider.SAML_PROVIDER_NAME)
}

func (c *Config) InputConfig() {
	utils.InputString(&c.Provider.IDP, "IDP")
	utils.InputString(&c.Provider.IDP_LOGIN_URL, "IDP Login URL")
	utils.InputString(&c.Provider.SAML_PROVIDER_NAME, "SAML Provider Name")
	utils.InputString(&c.Provider.AUTH_API, "Auth API")

	utils.InputInt64(&c.Provider.SESSION_DURATION, "AWS Session Duration")
	utils.InputInt64(&c.Agent.Expiration, "Agent Expiration")

	c.Validate()

	// input existed roles
	for i := range c.Roles {
		utils.InputString(&c.Roles[i].Name, "AWS Role Name")
		utils.InputString(&c.Roles[i].ARN, "AWS Role ARN")
	}

	// input new roles, exits when both input are empty string
	for {
		exitCnt := 0
		r := role{}

		utils.InputString(&r.Name, "AWS Role Name")
		if r.Name == "" {
			exitCnt++
		}

		utils.InputString(&r.ARN, "AWS Role ARN")
		if r.ARN == "" {
			exitCnt++
		}

		if exitCnt == 2 {
			break
		}

		c.Roles = append(c.Roles, r)
	}
}
