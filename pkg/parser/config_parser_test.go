package parser

import (
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			"Should load config successfully with valid config",
			args{"../../test-fixtures/valid-config.toml"},
			&Config{
				provider{
					IDP:              "okta",
					IdpLoginURL:      "idp_login_url",
					SamlProviderName: "saml_provide_name",
					AuthAPI:          "auth_api",
					SessionDuration:  3600,
				},
				agent{3300},
				[]role{
					{
						"adhoc",
						"adhoc",
					},
				},
			},
			false,
		},
		{
			"Should load config failure with validation error",
			args{"../../test-fixtures/validation-error-config.toml"},
			nil,
			true,
		},
		{
			"Should load config failure with invalid toml format",
			args{"../../test-fixtures/invalid-config.toml"},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadConfig(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_GetARN(t *testing.T) {
	type fields struct {
		Provider provider
		Agent    agent
		Roles    []role
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantARN string
	}{
		{
			"Should return existed arn with true",
			fields{
				provider{},
				agent{},
				[]role{
					{
						"name1",
						"arn1",
					},
					{
						"name2",
						"arn2",
					},
					{
						"name3",
						"arn3",
					},
				},
			},
			args{
				"name2",
			},
			true,
			"arn2",
		},
		{
			"Should return empty arn with false",
			fields{
				provider{},
				agent{},
				[]role{
					{
						"name1",
						"arn1",
					},
				},
			},
			args{
				"name2",
			},
			false,
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Provider: tt.fields.Provider,
				Agent:    tt.fields.Agent,
				Roles:    tt.fields.Roles,
			}
			got, gotARN := c.GetARN(tt.args.name)
			if got != tt.want {
				t.Errorf("Config.GetARN() got = %v, want %v", got, tt.want)
			}
			if gotARN != tt.wantARN {
				t.Errorf("Config.GetARN() got1 = %v, want %v", gotARN, tt.wantARN)
			}
		})
	}
}

func TestConfig_GetPrincipalArn(t *testing.T) {
	type fields struct {
		Provider provider
		Agent    agent
		Roles    []role
	}
	type args struct {
		arn string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"Should return the correct principal arn string",
			fields{
				provider{SamlProviderName: "saml_provide_name"},
				agent{},
				[]role{},
			},
			args{
				"arn:aws:iam::12345678910:role/test",
			},
			"arn:aws:iam::12345678910:saml-provider/saml_provide_name",
		},
		{
			"Should return empty string when not found",
			fields{
				provider{SamlProviderName: "saml_provide_name"},
				agent{},
				[]role{},
			},
			args{
				"arn:aws:iam12345678910:role/test",
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Provider: tt.fields.Provider,
				Agent:    tt.fields.Agent,
				Roles:    tt.fields.Roles,
			}
			if got := c.GetPrincipalArn(tt.args.arn); got != tt.want {
				t.Errorf("Config.GetPrincipalArn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_validate(t *testing.T) {
	type fields struct {
		Provider provider
		Agent    agent
		Roles    []role
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"Should pass expiration field validation without no error",
			fields{
				provider{SessionDuration: 3600},
				agent{3300},
				[]role{},
			},
			false,
		},
		{
			"Should fail expiration field validation with error",
			fields{
				provider{SessionDuration: 3300},
				agent{3300},
				[]role{},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Provider: tt.fields.Provider,
				Agent:    tt.fields.Agent,
				Roles:    tt.fields.Roles,
			}
			if err := c.validate(); (err != nil) != tt.wantErr {
				t.Errorf("Config.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
