package utils

import (
	"testing"
)

type mockScanner struct {
	text string
}

func (m *mockScanner) Scan() bool {
	return true
}

func (m *mockScanner) Text() string {
	return m.text
}

type mockReadPass struct {
	pass []byte
	err  error
}

func (m *mockReadPass) ReadPassword(int) ([]byte, error) {
	return m.pass, m.err
}

func TestAskUserInput(t *testing.T) {
	type args struct {
		scanner  iScanner
		readPass *mockReadPass
	}
	tests := []struct {
		name     string
		args     args
		wantUser string
		wantPass string
	}{
		{
			"Should return the same values when input is valid",
			args{
				&mockScanner{
					"hello",
				},
				&mockReadPass{
					[]byte("world"),
					nil,
				},
			},
			"hello",
			"world",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// patch
			readPassword = tt.args.readPass.ReadPassword
			user, pass := AskUserInput(tt.args.scanner)
			if user != tt.wantUser {
				t.Errorf("AskUserInput() user = %v, want %v", user, tt.wantUser)
			}
			if pass != tt.wantPass {
				t.Errorf("AskUserInput() pass = %v, want %v", pass, tt.wantPass)
			}
		})
	}
}

func TestInputString(t *testing.T) {
	type args struct {
		scanner  iScanner
		original string
		name     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Should return original string when no input and empty original",
			args{
				&mockScanner{
					text: "",
				},
				"",
				"IDP",
			},
			"",
		},
		{
			"Should return original string when no input and non empty original",
			args{
				&mockScanner{
					text: "",
				},
				"original",
				"IDP",
			},
			"original",
		},
		{
			"Should return inputted string",
			args{
				&mockScanner{
					text: "adhoc",
				},
				"original",
				"IDP",
			},
			"adhoc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InputString(tt.args.scanner, tt.args.original, tt.args.name); got != tt.want {
				t.Errorf("InputString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInputInt64(t *testing.T) {
	type args struct {
		scanner  iScanner
		original int64
		name     string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			"Should return default expiration when no input",
			args{
				&mockScanner{
					text: "",
				},
				0,
				"Expiration",
			},
			3300,
			false,
		},
		{
			"Should return default duration when no input",
			args{
				&mockScanner{
					text: "",
				},
				0,
				"Duration",
			},
			3600,
			false,
		},
		{
			"Should return inputted valid value",
			args{
				&mockScanner{
					text: "100",
				},
				0,
				"Duration",
			},
			100,
			false,
		},
		{
			"Should return error when to input non number string",
			args{
				&mockScanner{
					text: "number",
				},
				0,
				"Duration",
			},
			0,
			true,
		},
		{
			"Should return original value when no new input",
			args{
				&mockScanner{
					text: "",
				},
				100,
				"Duration",
			},
			100,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InputInt64(tt.args.scanner, tt.args.original, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("InputInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("InputInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}
