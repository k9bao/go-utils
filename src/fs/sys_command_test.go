package fs

import (
	"context"
	"os"
	"reflect"
	"testing"
)

func TestRunSysCommand(t *testing.T) {
	type args struct {
		ctx      context.Context
		commands []string
		envs     map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"normal", args{context.Background(), []string{"pwd"}, nil}, false},
		{"empty", args{context.Background(), []string{}, nil}, true},
		{"not exist", args{context.Background(), []string{"i am not an command"}, nil}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RunSysCommand(tt.args.ctx, tt.args.commands, tt.args.envs); (err != nil) != tt.wantErr {
				t.Errorf("RunSysCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRunSysCommandRet(t *testing.T) {
	type args struct {
		ctx      context.Context
		commands []string
		envs     map[string]string
	}
	dir, _ := os.Getwd()

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"normal", args{context.Background(), []string{"pwd"}, nil}, []byte(dir + "\n"), false},
		{"empty", args{context.Background(), []string{}, nil}, nil, true},
		{"not exist", args{context.Background(), []string{"i am not an command"}, nil}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RunSysCommandRet(tt.args.ctx, tt.args.commands, tt.args.envs)
			if (err != nil) != tt.wantErr {
				t.Errorf("RunSysCommandRet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%v, RunSysCommandRet() = %v, want %v", tt.name, string(got), string(tt.want))
			}
		})
	}
}
