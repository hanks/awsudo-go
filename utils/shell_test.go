package utils

import "testing"

type mockExec struct {
}

func (f *mockExec) execv(argv0 string, argv []string, envv []string) error {
	return nil
}

func TestExecCommand(t *testing.T) {
	type args struct {
		cmds  []string
		execv func(argv0 string, argv []string, envv []string) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Should return error when to execute invalid command",
			args{
				[]string{"notfound", "-alh", "."},
				new(mockExec).execv,
			},
			true,
		},
		{
			"Should execute valid command correctly",
			args{
				[]string{"ls", "-alh", "."},
				new(mockExec).execv,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execv = tt.args.execv
			if err := ExecCommand(tt.args.cmds); (err != nil) != tt.wantErr {
				t.Errorf("ExecCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
