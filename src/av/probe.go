package av

import (
	"context"
	"encoding/json"
	"path"

	"go-utils/src/tools/algorithm"
	"go-utils/src/tools/fs"
)

var ffprobeBin = "ffprobe"

// Disposition 布局, 具体参数含义详见 ffprobe 描述 https://ffmpeg.org/ffprobe.html
type Disposition struct {
	Default         int `json:"default"`
	Dub             int `json:"dub"`
	Original        int `json:"original"`
	Comment         int `json:"comment"`
	Lyrics          int `json:"lyrics"`
	Karaoke         int `json:"karaoke"`
	Forced          int `json:"forced"`
	HearingImpaired int `json:"hearing_impaired"`
	VisualImpaired  int `json:"visual_impaired"`
	CleanEffects    int `json:"clean_effects"`
	AttachedPic     int `json:"attached_pic"`
	TimedThumbnails int `json:"timed_thumbnails"`
}

// Tags tag信息，具体参数含义详见 ffprobe 描述 https://ffmpeg.org/ffprobe.html
type Tags struct {
	Language       string `json:"language"`
	AlphaMode      string `json:"ALPHA_MODE"`
	AlphaModeLower string `json:"alpha_mode"`
	Duration       string `json:"DURATION"`
}

// SideDataList side data 具体参数含义详见 ffprobe 描述 https://ffmpeg.org/ffprobe.html
type SideDataList struct {
	SideDataType  string  `json:"side_data_type,omitempty"`
	DisplayMatrix string  `json:"displaymatrix,omitempty"`
	Rotation      float32 `json:"rotation,omitempty"`
}

// CodecType 编码类型 描述 https://ffmpeg.org/ffprobe.html
type CodecType string

const (
	// CodecTypeVideo 视频类型
	CodecTypeVideo CodecType = "video" // 视频
	// CodecTypeAudio 音频类型
	CodecTypeAudio CodecType = "audio" // 音频
)

// Streams 流信息，具体参数含义详见 ffprobe 描述 https://ffmpeg.org/ffprobe.html
type Streams struct {
	Index              int            `json:"index"`
	CodecName          string         `json:"codec_name"`
	CodecLongName      string         `json:"codec_long_name"`
	Profile            string         `json:"profile,omitempty"`
	CodecType          CodecType      `json:"codec_type"`
	CodecTimeBase      string         `json:"codec_time_base"`
	CodecTagString     string         `json:"codec_tag_string"`
	CodecTag           string         `json:"codec_tag"`
	Width              int            `json:"width,omitempty"`
	Height             int            `json:"height,omitempty"`
	CodedWidth         int            `json:"coded_width,omitempty"`
	CodedHeight        int            `json:"coded_height,omitempty"`
	ClosedCaptions     int            `json:"closed_captions,omitempty"`
	HasBFrames         int            `json:"has_b_frames,omitempty"`
	SampleAspectRatio  string         `json:"sample_aspect_ratio,omitempty"`
	DisplayAspectRatio string         `json:"display_aspect_ratio,omitempty"`
	PixFmt             string         `json:"pix_fmt,omitempty"`
	Level              int            `json:"level,omitempty"`
	ColorRange         string         `json:"color_range,omitempty"`
	Refs               int            `json:"refs,omitempty"`
	RFrameRate         string         `json:"r_frame_rate"`
	AvgFrameRate       string         `json:"avg_frame_rate"`
	TimeBase           string         `json:"time_base"`
	StartPts           int            `json:"start_pts"`
	StartTime          string         `json:"start_time"`
	Disposition        Disposition    `json:"disposition"`
	Tags               Tags           `json:"tags,omitempty"`
	SideDataList       []SideDataList `json:"side_data_list,omitempty"`
	SampleFmt          string         `json:"sample_fmt,omitempty"`
	SampleRate         string         `json:"sample_rate,omitempty"`
	Channels           int            `json:"channels,omitempty"`
	ChannelLayout      string         `json:"channel_layout,omitempty"`
	BitsPerSample      int            `json:"bits_per_sample,omitempty"`
}

// FormalTags 格式tag，具体参数含义详见 ffprobe 描述 https://ffmpeg.org/ffprobe.html
type FormalTags struct {
	Encoder string `json:"ENCODER"`
}

// Format 格式，具体参数含义详见 ffprobe 描述 https://ffmpeg.org/ffprobe.html
type Format struct {
	Filename       string     `json:"filename"`
	NbStreams      int        `json:"nb_streams"`
	NbPrograms     int        `json:"nb_programs"`
	FormatName     string     `json:"format_name"`
	FormatLongName string     `json:"format_long_name"`
	StartTime      string     `json:"start_time"`
	Duration       string     `json:"duration"`
	Size           string     `json:"size"`
	BitRate        string     `json:"bit_rate"`
	ProbeScore     int        `json:"probe_score"`
	Tags           FormalTags `json:"tags"`
}

