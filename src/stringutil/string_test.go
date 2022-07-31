package stringutil

import (
	"testing"
)

// TestRandomLengthString test random length string
func TestRandomLengthString(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"0-empty", args{0}, ""},
		{"1-normal", args{3}, "abc"},
		{"1-normal", args{1}, "a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RandomLengthString(tt.args.length)
			if len(got) != len(tt.want) {
				t.Errorf("RandomLengthString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetUUID(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"normal"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUUID(); got == "" {
				t.Errorf("getUUID() = %v", got)
			}
		})
	}
}
