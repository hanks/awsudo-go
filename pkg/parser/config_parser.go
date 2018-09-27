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
	IDP              string `toml:"idp"`
	IdpLoginURL      string `toml:"idp_login_url"`
	SamlProviderName string `toml:"saml_provide_name"`
	AuthAPI          string `toml:"auth_api"`
	SessionDuration  int64  `toml:"session_duration"`
}

type role struct {
	Name string `toml:"name"`
	ARN  string `toml:"arn"`
}

type agent struct {
	Expiration int64 `toml:"expiration"`
}

// Config is the struct for awsudo config file
type Config struct {
	Provider provider `toml:"provider"`
	Agent    agent    `toml:"agent"`
	Roles    []role   `toml:"roles"`
}

// LoadConfig is to create Config instance by loading from config file
func LoadConfig(path string) (*Config, error) {
	var conf Config
	_, err := toml.DecodeFile(path, &conf)
	if err != nil {
		return nil, err
	}

	conf.Validate()

	return &conf, nil
}

// WriteConfig is to write config struct data to config file
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

// Validate is to do easy validation for config items
func (c *Config) Validate() {
	if c.Agent.Expiration > c.Provider.SessionDuration {
		log.Fatalf("Agent expiration (%d) should be smaller than session duration (%d), please check config file",
			c.Agent.Expiration,
			c.Provider.SessionDuration,
		)
	}
}

// GetARN is to get arn from role name
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

// GetPrincipalArn is to get PrincipalArn from arn string
func (c *Config) GetPrincipalArn(arn string) string {
	r := regexp.MustCompile(`.*::(?P<AccountID>\d+):.*`)
	groups := r.FindStringSubmatch(arn)
	if len(groups) == 0 {
		return ""
	}

	return fmt.Sprintf("arn:aws:iam::%s:saml-provider/%s", groups[1], c.Provider.SamlProviderName)
}

// InputConfig is to accept user intput to set each items for config
func (c *Config) InputConfig() {
	utils.InputString(&c.Provider.IDP, "IDP")
	utils.InputString(&c.Provider.IdpLoginURL, "IDP Login URL")
	utils.InputString(&c.Provider.SamlProviderName, "SAML Provider Name")
	utils.InputString(&c.Provider.AuthAPI, "Auth API")

	utils.InputInt64(&c.Provider.SessionDuration, "AWS Session Duration")
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
