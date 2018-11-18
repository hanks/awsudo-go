package utils

import "testing"

func TestGetAbsPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"Should return empty path when input is empty path",
			args{""},
			"",
			false,
		},
		{
			"Should expand home relative path correctly",
			args{"~/Desktop"},
			"/root/Desktop",
			false,
		},
		{
			"Should return error when to expand invalid home relative path",
			args{"~Desktop"},
			"",
			true,
		},
		{
			"Should return abs path when to input non home relative path",
			args{"Desktop"},
			"/go/src/github.com/hanks/awsudo-go/utils/Desktop",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAbsPath(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAbsPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetAbsPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
