package fs

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func TestGetParentAbsDir(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"dir_not_slash", args{"/tmp/dir1"}, "/tmp", false},
		{"dri_have_slash", args{"/tmp/dir1/"}, "/tmp/dir1", false},
		{"file", args{"/tmp/dir1/file.txt"}, "/tmp/dir1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetParentAbsDir(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetParentAbsDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetParentAbsDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsDir(t *testing.T) {

	tmpFile, err := ioutil.TempFile(os.TempDir(), "simple")
	if err != nil {
		return
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte("helloword"))
	if err != nil {
		return
	}
	tmpFile.Close()

	type args struct {
		file string
	}
	testsDir := []struct {
		name string
		args args
		want bool
	}{
		{"file", args{tmpFile.Name()}, false},
		{"not_exist", args{"/tmp/notexist"}, false},
		{"dir", args{os.TempDir()}, true},
		{".", args{"."}, true},
	}
	for _, tt := range testsDir {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDir(tt.args.file); got != tt.want {
				t.Errorf("[%v] IsDir() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestIsFile(t *testing.T) {

	tmpFile, err := ioutil.TempFile(os.TempDir(), "simple")
	if err != nil {
		return
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte("helloword"))
	if err != nil {
		return
	}
	tmpFile.Close()

	type args struct {
		file string
	}
	testsFile := []struct {
		name string
		args args
		want bool
	}{
		{"file", args{tmpFile.Name()}, true},
		{"not_exist", args{"/tmp/notexist"}, false},
		{"dir", args{os.TempDir()}, false},
	}
	for _, tt := range testsFile {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFile(tt.args.file); got != tt.want {
				t.Errorf("IsDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFileName(t *testing.T) {
	type args struct {
		inputPath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"normal", args{"/usr/local/test.txt"}, "test"},
		{"normal2", args{"/usr/local/test"}, "test"},
		{"dir", args{"/tmp"}, "tmp"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFileName(tt.args.inputPath); got != tt.want {
				t.Errorf("GetFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocalCopy(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "simple")
	if err != nil {
		return
	}
	defer os.Remove(tmpFile.Name())
	_, err = tmpFile.Write([]byte("helloword"))
	if err != nil {
		return
	}
	tmpFile.Close()
	dirName, _ := ioutil.TempDir(os.TempDir(), "")

	type args struct {
		src string
		dst string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{"same", args{tmpFile.Name(), tmpFile.Name()}, 9, false},
		{"normal", args{tmpFile.Name(), tmpFile.Name() + ".bak"}, 9, false},
		{"file_to_dir", args{tmpFile.Name(), dirName}, 0, true},
		{"dir_to_file", args{os.TempDir(), tmpFile.Name() + ".bak"}, 0, true},
		{"src_not_exist", args{"not exist", ""}, 0, true},
		{"dst_not_exist", args{tmpFile.Name(), ""}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LocalCopy(tt.args.src, tt.args.dst)
			if (err != nil) != tt.wantErr {
				t.Errorf("LocalCopy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("%v, LocalCopy() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestGetNameWithNewExt(t *testing.T) {
	type args struct {
		u   string
		ext string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"same", args{"/temp/1.mp3", ".mp3"}, "/temp/1.mp3"},
		{"notSame", args{"/temp/1.mp3", ".mp4"}, "/temp/1.mp4"},
		{"emptyExt", args{"/temp/abc", ".mp4"}, "/temp/abc.mp4"},
		{"empty", args{"", ".mp4"}, ".mp4"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNameWithNewExt(tt.args.u, tt.args.ext); got != tt.want {
				t.Errorf("ReplaceExt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddNonceParam(t *testing.T) {
	testTime := time.Date(1974, time.May, 19, 1, 2, 3, 4, time.UTC)

	type args struct {
		u string
		t time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"empty_url", args{"", testTime}, "", true},
		{"normal", args{"http://www.tencent.com", testTime}, "http://www.tencent.com?noncestr=", false},
		{"illegal_url", args{"http://www.ten cent.com", testTime}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddNonceParam(tt.args.u, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddNonceParam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.HasPrefix(got, tt.want) {
				t.Errorf("AddNonceParam() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileNameAppend(t *testing.T) {
	type args struct {
		inputPath  string
		appText    string
		defaultExt string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{"", "", ""}, ""},
		{"empty", args{"", "_1", ""}, "_1"},
		{"normal_path", args{"/usr/local/1.txt", "_1", ""}, "/usr/local/1_1.txt"},
		{"normal_file", args{"1.txt", "_1", ""}, "1_1.txt"},
		{"normal_relate", args{"test/1.txt", "_1", ".txt2"}, "test/1_1.txt"},
		{"default_text", args{"test/1", "_1", ".txt2"}, "test/1_1.txt2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FileNameAppend(tt.args.inputPath, tt.args.appText, tt.args.defaultExt); got != tt.want {
				t.Errorf("%v FileNameAppend() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{"empty", "", ""},
		{"local_path", "/qq/0.txt", "/qq/0.txt"},
		{"relative_url", "qq.com/1.txt", "http://qq.com/1.txt"},
		{"no_scheme", "//qq.com/2.txt", "http://qq.com/2.txt"},
		{"http_scheme", "http://qq.com/3.txt", "http://qq.com/3.txt"},
		{"https_scheme", "https://qq.com/4.txt", "https://qq.com/4.txt"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeURL(tt.args); got != tt.want {
				t.Errorf("%v NormalizeURL() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