// ProbeInfo probe信息，具体参数含义详见 ffprobe 描述 https://ffmpeg.org/ffprobe.html
type ProbeInfo struct {
	InputPath string    `json:"-"`
	Streams   []Streams `json:"streams"`
	Format    Format    `json:"format"`
}

// GetVideoStream 获取视频流
func (p *ProbeInfo) GetVideoStream() *Streams {
	for _, s := range p.Streams {
		if s.CodecType == CodecTypeVideo {
			return &s
		}
	}
	return nil
}

// GetAudioStream 获取音频流
func (p *ProbeInfo) GetAudioStream() *Streams {
	for _, s := range p.Streams {
		if s.CodecType == CodecTypeAudio {
			return &s
		}
	}
	return nil
}

// GetVideoCodec 获取视频编码
func (p *ProbeInfo) GetVideoCodec() string {
	s := p.GetVideoStream()
	if s != nil {
		return s.CodecName
	}
	return ""
}

// GetAudioCodec 获取音频编码
func (p *ProbeInfo) GetAudioCodec() string {
	s := p.GetAudioStream()
	if s != nil {
		return s.CodecName
	}
	return ""
}

// HasAlpha 判断是否包含 alpha 通道
func (p *ProbeInfo) HasAlpha() bool {
	s := p.GetVideoStream()
	if s != nil {
		if s.Tags.AlphaMode+s.Tags.AlphaModeLower == "1" {
			return true
		}
	}
	return false
}

// GetFormatDuration 获取文件时长
func (p *ProbeInfo) GetFormatDuration() float64 {
	dur := algorithm.ParseFloat(p.Format.Duration, 0)
	return dur
}

// getDefaultExt 获取当前流的默认后缀
func (p *ProbeInfo) getDefaultExt() string {
	defaultExt := path.Ext(p.InputPath)
	if len(defaultExt) == 0 {
		defaultExt = formatDefault.getExt()
	}
	return defaultExt
}

// GetSuggestedExtFromCodec 根据流里的已有信息获取他的建议文件后嘴
func (p *ProbeInfo) GetSuggestedExtFromCodec() string {
	audioCodec := p.GetAudioCodec()
	videoCodec := p.GetVideoCodec()
	suggestFormatFromVideo := videoCodecFormats[codecName(videoCodec)]
	suggestFormatFromAudio := audioCodecFormats[codecName(audioCodec)]
	defaultExt := p.getDefaultExt()

	if len(videoCodec) == 0 { // 仅有音频流时根据音频决定文件后缀
		return getExt(suggestFormatFromAudio, 0, defaultExt)
	} else if len(audioCodec) == 0 { // 仅有视频流时根据视频决定文件后缀
		return getExt(suggestFormatFromVideo, 0, defaultExt)
	}

	if len(suggestFormatFromVideo) == 0 { // 未匹配到视频后缀，则使用音频次优后缀，次优后缀同时可容纳音频和视频
		return getExt(suggestFormatFromAudio, 1, defaultExt)
	} else if len(suggestFormatFromAudio) == 0 { // 未匹配到音频后缀，则使用视频最优后缀
		return getExt(suggestFormatFromVideo, 0, defaultExt)
	}

	for _, audioExt := range suggestFormatFromAudio[1:] {
		for _, videoExt := range suggestFormatFromVideo {
			if videoExt == audioExt {
				return videoExt.getExt()
			}
		}
	}
	return defaultExt
}

func getExt(slice []formatName, index int, defaultExt string) string {
	if slice == nil || index > len(slice) {
		return defaultExt
	}
	return slice[index].getExt()
}

// GetSuggestedExtFromFormat 获取 Format 对应的后缀
func (p *ProbeInfo) GetSuggestedExtFromFormat() string {
	if allFormat[formatName(p.Format.FormatName)] {
		return formatName(p.Format.FormatName).getExt()
	}
	return ""
}

// Probe get video probe info
func Probe(ctx context.Context, inputPath string) (*ProbeInfo, error) {
	cmd := []string{
		ffprobeBin,
		"-loglevel", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		inputPath,
	}
	jsonStr, err := fs.RunSysCommandRet(ctx, cmd, nil)
	if err != nil {
		return nil, err
	}
	info := &ProbeInfo{}
	unmarsha1Err := json.Unmarshal(jsonStr, info)
	if unmarsha1Err != nil {
		return nil, unmarsha1Err
	}
	info.InputPath = inputPath
	return info, nil
}
