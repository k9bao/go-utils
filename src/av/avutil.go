// Package av 提供audio/video常用函数封装，比如获取时长等
package av

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"errors"
	"go-utils/src/tools/fs"
	
)

// 可执行程序及部分后缀
var (
	ffmpegBin = "ffmpeg"

	Mp4Ext = formatMp4.getExt()
	Mp3Ext = formatMp3.getExt()
)

// 定义一些常用变量
const (
	soundstrethBin = "soundstretch"

	pixFmtYUVA420p = "yuva420p"
)

// fileHeader 文件头常量
type fileHeader string

// 视频文件开头几个字节内容，增加内容同时修改 fileHeader2ExtMap 内容
const (
	ID3v2 fileHeader = "ID3"
	FLAC1 fileHeader = "FLAC"
	FLAC2 fileHeader = "\177FLAC"
)

// maxFileHeaderLen 最大文件头字节数
var maxFileHeaderLen int

// 文件头到后缀的映射
var fileHeader2ExtMap = map[fileHeader]string{
	ID3v2: Mp3Ext,
	FLAC1: Mp3Ext,
	FLAC2: Mp3Ext,
}

// CutMedia cut video/audio clip
func CutMedia(ctx context.Context, inputPath, outputPath string, start, dur time.Duration) (string, error) {
	info, err := Probe(ctx, inputPath)
	if err != nil {
		return "", err
	}
	var (
		encodeName string
		decodeName string
		pixFmt     string
	)
	if info.GetVideoCodec() == string(codecVp8) {
		encodeName = "libvpx"
		decodeName = encodeName
		if info.HasAlpha() {
			pixFmt = pixFmtYUVA420p
		}
	} else if info.GetVideoCodec() == string(codecVp9) {
		encodeName = "libvpx-vp9"
		decodeName = encodeName
		if info.HasAlpha() {
			pixFmt = pixFmtYUVA420p
		}
	}

	cmd := []string{
		ffmpegBin, "-y",
		"-loglevel", "error",
	}
	inputOpt := []string{
		"-ss", strconv.FormatInt(start.Microseconds(), 10) + "us",
		"-t", strconv.FormatInt(dur.Microseconds(), 10) + "us",
	}
	if decodeName != "" {
		inputOpt = append(inputOpt, "-c:v", decodeName)
	}
	outputOpt := []string{
		"-i", inputPath,
	}
	if encodeName != "" {
		outputOpt = append(outputOpt, "-c:v", encodeName)
	}
	if pixFmt != "" {
		outputOpt = append(outputOpt, "-pix_fmt", pixFmt)
	}

	newOutputPath := fs.GetNameWithNewExt(outputPath, info.GetSuggestedExtFromCodec())

	cmd = append(cmd, inputOpt...)
	cmd = append(cmd, outputOpt...)
	cmd = append(cmd, newOutputPath)
	return newOutputPath, fs.RunSysCommand(ctx, cmd, nil)
}

// MediaReverse reverse video
func MediaReverse(ctx context.Context, inputPath, outputPath string) error {
	cmd := []string{ffmpegBin, "-y", "-i", inputPath, "-vf", "reverse", "-af", "areverse", outputPath}
	return fs.RunSysCommand(ctx, cmd, nil)
}

// MergeTS convert m3u8 to mp4
func MergeTS(ctx context.Context, inputPath, outputPath string) error {
	cmd := []string{ffmpegBin, "-y", "-loglevel", "error", "-i", inputPath, "-c", "copy", outputPath}
	return fs.RunSysCommand(ctx, cmd, nil)
}

// ConvertToWav convert video to wav format
func ConvertToWav(ctx context.Context, inputPath string, start, dur time.Duration, outputPath string) error {
	cmd := []string{
		ffmpegBin,
		"-y",
		"-loglevel", "error",
		"-ss", strconv.FormatInt(start.Microseconds(), 10) + "us",
		"-t", strconv.FormatInt(dur.Microseconds(), 10) + "us",
		"-i", inputPath,
		"-f", "wav",
		outputPath,
	}
	return fs.RunSysCommand(ctx, cmd, nil)
}

