package fs

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/bmizerany/assert"
)

var (
	sampleBuffer = []byte("helloword")
)

func TestGetFileSizeNoErr(t *testing.T) {
	tmpFile := path.Join(os.TempDir(), "simple")
	if err := ioutil.WriteFile(tmpFile, sampleBuffer, os.ModePerm); err != nil {
		return
	}
	defer os.Remove(tmpFile)

	notExistFile := path.Join(os.TempDir(), "13017ad2-a411-11ec-b909-0242ac120002")
	if IsFile(notExistFile) {
		os.Remove(notExistFile)
	}

	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{"normal", args{tmpFile}, 9},
		{"not_exist", args{notExistFile}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFileSizeNoErr(tt.args.filename); got != tt.want {
				t.Errorf("GetFileSizeNoErr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadFileWithSize(t *testing.T) {
	tmpFile := path.Join(os.TempDir(), "simple")
	if err := ioutil.WriteFile(tmpFile, sampleBuffer, os.ModePerm); err != nil {
		return
	}
	defer os.Remove(tmpFile)

	notExistFile := path.Join(os.TempDir(), "13017ad2-a411-11ec-b909-0242ac120002")
	if IsFile(notExistFile) {
		os.Remove(notExistFile)
	}
	halfSampleBuffer := len(sampleBuffer) >> 1

	type args struct {
		filename string
		pos      int64
		size     int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    []byte
	}{
		{"not_exist", args{notExistFile, 0, 0}, true, nil},
		{"invarg_seed", args{tmpFile, -1, 20}, true, nil},
		{"failed_to_read", args{tmpFile, 100, 20}, true, nil},
		{"read_lt_want", args{tmpFile, int64(halfSampleBuffer), len(sampleBuffer)}, true, nil},
		{"normal", args{tmpFile, 0, len(sampleBuffer)}, false, sampleBuffer},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadFileWithSize(tt.args.filename, tt.args.pos, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFileWithSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
