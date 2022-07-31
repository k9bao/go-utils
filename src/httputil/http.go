// Package httputil 对http常用功能的封装，比如重试，重定向等
package httputil

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"errors"
	"go-utils/src/config"
	"errors"
	
	"github.com/google/uuid"
)

// GetDefaultHeader 获取默认构造的 header
func GetDefaultHeader() http.Header {
	defaultHeader := make(http.Header)
	defaultHeader.Set("Accept-Language", "en-US,en;q=0.8,zh-CN;q=0.6,zh;q=0.4")
	defaultHeader.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3")
	defaultHeader.Add("User-Agent", "AppleWebKit/537.36 (KHTML, like Gecko)")
	defaultHeader.Add("User-Agent", "Chrome/56.0.2924.87 Safari/537.36")
	defaultHeader.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	return defaultHeader
}

// getNewRequest NewRequest 简单封装
func getNewRequest(
	ctx context.Context,
	method, url string,
	header http.Header,
	data *bytes.Buffer,
) (*http.Request, error) {
	var req *http.Request
	var err error
	if data != nil {
		req, err = http.NewRequestWithContext(ctx, method, url, data)
	} else {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
	}
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s %v => %+v", url, data, err))
	}
	if header != nil {
		req.Header = header
	}

	return req, nil
}

// getResp NewRequest+Do 简单封装
func getResp(
	ctx context.Context,
	method, url string,
	header http.Header,
	data []byte,
) (*http.Response, error) {
	reqBody := bytes.NewBuffer(data)
	req, err := getNewRequest(ctx, method, url, header, reqBody)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Do fail: %v %+v %s => %+v", url, req.Header, string(data), err))
	}
	if resp.StatusCode >= 400 {
		logs.Log.Wainf("%+v %+v %s => %+v", url, req.Header, string(data), resp)
	}
	return resp, nil
}

func find(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// getRespRedirect NewRequest 简单封装，支持重定向
func getRespRedirect(
	ctx context.Context,
	method, url string,
	header http.Header,
	data []byte,
) (*http.Response, string, error) {
	codes := []int{301, 302, 303, 307, 308}
	var resp *http.Response
	var err error
	maxRedirectTimes := int(config.ServerCnf.ControlConfig.MaxRedirectCounts)
	redirectURL := url
	for i := 0; i < maxRedirectTimes; i++ {
		resp, err = getResp(ctx, method, redirectURL, header, data)
		if err != nil {
			return nil, redirectURL, errorFromResponse(resp, err)
		}
		if find(codes, resp.StatusCode) {
			redirectURL = resp.Header["Location"][0]
			resp.Body.Close()
		} else {
			break
		}
	}
	if resp.StatusCode > 299 {
		// the body was closed in redirect codes
		respErr := errorFromResponse(resp, err)
		if !find(codes, resp.StatusCode) {
			resp.Body.Close()
		}
		return nil, redirectURL, respErr
	}
	return resp, redirectURL, nil
}

func errorFromResponse(resp *http.Response, err error) error {
	if resp == nil && err == nil {
		return nil
	}

	var errMsg string

	if err != nil {
		if e, ok := err.(*errors.Error); ok {
			return e
		}
		errMsg += fmt.Sprintf("%+v", err)
	}

	errCode := errors.RetResponseCodeErr
	if resp != nil {
		errCode = resp.StatusCode + errors.RetHTTPStatusCodeBegin
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errMsg += fmt.Sprintf("%+v", resp)
		} else {
			errMsg += fmt.Sprintf("%+v", string(content))
		}
	}

	if len(errMsg) == 0 {
		errMsg = "unknown error"
	}
	return errors.New(errCode, errMsg)
}

// GetHTTPFileName 获取默认文件名 ext:默认后缀比如".mp4"
func GetHTTPFileName(uri string, resp *http.Response, defaultExt string, pre string) string {
	if pre == "" {
		pre = uuid.New().String()
	}
	pre += "_"
	getDispos := resp.Header.Get("Content-Disposition")
	if getDispos != "" {
		_, params, err := mime.ParseMediaType(getDispos)
		if err == nil {
			if name, ok := params["filename"]; ok {
				return pre + name
			}
		}
	}

	u, err := url.Parse(uri)
	if err == nil {
		index := strings.LastIndex(u.Path, "/")
		if index != -1 {
			name := u.Path[index+1:]
			return pre + name
		}
	}

	return pre + defaultExt
}

// TryCountGetRespRedirect 多次尝试 GetRespRedirect
func TryCountGetRespRedirect(
	ctx context.Context,
	method, url string,
	header http.Header,
	data []byte,
) (resp *http.Response, redirectURL string, err error) {
	count := int(config.ServerCnf.ControlConfig.HTTPRequestRetryCounts)
	sleep := time.Second
	for i := 0; i < count; i++ {
		resp, redirectURL, err = getRespRedirect(ctx, method, url, header, data)
		if err == nil { // ok
			return
		}

		time.Sleep(sleep)
		logs.Log.Wainf("redirect: %v %+v %s retry count %d => %+v", url, header, string(data), i, err)
	}
	return nil, "", errorFromResponse(resp, err)
}

// AcceptRange 判断 url 是否支持按照字节下载，通过 GET 方法判断
func AcceptRange(ctx context.Context, url string) (fileSize int64, fileName, redirectURL string, err error) {
	header := GetDefaultHeader()
	header.Set("Range", "bytes=0-1")
	var resp *http.Response
	resp, redirectURL, err = TryCountGetRespRedirect(ctx, http.MethodGet, url, header, nil)
	if err != nil {
		return
	}
	fileName = GetHTTPFileName(url, resp, ".mp4", "")
	// 检查是否支持 断点续传
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Ranges
	if resp.Header.Get("Accept-Ranges") != "bytes" && resp.StatusCode != 206 {
		err = errors.New(fmt.Sprintf("not support Ranges, code=%v", resp.StatusCode))
		return
	}
	defer resp.Body.Close()

	// bytes 3600-5000/5000
	contentRange := resp.Header.Get("Content-Range")
	index := strings.LastIndex(contentRange, "/")
	fileSizeText := contentRange[index+1:]
	fileSize, err = strconv.ParseInt(fileSizeText, 10, 64)
	return
}

// PostFile 发送文件
func PostFile(ctx context.Context, fieldname, file, url string, header http.Header) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	formFile, err := writer.CreateFormFile(fieldname, file)
	if err != nil {
		errWarp := errors.New(fmt.Sprintf("Create form file failed: %v, %v, %+v", file, url, err))
		logs.Log.Error(errWarp)
		return nil, errWarp
	}

	srcFile, err := os.Open(file)
	if err != nil {
		errWarp := errors.New(fmt.Sprintf("Open source file failed: %v, %v, %+v", file, url, err))
		logs.Log.Error(errWarp)
		return nil, errWarp
	}
	defer srcFile.Close()
	_, err = io.Copy(formFile, srcFile)
	if err != nil {
		errWarp := errors.New(fmt.Sprintf("Write to form file falied: %v, %v, %+v", file, url, err))
		logs.Log.Error(errWarp)
		return nil, errWarp
	}

	writer.Close()

	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", writer.FormDataContentType())
	resp, _, err := TryCountGetRespRedirect(ctx, http.MethodPost, url, header, buf.Bytes())
	if err != nil {
		errWarp := errors.New(fmt.Sprintf("TryCountGetRespRedirect fail: %v %+v %+v", url, resp, err))
		logs.Log.Error(errWarp)
		return nil, errWarp
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
