package parser

import (
	"bufio"
	"fmt"
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

var newscanner = bufio.NewScanner

// LoadConfig is to create Config instance by loading from config file
func LoadConfig(path string) (*Config, error) {
	var conf Config
	_, err := toml.DecodeFile(path, &conf)
	if err != nil {
		return nil, err
	}

	err = conf.validate()
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

// WriteConfig is to write config struct data to config file
func (c *Config) WriteConfig(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("Create config (%s) error, please check: %v", path, err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	encoder := toml.NewEncoder(w)
	encoder.Indent = ""
	err = encoder.Encode(c)
	if err != nil {
		return fmt.Errorf("Write config (%s) error, please check: %v", path, err)
	}

	return nil
}

// validate is to do easy validation for config items
func (c *Config) validate() error {
	if c.Agent.Expiration >= c.Provider.SessionDuration {
		return fmt.Errorf("Agent expiration (%d) should be smaller than session duration (%d), please check config file",
			c.Agent.Expiration,
			c.Provider.SessionDuration,
		)
	}

	return nil
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
func (c *Config) InputConfig() error {
	scanner := newscanner(os.Stdin)

	c.Provider.IDP = utils.InputString(scanner, c.Provider.IDP, "IDP")
	c.Provider.IdpLoginURL = utils.InputString(scanner, c.Provider.IdpLoginURL, "IDP Login URL")
	c.Provider.SamlProviderName = utils.InputString(scanner, c.Provider.SamlProviderName, "SAML Provider Name")
	c.Provider.AuthAPI = utils.InputString(scanner, c.Provider.AuthAPI, "Auth API")

	c.Provider.SessionDuration, _ = utils.InputInt64(scanner, c.Provider.SessionDuration, "AWS Session Duration")
	c.Agent.Expiration, _ = utils.InputInt64(scanner, c.Agent.Expiration, "Agent Expiration")

	err := c.validate()
	if err != nil {
		return err
	}

	// input existed roles
	for i := range c.Roles {
		c.Roles[i].Name = utils.InputString(scanner, c.Roles[i].Name, "AWS Role Name")
		c.Roles[i].ARN = utils.InputString(scanner, c.Roles[i].ARN, "AWS Role ARN")
	}

	// input new roles, exits when both input are empty string
	for {
		exitCnt := 0
		r := role{}

		r.Name = utils.InputString(scanner, r.Name, "AWS Role Name")
		if r.Name == "" {
			exitCnt++
		}

		r.ARN = utils.InputString(scanner, r.ARN, "AWS Role ARN")
		if r.ARN == "" {
			exitCnt++
		}

		if exitCnt == 2 {
			break
		}

		c.Roles = append(c.Roles, r)
	}

	return nil
}
