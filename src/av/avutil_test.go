package av

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
	"time"

	"errors"
	"go-utils/src/tools/fs"
	"go-utils/src/logs"
	
	"github.com/agiledragon/gomonkey/v2"
)

func TestCutMedia(t *testing.T) {
	videoFile := getTestAsset()
	if videoFile == "" {
		t.Fatal("getTestAsset fail")
	}
	notExistPath := path.Join(os.TempDir(), "f313df34-32e1-4a4c-9472-0c1cfc8abf0b")

	type args struct {
		ctx        context.Context
		inputPath  string
		outputPath string
		start      time.Duration
		dur        time.Duration
	}
	outFile := path.Join(os.TempDir(), "30c0a124-a59c-11ec-b909-0242ac120002.mp4")
	outFile2 := path.Join(os.TempDir(), "318a46da-a5bc-11ec-b909-0242ac120002.mp4")

	ctx := context.Background()
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"normal",
			args{ctx, videoFile, outFile, 3 * time.Second, 5 * time.Second},
			path.Join(os.TempDir(), "30c0a124-a59c-11ec-b909-0242ac120002.webm"),
			false,
		},
		{
			"normal_over_end",
			args{ctx, videoFile, outFile2, 8 * time.Second, 5 * time.Second},
			path.Join(os.TempDir(), "318a46da-a5bc-11ec-b909-0242ac120002.webm"),
			false,
		},
		{"empty", args{ctx, "", outFile, 0, time.Second}, "", true},
		{"notexist", args{ctx, notExistPath, outFile, 0, time.Second}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CutMedia(tt.args.ctx, tt.args.inputPath, tt.args.outputPath, tt.args.start, tt.args.dur)
			if (err != nil) != tt.wantErr {
				t.Errorf("CutMedia() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CutMedia() = %v, want %v", got, tt.want)
			} else {
				if fs.IsFile(got) {
					os.Remove(got)
				}
			}
		})
	}
}

