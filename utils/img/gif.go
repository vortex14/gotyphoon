package img

import (
	"image"
	"image/gif"
	"os"
)

// WriteGif Write a img2gif file from a paletted image slice
// d = delay in 100ths of a second per frame.
func WriteGif(im *[]*image.Paletted, d int, exportPath string) error {
	//
	g := &gif.GIF{}

	for _, i := range *im {
		g.Image = append(g.Image, i)
		g.Delay = append(g.Delay, d)
	}
	f, err := os.Create(exportPath)

	if err != nil {
		println(err)
	}
	defer f.Close()
	return gif.EncodeAll(f, g)
}


// BuildGif Executes the functions above in the right order.
// Takes an array of file paths pointing to images as input.
// p is a path to the output file.
// fps: frames per second.
func BuildGif(files []string, fps int, exportPath string) error {
	imgs := ReadImages(files)

	im_p := EncodeImgPaletted(&imgs)
	println(im_p, len(files), len(imgs))
	return WriteGif(&im_p, 100 / fps, exportPath)
}
