package utils

import "testing"
import "github.com/stretchr/testify/mock"


type MyMockedObject struct{
	mock.Mock
  }
func TestGetAbsPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAbsPath(tt.args.path); got != tt.want {
				t.Errorf("GetAbsPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
