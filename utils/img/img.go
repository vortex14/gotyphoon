package img

import (
	"bytes"
	"image"
	"image/gif"
)


// EncodeImgPaletted Encode images in *image.Paletted by encoding and decoding to img2gif and image.
// Encode and decode is necessary to convert jpeg and png to img2gif.
func EncodeImgPaletted(images *[]image.Image) []*image.Paletted {
	// Gif options
	opt := gif.Options{}
	var g []*image.Paletted

	for _, im := range *images {
		b := bytes.Buffer{}
		// Write img2gif file to buffer.
		err := gif.Encode(&b, im, &opt)

		if err != nil {
			println(err)
		}
		// Decode img2gif file from buffer to img.
		img, err := gif.Decode(&b)

		if err != nil {
			println(err)
		}

		// Cast img.
		i, ok := img.(*image.Paletted)
		if ok {
			g = append(g, i)
		}
	}
	return g
}

