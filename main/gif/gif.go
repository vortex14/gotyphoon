package main

import (
	"github.com/vortex14/gotyphoon/utils"
	"github.com/vortex14/gotyphoon/utils/img"
	"path/filepath"
)

func main()  {
	path := utils.GetCurrentDir()
	path1 := filepath.Join(path, "1.png")
	path2 := filepath.Join(path, "2.png")
	files := []string {path1,
		path2,
	}

	fps := 2
	out := "./out.gif"

	err := img.BuildGif(files, fps, out)
	if err != nil {
		return 
	}
}
