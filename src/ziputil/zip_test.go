package ziputil

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestZip(t *testing.T) {
	// zip-xxx
	//    1.txt
	//    in-xxx
	//       2.txt
	tempDir, err := ioutil.TempDir("", "zip-")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer os.RemoveAll(tempDir)

	if err = ioutil.WriteFile(path.Join(tempDir, "1.txt"), []byte{1, 2, 3, 4, 5}, os.ModePerm); err != nil {
		fmt.Println(err)
		return
	}
	tempDirIn, err := ioutil.TempDir(tempDir, "in-")
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = ioutil.WriteFile(path.Join(tempDirIn, "2.txt"), []byte{5, 2, 3, 4, 5}, os.ModePerm); err != nil {
		fmt.Println(err)
		return
	}

	zipFileName := path.Join(tempDir, "test.zip")

	type args struct {
		srcDir      string
		zipFilename string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"failed_to_create_zip_file", args{"/test/to/not_existed.txt", "/test/not_existed.zip"}, true},
		{"failed_to_walk", args{"/test/to/not_existed.txt", zipFileName}, true},
		{"normal", args{tempDir, zipFileName}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Zip(tt.args.srcDir, tt.args.zipFilename); (err != nil) != tt.wantErr {
				t.Errorf("Zip() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