// SoundstretchProcess soundtouch, pitch=n : Change sound pitch by n semitones (n=-60..+60 semitones)
// input and output only support wav
func SoundstretchProcess(ctx context.Context, inputPath, outputPath string, pitch float32) error {
	if pitch < -60 || pitch > 60 {
		return errors.New("pitch para is error,pitch should in（-60,60)")
	}
	if !strings.HasSuffix(inputPath, ".wav") || !strings.HasSuffix(outputPath, ".wav") {
		return errors.New("only support .wav")
	}
	cmd := []string{soundstrethBin, inputPath, outputPath, fmt.Sprintf("-pitch=%f", pitch)}
	return fs.RunSysCommand(ctx, cmd, nil)
}

// FormatConvert  ffmpeg -i inputPath -c copy outputPath
func FormatConvert(ctx context.Context, inputPath, ext string) (string, error) {
	outputPath := filepath.Join(filepath.Dir(inputPath), fs.GetFileName(inputPath)+"_convert"+ext)
	cmd := []string{
		ffmpegBin,
		"-y",
		"-loglevel", "error",
		"-i", inputPath,
		"-c", "copy",
		"-strict", "-2",
		outputPath,
	}
	return outputPath, fs.RunSysCommand(ctx, cmd, nil)
}

// ResizeMediaFitIn resize media file in the specified size
func ResizeMediaFitIn(ctx context.Context, inputPath string, width int, height int, defaultExt string) (string, error) {
	probeInfo, err := Probe(ctx, inputPath)
	if err != nil {
		logs.Log.Errorf("failed to probe %s err = %+v", inputPath, err)
		return "", err
	}

	for _, stream := range probeInfo.Streams {
		if stream.CodecType != CodecTypeVideo {
			continue
		}
		if stream.CodedWidth <= width && stream.CodedHeight <= height {
			continue
		}

		// 2的倍数
		newWidth := -2
		newHeight := -2
		if stream.CodedWidth*height > stream.CodedHeight*width {
			newWidth = width & ^1
		} else {
			newHeight = height & ^1
		}

		if defaultExt == "" {
			defaultExt = probeInfo.GetSuggestedExtFromCodec()
		}

		outputPath := fs.FileNameAppend(inputPath, fmt.Sprintf("_%vx%v", newWidth, newHeight), defaultExt)
		cmd := []string{
			ffmpegBin,
			"-y",
			"-loglevel", "error",
			"-i", inputPath,
			"-vf", fmt.Sprintf("scale=%d:%d", newWidth, newHeight),
			"-strict", "-2",
			outputPath,
		}
		if err = fs.RunSysCommand(ctx, cmd, nil); err != nil {
			logs.Log.Errorf("failed to resize %s %dx%d => %s err = %+v", inputPath, newWidth, newHeight, outputPath, err)
			return "", err
		}
		return outputPath, nil
	}
	return inputPath, nil
}

// GetSuggestedExtFromContent 从文件头中获取文件后缀
func GetSuggestedExtFromContent(ctx context.Context, filePath string) string {
	ext := filepath.Ext(filePath)
	fileHeader, err := fs.ReadFileWithSize(filePath, 0, maxFileHeaderLen)
	if err != nil {
		logs.Log.Wainf("failed to ReadFileWithSize: %s err %+v", filePath, err)
		return ext
	}
	// logs.Log.Debugf("ReadFileWithSize %v from %s", fileHeader, filePath)
	for k, v := range fileHeader2ExtMap {
		if strings.ToUpper(string(fileHeader[:len(k)])) == string(k) {
			return v
		}
	}
	return ext
}

// init 计算文件头最大长度
func init() {
	for k := range fileHeader2ExtMap {
		if len(k) > maxFileHeaderLen {
			maxFileHeaderLen = len(k)
		}
	}
}
