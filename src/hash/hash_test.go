package hash

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestFMD5(t *testing.T) {
	outFile := path.Join(os.TempDir(), "filename")
	if err := ioutil.WriteFile(outFile, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}, os.ModePerm); err != nil {
		return
	}

	defer os.Remove(outFile)

	type args struct {
		localFilePath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"notexist", args{"notexist"}, "", true},
		{"normal", args{outFile}, "8596c1af55b14b7b320112944fcb8536", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FMD5(tt.args.localFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("[%v] FMD5() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("[%v] FMD5() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestFSHA1(t *testing.T) {
	outFile := path.Join(os.TempDir(), "filename")
	if err := ioutil.WriteFile(outFile, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}, os.ModePerm); err != nil {
		return
	}
	defer os.Remove(outFile)

	type args struct {
		localFilePath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"notexist", args{"notexist"}, "", true},
		{"normal", args{outFile}, "b6c511873b07a73513161b142d344b7b845cacef", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FSHA1(tt.args.localFilePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("FSHA1() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FSHA1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMD5(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"empty", args{""}, "d41d8cd98f00b204e9800998ecf8427e", false},
		{"normal", args{"123456789"}, "25f9e794323b453885f5181f1b624d0b", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MD5(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("MD5() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MD5() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMd5String(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{""}, "d41d8cd98f00b204e9800998ecf8427e"},
		{"normal", args{"123456789"}, "25f9e794323b453885f5181f1b624d0b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Md5String(tt.args.s); got != tt.want {
				t.Errorf("Md5String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSHA1(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"empty", args{""}, "da39a3ee5e6b4b0d3255bfef95601890afd80709", false},
		{"normal", args{"123456789"}, "f7c3bc1d808e04732adf679965ccc34ca7ae3441", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SHA1(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("SHA1() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("[%v] SHA1() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
