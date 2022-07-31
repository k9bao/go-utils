package av

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"go-utils/src/tools/fs"
	"go-utils/src/tools/httputil"
	
)

const webmExt = ".webm"

func getTestAsset() string {
	videoFile, err := download("vp9_alpha.webm")
	if err != nil {
		return ""
	}
	sysType := runtime.GOOS
	if sysType != "darwin" { // mac权限问题，需提前自备
		ffmpegBin, err = download("ffmpeg")
		if err != nil {
			return ""
		}
		_ = os.Chmod(ffmpegBin, 0777)
		ffprobeBin, err = download("ffprobe")
		if err != nil {
			return ""
		}
	}

	_ = os.Chmod(ffprobeBin, 0777)
	return videoFile
}

func download(fileName string) (string, error) {
	ctx := context.Background()
	dstPath := filepath.Join(os.TempDir(), fileName)
	if fs.IsFile(dstPath) {
		return dstPath, nil
	}

	urlParent := "https://cpom-zenvideo-zenfile-1258344701.cos.ap-beijing.myqcloud.com" +
		"/hunchai/sample_videos/render_server_plus/"
	url := urlParent + fileName

	header := httputil.GetDefaultHeader()
	resp, _, err := httputil.TryCountGetRespRedirect(ctx, http.MethodGet, url, header, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if fs.IsDir(dstPath) {
		name := httputil.GetHTTPFileName(url, resp, Mp4Ext, "")
		dstPath = filepath.Join(dstPath, name)
	}

	fileHander, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer fileHander.Close()
	_, err = io.Copy(fileHander, resp.Body)
	if err != nil {
		logs.Log.Errorf("url=%v,dstPath=%v,err=%+v", url, dstPath, err)
		return "", err
	}
	logs.Log.Infof(
		"success download. [%v, %v]%v, %v",
		resp.ContentLength, fs.GetFileSizeNoErr(dstPath), url, dstPath,
	)

	return dstPath, nil
}

func TestProbeErr(t *testing.T) {
	type args struct {
		ctx       context.Context
		inputPath string
	}
	tests := []struct {
		name    string
		args    args
		want    *ProbeInfo
		wantErr bool
	}{
		{"nil", args{context.Background(), ""}, nil, true},
		{"not exist", args{context.Background(), "notexist.mp4"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Probe(tt.args.ctx, tt.args.inputPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Probe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Probe() = %v, want %v", got, tt.want)
			}
		})
	}
}

func assertTrue(t *testing.T, ok bool, desc string) {
	if !ok {
		t.Fatal(desc)
	}
}

func TestProbe(t *testing.T) {
	videoFile := getTestAsset()
	if videoFile == "" {
		t.Fatal("getTestAsset fail")
	}
	prob, err := Probe(context.Background(), videoFile)
	if err != nil {
		t.Fatal("probe fail")
	}
	assertTrue(t, prob.GetVideoCodec() == "vp9", fmt.Sprintf("GetVideoCodec fail %v", prob.GetVideoCodec()))
	assertTrue(t, prob.GetAudioCodec() == "opus", fmt.Sprintf("GetAudioCodec fail %v", prob.GetAudioCodec()))
	assertTrue(t, prob.GetFormatDuration() == 10, fmt.Sprintf("GetFormatDuration fail %v", prob.GetFormatDuration()))
	assertTrue(t, prob.getDefaultExt() == webmExt, fmt.Sprintf("getDefaultExt fail %v", prob.getDefaultExt()))
	assertTrue(
		t,
		prob.GetSuggestedExtFromCodec() == webmExt,
		fmt.Sprintf("GetSuggestedExtFromCodec fail %v", prob.GetSuggestedExtFromCodec()),
	)
	assertTrue(t, prob.HasAlpha(), fmt.Sprintf("HasAlpha fail %v", prob.HasAlpha()))
}

func TestProbeInfo_GetSuggestedExtFromCodec(t *testing.T) {
	type fields struct {
		InputPath string
		Streams   []Streams
		Format    Format
	}
	info, err := Probe(context.Background(), "/data/code/render_server_plus/ID3-ffprobe-err.mp4")
	if err != nil {
		return
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"normal", fields{info.InputPath, info.Streams, info.Format}, "mp3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ProbeInfo{
				InputPath: tt.fields.InputPath,
				Streams:   tt.fields.Streams,
				Format:    tt.fields.Format,
			}
			if got := p.GetSuggestedExtFromCodec(); got != tt.want {
				t.Errorf("ProbeInfo.GetSuggestedExtFromCodec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getExt(t *testing.T) {
	type args struct {
		slice      []formatName
		index      int
		defaultExt string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"normal", args{[]formatName{formatMp4, formatMp3}, 0, formatMp4.getExt()}, formatMp4.getExt()},
		{"err_out", args{nil, 1, formatMp4.getExt()}, formatMp4.getExt()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getExt(tt.args.slice, tt.args.index, tt.args.defaultExt); got != tt.want {
				t.Errorf("getExt() = %v, want %v", got, tt.want)
			}
		})
	}
}
