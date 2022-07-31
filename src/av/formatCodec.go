package av

// formatName 视频文件封装的 format 名称
type formatName string

// 导出格式
const (
	formatMp4     formatName = "mp4"
	formatMov     formatName = "mov"
	format3gp     formatName = "3gp"
	format3g2     formatName = "3g2"
	formatWebm    formatName = "webm"
	formatOpus    formatName = "opus"
	formatVorbis  formatName = "ogg"
	formatMp3     formatName = "mp3"
	formatM4a     formatName = "m4a"
	formatAc3     formatName = "ac3"
	formatAmr     formatName = "amr"
	formatGif     formatName = "gif"
	formatWebp    formatName = "webp"
	formatPng     formatName = "png"
	formatTs      formatName = "mpegts"
	formatDefault formatName = formatMp4
)

func (f formatName) getExt() string {
	if f == formatTs {
		return ".ts"
	}
	return "." + string(f)
}

var allFormat = map[formatName]bool{
	formatMp4:    true,
	formatMov:    true,
	format3gp:    true,
	format3g2:    true,
	formatWebm:   true,
	formatOpus:   true,
	formatVorbis: true,
	formatMp3:    true,
	formatM4a:    true,
	formatAc3:    true,
	formatAmr:    true,
	formatGif:    true,
	formatWebp:   true,
	formatPng:    true,
	formatTs:     true,
}

// codecName 视频编解码名称
type codecName string

const (
	codecH264 codecName = "h264"
	codecH265 codecName = "h265"
	codecVp8  codecName = "vp8"
	codecVp9  codecName = "vp9"
	codecGif  codecName = "gif"
	codecWebp codecName = "webp"
	codecPng  codecName = "png"

	codecMp3    codecName = "mp3"
	codecVorbis codecName = "vorbis"
	codecOpus   codecName = "opus"
	codecAac    codecName = "aac"
	codecAc3    codecName = "ac3"
	codecAmrNb  codecName = "amr_nb"
	codecAmrWb  codecName = "amr_wb"
)

var (
	// 视频和音频的建议后缀列表，越靠前权重越大
	videoCodecFormats = map[codecName][]formatName{
		codecH264: {formatMp4, formatMov, format3gp},
		codecH265: {formatMp4, formatMov, format3gp},
		codecVp8:  {formatWebm, formatMov, formatMp4},
		codecVp9:  {formatWebm, formatMp4, formatMov},
		codecGif:  {formatGif, formatMp4},
		codecWebp: {formatWebp},
		codecPng:  {formatPng},
	}
	// 对音频而言第一个元素表示有且仅有当前音频时的后缀
	audioCodecFormats = map[codecName][]formatName{
		codecMp3:    {formatMp3, formatMp4},
		codecVorbis: {formatVorbis, formatWebm, formatMp4},
		codecOpus:   {formatOpus, formatWebm, formatMp4},
		codecAac:    {formatM4a, formatMp4, formatMov},
		codecAc3:    {formatAc3, formatMp4, formatMov},
		codecAmrNb:  {formatAmr, format3gp, format3g2, formatMov},
		codecAmrWb:  {formatAmr, format3gp, format3g2, formatMov},
	}
)
