package img

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

func Read(path string) image.Image {
	f, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
	}

	var img image.Image

	switch "" {
	case "PNG":
		img, err = png.Decode(f)
	default:
		img, _, err = image.Decode(f)
	}


	return img


}

// ReadImages Read images from a slice with file locations.
func ReadImages(paths []string) []image.Image {
	var im []image.Image

	for _, path := range paths {
		f, err := os.Open(path)

		if err != nil {
			fmt.Println(err)
		}
		var img image.Image
		var errF error
		//println(GuessImageMimeTypes(f))
		switch "" {
		case "PNG":
			img, errF = png.Decode(f)
		default:
			img, _, errF = image.Decode(f)
		}


		if errF != nil {
			fmt.Println(">>>",err)
		}

		im = append(im, img)
		err = f.Close()
		if err != nil {
			return nil
		}
	}
	return im
}
