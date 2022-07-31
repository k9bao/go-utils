package av

import "testing"

func Test_formatName_getExt(t *testing.T) {
	tests := []struct {
		name string
		f    formatName
		want string
	}{
		{"fileFormatMp4", formatMp4, ".mp4"},
		{"fileFormatMov", formatMov, ".mov"},
		{"fileFormat3gp", format3gp, ".3gp"},
		{"fileFormat3g2", format3g2, ".3g2"},
		{"fileFormatWebm", formatWebm, ".webm"},
		{"fileFormatOpus", formatOpus, ".opus"},
		{"fileFormatVorbis", formatVorbis, ".ogg"},
		{"fileFormatMp3", formatMp3, ".mp3"},
		{"fileFormatM4a", formatM4a, ".m4a"},
		{"fileFormatAc3", formatAc3, ".ac3"},
		{"fileFormatAmr", formatAmr, ".amr"},
		{"fileFormatGif", formatGif, ".gif"},
		{"fileFormatWebp", formatWebp, ".webp"},
		{"fileFormatPng", formatPng, ".png"},
		{"fileFormatTs", formatTs, ".ts"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.getExt(); got != tt.want {
				t.Errorf("formatName.getExt() = %v, want %v", got, tt.want)
			}
		})
	}
}
