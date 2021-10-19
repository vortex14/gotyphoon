package folder

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/vortex14/gotyphoon/interfaces"
)

type ExportOptions struct {
	TarPath string
	IsCompress bool
}

type Folder struct {
	Path string
	ExportOptions
}

func (f *Folder) Compress() error {
	return nil
}


func (f *Folder) Uncompress() error {

	return nil
}

func (f *Folder) GetDataFromDirectory() (error, interfaces.MapFileObjects ) {
	currentData := make(interfaces.MapFileObjects)

	files, err := ioutil.ReadDir(f.Path)
	if err != nil { return err, nil }

	for _, file := range files {
		typeFile := "file"
		if file.IsDir() { typeFile = "dir" }

		currentData[file.Name()] = &interfaces.FileObject{
			Type: typeFile,
			Path: file.Name(),
		}
	}

	return nil, currentData
}

func (f *Folder) IsExists(required []string) (error,bool)  {
	var status = true
	err, dataDir := f.GetDataFromDirectory()
	if err != nil { return err, false}

	for _, reqFile := range required {
		if _, ok := dataDir[reqFile]; !ok {
			status = false
		}

	}

	return nil, status
}

func (f *Folder) IsExist(path string) bool  {
	var status = false
	if _, err := os.Stat(filepath.Join(f.Path, path)); !os.IsNotExist(err) {
		status = true
	}

	return status
}