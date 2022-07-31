package algorithm

import (
	"testing"
)

func TestMaxInt64(t *testing.T) {
	type args struct { 
		a int64
		b int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{"a>b", args{300, 100}, 300},
		{"a=b", args{300, 300}, 300},
		{"a<b", args{100, 300}, 300},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxInt64(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("MaxInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMinInt64(t *testing.T) {
	type args struct {
		a int64
		b int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{"a>b", args{300, 100}, 100},
		{"a=b", args{300, 300}, 300},
		{"a<b", args{100, 300}, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MinInt64(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("MinInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxInt(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"a>b", args{300, 100}, 300},
		{"a=b", args{300, 300}, 300},
		{"a<b", args{100, 300}, 300},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxInt(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("MaxInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMinInt(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"a>b", args{300, 100}, 100},
		{"a=b", args{300, 300}, 300},
		{"a<b", args{100, 300}, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MinInt(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("MinInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseFloat(t *testing.T) {
	var defVal float64 = 10
	type args struct {
		s          string
		defaultVal float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"normal", args{"3.14", defVal}, 3.14},
		{"0", args{"0", defVal}, 0},
		{"empty", args{"", defVal}, defVal},
		{"3.14a", args{"3.14a", defVal}, defVal},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseFloat(tt.args.s, tt.args.defaultVal); got != tt.want {
				t.Errorf("[%v]ParseFloat() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestParseFloatDecimals(t *testing.T) {
	var defVal float64 = 10
	type args struct {
		s          string
		defaultVal float64
		decimals   int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"normal", args{"3.14", defVal, 1}, 3.1},
		{"0", args{"0", defVal, 1}, 0},
		{"empty", args{"", defVal, 1}, defVal},
		{"3.14a", args{"3.14a", defVal, 1}, defVal},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseFloatDecimals(tt.args.s, tt.args.defaultVal, tt.args.decimals); got != tt.want {
				t.Errorf("ParseFloatDecimals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseInt(t *testing.T) {
	var defVal int64 = 10
	type args struct {
		s          string
		defaultVal int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{"normal", args{"3", defVal}, 3},
		{"float", args{"3.14", defVal}, defVal},
		{"0", args{"0", defVal}, 0},
		{"empty", args{"", defVal}, defVal},
		{"3.14a", args{"3.14a", defVal}, defVal},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseInt(tt.args.s, tt.args.defaultVal); got != tt.want {
				t.Errorf("[%v]ParseInt() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