func TestMediaReverse(t *testing.T) {
	videoFile := getTestAsset()
	if videoFile == "" {
		t.Fatal("getTestAsset fail")
	}

	outFile := path.Join(os.TempDir(), "out.mp4")
	type args struct {
		ctx        context.Context
		inputPath  string
		outputPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"normal", args{context.Background(), videoFile, outFile}, false},
		{"empty", args{context.Background(), "", outFile}, true},
		{"notexist", args{context.Background(), "notexist", outFile}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MediaReverse(tt.args.ctx, tt.args.inputPath, tt.args.outputPath); (err != nil) != tt.wantErr {
				t.Errorf("MediaReverse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMergeTS(t *testing.T) {
	type args struct {
		ctx        context.Context
		inputPath  string
		outputPath string
	}
	outFile := path.Join(os.TempDir(), "out.mp4")
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{context.Background(), "", outFile}, true},
		{"notexist", args{context.Background(), "notexist", outFile}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MergeTS(tt.args.ctx, tt.args.inputPath, tt.args.outputPath); (err != nil) != tt.wantErr {
				t.Errorf("MergeTS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConvertToWav(t *testing.T) {
	type args struct {
		ctx        context.Context
		inputPath  string
		start      time.Duration
		dur        time.Duration
		outputPath string
	}
	outFile := path.Join(os.TempDir(), "out.mp4")
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{context.Background(), "", 0, 0, outFile}, true},
		{"notexist", args{context.Background(), "notexist", 0, 0, outFile}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ConvertToWav(
				tt.args.ctx,
				tt.args.inputPath,
				tt.args.start,
				tt.args.dur,
				tt.args.outputPath,
			); (err != nil) != tt.wantErr {
				t.Errorf("ConvertToWav() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSoundstretchProcess(t *testing.T) {
	type args struct {
		ctx        context.Context
		inputPath  string
		outputPath string
		pitch      float32
	}
	outFile := path.Join(os.TempDir(), "out.mp4")
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{context.Background(), "", outFile, 1}, true},
		{"notexist", args{context.Background(), "notexist", outFile, 1}, true},
		{"para fail", args{context.Background(), "", outFile, 80}, true},
		{"para fail", args{context.Background(), "notexist.wav", "output.wav", 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SoundstretchProcess(
				tt.args.ctx,
				tt.args.inputPath,
				tt.args.outputPath,
				tt.args.pitch,
			); (err != nil) != tt.wantErr {
				t.Errorf("SoundstretchProcess() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFormatConvert(t *testing.T) {
	videoFile := getTestAsset()
	if videoFile == "" {
		t.Fatal("getTestAsset fail")
	}

	type args struct {
		ctx       context.Context
		inputPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"normal", args{context.Background(), videoFile}, false},
		{"empty", args{context.Background(), ""}, true},
		{"notexist", args{context.Background(), "notexist"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatConvert(tt.args.ctx, tt.args.inputPath, Mp4Ext)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatConvert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) <= 0 {
				t.Errorf("FormatConvert() = %v", got)
			}
		})
	}
}

func TestResizeMedia(t *testing.T) {

	audioStream := Streams{
		CodecType: CodecTypeAudio,
	}
	smallSizeVideoStream := Streams{
		CodecType:   CodecTypeVideo,
		CodedWidth:  640,
		CodedHeight: 360,
	}
	// landscape video stream
	landscapeVideoStream := Streams{
		CodecType:   CodecTypeVideo,
		CodedWidth:  2560,
		CodedHeight: 1000,
	}
	// portrait video stream
	portraitVideoStream := Streams{
		CodecType:   CodecTypeVideo,
		CodedWidth:  1000,
		CodedHeight: 2560,
	}

	patches := gomonkey.ApplyFuncSeq(Probe, []gomonkey.OutputCell{
		{Values: gomonkey.Params{nil, errors.ErrUnknown}, Times: 1},
		{Values: gomonkey.Params{&ProbeInfo{}, nil}, Times: 1},
		{
			Values: gomonkey.Params{
				&ProbeInfo{
					Streams: []Streams{
						// skip no video
						audioStream,
						// skip small resource video
						smallSizeVideoStream,
						// failed to run command
						landscapeVideoStream,
					},
				},
				nil,
			},
			Times: 1,
		},

		// landscape video
		{
			Values: gomonkey.Params{
				&ProbeInfo{
					Streams: []Streams{
						landscapeVideoStream,
					},
				},
				nil,
			},
			Times: 1,
		},

		// portrait video
		{
			Values: gomonkey.Params{
				&ProbeInfo{
					Streams: []Streams{
						portraitVideoStream,
					},
				},
				nil,
			},
			Times: 2,
		},
	})

	patches = patches.ApplyFuncSeq(fs.RunSysCommand, []gomonkey.OutputCell{
		// failed
		{Values: gomonkey.Params{errors.ErrUnknown}, Times: 1},
		// landscape + portrait + portrait default ext
		{Values: gomonkey.Params{nil}, Times: 3},
	})
	defer patches.Reset()

	canvasWidth := 1280
	canvasHeight := 720
	testMediaPath := "/test/to/path"
	// default with .mp4 extension
	testLandscapeOutputPath := "/test/to/path_1280x-2.mp4"
	testPortraitOutputPath := "/test/to/path_-2x720.mp4"
	testPortraitOutputPathWithExt := "/test/to/path_-2x720.png"
	tests := []struct {
		name       string
		defaultExt string
		wantErr    bool
		wantOutput string
	}{
		{"failed_to_probe", "", true, ""},
		{"empty_streams", "", false, testMediaPath},
		{"failed_to_resize", "", true, ""},
		{"landscape", "", false, testLandscapeOutputPath},
		{"portrait", "", false, testPortraitOutputPath},
		{"portrait_with_ext", ".png", false, testPortraitOutputPathWithExt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResizeMediaFitIn(context.Background(), testMediaPath, canvasWidth, canvasHeight, tt.defaultExt)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResizeMedia() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantOutput) {
				t.Errorf("got = %v, want %v", got, tt.wantOutput)
			}
		})
	}
}

func createTempMp4FilewithString(ctx context.Context, dir, content string) (string, error) {
	tmpFile, err := ioutil.TempFile(dir, "*.mp4")
	if err != nil {
		logs.Log.Errorf("failed to create temp file: %+v", err)
		return "", err
	}
	defer tmpFile.Close()
	if _, err = tmpFile.WriteString(content); err != nil {
		logs.Log.Errorf("failed to write string to file %s", tmpFile.Name())
		return "", err
	}
	return tmpFile.Name(), nil
}

func TestGetSuggesedExtFromContent(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "render-")
	if err != nil {
		logs.Log.Errorf("failed to create temp dir: %+v", err)
		return
	}
	defer os.RemoveAll(tmpDir)

	allFileHeader := []fileHeader{ID3v2, FLAC1, FLAC2}
	allExt := []string{Mp3Ext, Mp3Ext, Mp3Ext}

	tmpMp4File, err := createTempMp4FilewithString(context.Background(), tmpDir, "mp4xxx")
	if err != nil {
		logs.Log.Errorf("failed to create temp file: %+v", err)
		return
	}

	type test struct {
		name string
		args string
		want string
	}
	tests := []test{
		{"not_exists", "/test/111.mp4", Mp4Ext},
		{"common_mp4_ext", tmpMp4File, Mp4Ext},
	}

	for idx, fh := range allFileHeader {
		tmpID3v2File, err := createTempMp4FilewithString(context.Background(), tmpDir, string(fh)+"xxx")
		if err != nil {
			logs.Log.Errorf("failed to create temp file: %+v", err)
			return
		}
		tests = append(tests, test{
			string(fh),
			tmpID3v2File,
			allExt[idx],
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSuggestedExtFromContent(context.Background(), tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}
