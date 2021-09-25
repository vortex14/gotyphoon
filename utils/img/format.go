package img

import (
	"image"
	"io"
	"mime"
)

// guessImageFormat Guess image format from gif/jpeg/png/webp
func guessImageFormat(r io.Reader) (format string, err error) {
	_, format, err = image.DecodeConfig(r)
	return
}

// GuessImageMimeTypes image mime types from gif/jpeg/png/webp
func GuessImageMimeTypes(r io.Reader) string {
	format, _ := guessImageFormat(r)
	if format == "" {
		return ""
	}
	return mime.TypeByExtension("." + format)
}
