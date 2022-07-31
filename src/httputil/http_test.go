package httputil

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"reflect"
	"testing"

	"errors"
	"go-utils/src/config"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/google/uuid"
)

func Test_GetDefaultHeader(t *testing.T) {
	tests := []struct {
		name string
		want http.Header
	}{
		{"first", GetDefaultHeader()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDefaultHeader(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDefaultHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getNewRequest(t *testing.T) {
	testCount := 4
	patchesNewRequest := gomonkey.ApplyFuncSeq(http.NewRequestWithContext, []gomonkey.OutputCell{
		{Values: gomonkey.Params{&http.Request{}, nil}, Times: testCount - 1},
		{Values: gomonkey.Params{nil, errors.New("NewRequestWithContext fail")}, Times: 1},
	})
	defer patchesNewRequest.Reset()

	ctx := context.Background()

	type args struct {
		ctx    context.Context
		method string
		url    string
		header http.Header
		data   *bytes.Buffer
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Request
		wantErr bool
	}{
		{"normal_no_data", args{ctx, "", "", nil, nil}, &http.Request{}, false},
		{"normal_have_data", args{ctx, "", "", nil, &bytes.Buffer{}}, &http.Request{}, false},
		{"normal_with_header",
			args{ctx, "", "", GetDefaultHeader(), nil},
			&http.Request{Header: GetDefaultHeader()},
			false,
		},
		{"error_new_request", args{ctx, "", "", nil, nil}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getNewRequest(tt.args.ctx, tt.args.method, tt.args.url, tt.args.header, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("getNewRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[%v] getNewRequest() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_getResp(t *testing.T) {
	patchesGetNewRequest := gomonkey.ApplyFunc(getNewRequest, func(
		ctx context.Context,
		method, url string,
		header http.Header,
		data *bytes.Buffer,
	) (*http.Request, error) {
		return &http.Request{Host: url}, nil
	})
	defer patchesGetNewRequest.Reset()

	ctx := context.Background()

	type args struct {
		ctx    context.Context
		method string
		url    string
		header http.Header
		data   []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr bool
	}{
		{"err_do", args{ctx, "", "", nil, nil}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getResp(tt.args.ctx, tt.args.method, tt.args.url, tt.args.header, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("[%v] getResp() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getResp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHTTPFileName(t *testing.T) {
	header := http.Header{}
	header.Set("Content-Disposition", "text/html; charset=utf-8; filename=hello.ts")
	type args struct {
		uri        string
		resp       *http.Response
		defaultExt string
		pre        string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"normal_uri_empty", args{"", &http.Response{}, ".mp4", "pre"}, "pre_.mp4"},
		{"normal", args{"http://tencent.com/test.mp3", &http.Response{}, ".mp4", "pre"}, "pre_test.mp3"},
		{"normal_no_ext", args{"http://tencent.com/test", &http.Response{}, ".mp4", "pre"}, "pre_test"},
		{"normal_resp", args{"", &http.Response{Header: header}, ".mp4", "pre"}, "pre_hello.ts"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHTTPFileName(tt.args.uri, tt.args.resp, tt.args.defaultExt, tt.args.pre); got != tt.want {
				t.Errorf("GetHTTPFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_find(t *testing.T) {
	type args struct {
		slice []int
		val   int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"normal_find", args{[]int{1, 2, 3}, 1}, true},
		{"normal_not_find", args{[]int{1, 2, 3}, 4}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := find(tt.args.slice, tt.args.val); got != tt.want {
				t.Errorf("find() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Read(p []byte) (n int, err error)
// Write(p []byte) (n int, err error)
// Close() error
// io.ReadCloser
type readWriteCloserImpl struct {
}

func (rc *readWriteCloserImpl) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (rc *readWriteCloserImpl) Write(p []byte) (n int, err error) {
	return 0, nil
}
func (rc *readWriteCloserImpl) Close() error {
	return nil
}

func Test_getRespRedirect(t *testing.T) {
	cfg := config.ServerConfig{}
	cfg.ControlConfig.MaxRedirectCounts = 2
	patchesRedirectCounts := gomonkey.ApplyGlobalVar(&config.ServerCnf, cfg)
	defer patchesRedirectCounts.Reset()

	header := http.Header{}
	header.Add("Location", "")

	bytes.NewBuffer([]byte{1, 2, 3})

	patchesGetResp := gomonkey.ApplyFuncSeq(getResp, []gomonkey.OutputCell{
		{Values: gomonkey.Params{nil, errors.New("getResp fail")}, Times: 1}, // getResp Fail
		{Values: gomonkey.Params{&http.Response{StatusCode: 301, Header: header, Body: new(readWriteCloserImpl)}, nil},
			Times: 1,
		}, // redirectOK
		{Values: gomonkey.Params{&http.Response{StatusCode: 200}, nil}, Times: 1}, // redirectOK
		{Values: gomonkey.Params{&http.Response{StatusCode: 404, Body: new(readWriteCloserImpl)}, nil},
			Times: 1,
		}, // 404
	})
	defer patchesGetResp.Reset()

	ctx := context.Background()

	type args struct {
		ctx    context.Context
		method string
		url    string
		header http.Header
		data   []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Response
		want1   string
		wantErr bool
	}{
		{"getResp Fail", args{ctx, "", "", nil, nil}, nil, "", true},
		{"redirectOK", args{ctx, "", "", nil, nil}, &http.Response{StatusCode: 200}, "", false},
		{"404", args{ctx, "", "", nil, nil}, nil, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := getRespRedirect(tt.args.ctx, tt.args.method, tt.args.url, tt.args.header, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("[%v] getRespRedirect() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[%v] getRespRedirect() got = %v, want %v", tt.name, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("[%v] getRespRedirect() got1 = %v, want %v", tt.name, got1, tt.want1)
			}
		})
	}
}

func Test_errorFromResponse(t *testing.T) {
	type args struct {
		resp *http.Response
		err  error
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"nil", args{nil, nil}, false},
		{"render_error", args{nil, errors.errorserverNoService}, true},
		{"system_error", args{&http.Response{StatusCode: 404, Body: new(readWriteCloserImpl)}, io.EOF}, true},
		{"system_resp_body_close_error", args{&http.Response{StatusCode: 404, Body: new(readWriteCloserImpl)}, io.EOF}, true},
		{"unknow_error", args{nil, errors.New("")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := errorFromResponse(tt.args.resp, tt.args.err); (err != nil) != tt.wantErr {
				t.Errorf("[%v] errorFromResponse() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

func TestTryCountGetRespRedirect(t *testing.T) {
	cfg := config.ServerConfig{}
	cfg.ControlConfig.HTTPRequestRetryCounts = 2
	patchesRedirectCounts := gomonkey.ApplyGlobalVar(&config.ServerCnf, cfg)
	defer patchesRedirectCounts.Reset()

	patchesGetRespRedirect := gomonkey.ApplyFuncSeq(getRespRedirect, []gomonkey.OutputCell{
		{Values: gomonkey.Params{nil, "", errors.New("getRespRedirect fail")}, Times: 1}, // normal
		{Values: gomonkey.Params{&http.Response{}, "", nil}, Times: 1},                   // normal
	})
	defer patchesGetRespRedirect.Reset()
	ctx := context.Background()

	type args struct {
		ctx    context.Context
		method string
		url    string
		header http.Header
		data   []byte
	}
	tests := []struct {
		name            string
		args            args
		wantResp        *http.Response
		wantRedirectURL string
		wantErr         bool
	}{
		{"normal", args{ctx, "", "", nil, nil}, &http.Response{}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, gotRedirectURL, err := TryCountGetRespRedirect(
				tt.args.ctx, tt.args.method, tt.args.url, tt.args.header, tt.args.data,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("TryCountGetRespRedirect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("TryCountGetRespRedirect() gotResp = %v, want %v", gotResp, tt.wantResp)
			}
			if gotRedirectURL != tt.wantRedirectURL {
				t.Errorf("TryCountGetRespRedirect() gotRedirectURL = %v, want %v", gotRedirectURL, tt.wantRedirectURL)
			}
		})
	}
}

func TestAcceptRange(t *testing.T) {

	header := http.Header{}
	header.Add("Accept-Ranges", "bytes")
	header.Add("Content-Range", "3600-5000/5000")
	patchesTryCountGetRespRedirect := gomonkey.ApplyFuncSeq(TryCountGetRespRedirect, []gomonkey.OutputCell{
		{Values: gomonkey.Params{nil, "", errors.New("getRespRedirect fail")}, Times: 1},
		{Values: gomonkey.Params{&http.Response{StatusCode: 200}, "", nil}, Times: 1},
		{
			Values: gomonkey.Params{
				&http.Response{Header: header, StatusCode: 206, Body: new(readWriteCloserImpl)}, "", nil,
			},
			Times: 1,
		},
	})
	defer patchesTryCountGetRespRedirect.Reset()

	patchesGetHTTPFileName := gomonkey.ApplyFuncSeq(GetHTTPFileName, []gomonkey.OutputCell{
		{Values: gomonkey.Params{""}, Times: 3},
	})
	defer patchesGetHTTPFileName.Reset()

	ctx := context.Background()

	type args struct {
		ctx context.Context
		url string
	}
	tests := []struct {
		name            string
		args            args
		wantFileSize    int64
		wantFileName    string
		wantRedirectURL string
		wantErr         bool
	}{
		{"error_TryCountGetRespRedirect", args{ctx, ""}, 0, "", "", true},
		{"normal_not_support", args{ctx, ""}, 0, "", "", true},
		{"normal_support", args{ctx, ""}, 5000, "", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFileSize, gotFileName, gotRedirectURL, err := AcceptRange(tt.args.ctx, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("AcceptRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotFileSize != tt.wantFileSize {
				t.Errorf("AcceptRange() gotFileSize = %v, want %v", gotFileSize, tt.wantFileSize)
			}
			if gotFileName != tt.wantFileName {
				t.Errorf("AcceptRange() gotFileName = %v, want %v", gotFileName, tt.wantFileName)
			}
			if gotRedirectURL != tt.wantRedirectURL {
				t.Errorf("AcceptRange() gotRedirectURL = %v, want %v", gotRedirectURL, tt.wantRedirectURL)
			}
		})
	}
}

func TestPostFile(t *testing.T) {
	patchesTryCountGetRespRedirect := gomonkey.ApplyFuncSeq(TryCountGetRespRedirect, []gomonkey.OutputCell{
		{Values: gomonkey.Params{&http.Response{StatusCode: 200, Body: new(readWriteCloserImpl)}, "", nil}, Times: 1},
		{Values: gomonkey.Params{nil, "", errors.New("getRespRedirect fail")}, Times: 1},
	})
	defer patchesTryCountGetRespRedirect.Reset()

	tmpFile := path.Join(os.TempDir(), uuid.New().String())
	if err := ioutil.WriteFile(tmpFile, []byte{1, 2, 3}, os.ModePerm); err != nil {
		return
	}
	defer os.Remove(tmpFile)

	ctx := context.Background()

	type args struct {
		ctx       context.Context
		fieldname string
		file      string
		url       string
		header    http.Header
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"normal", args{ctx, "fieldname", tmpFile, "", GetDefaultHeader()}, []byte{}, false},
		{"err_os.Open", args{ctx, "fieldname", "", "", nil}, nil, true},
		{"err_TryCountGetRespRedirect", args{ctx, "fieldname", tmpFile, "", nil}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PostFile(tt.args.ctx, tt.args.fieldname, tt.args.file, tt.args.url, tt.args.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("[%v] PostFile() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("[%v] PostFile() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
