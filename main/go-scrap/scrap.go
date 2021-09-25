package main
import (
	"fmt"
	"github.com/kbinani/screenshot"
	"image/png"
	"os"
	"time"
)

func main() {
	n := screenshot.NumActiveDisplays()
	s := 0
	for i := 0; i < n; i++ {
		for {
			bounds := screenshot.GetDisplayBounds(i)

			img, err := screenshot.CaptureRect(bounds)
			if err != nil {
				panic(err)
			}
			fileName := fmt.Sprintf("%d_%d:%dx%d.jpeg", i, s, bounds.Dx(), bounds.Dy())
			file, _ := os.Create(fileName)
			defer file.Close()
			png.Encode(file, img)



			fmt.Printf("#%d : %v \"%s\"\n", i, bounds, fileName)


			time.Sleep(time.Duration(1) * time.Second)

			s++
		}

	}
}