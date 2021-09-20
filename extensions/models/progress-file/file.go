package progress_file

import (
	"fmt"
	"github.com/vortex14/gotyphoon/extensions/bar"
	"io"
	"os"
	"sync"
)

type File struct {
	io.Reader
	bar *bar.Bar
	File *os.File

	closeOnce sync.Once

	total int64
	fileSize int64
	OnFinish func(f *os.File)

}

func (f *File) getFileSize() int64 {
	stat, _ := f.File.Stat()
	return stat.Size()
}

func (f *File) finish()  {
	f.closeOnce.Do(func() {
		f.bar.Finish()
		if f.OnFinish != nil { f.OnFinish(f.File) }
	})
}

func (f *File) Read(p []byte) (int, error) {
	n, err := f.Reader.Read(p)

	if  f.total == 0 {
		f.fileSize = f.getFileSize()

		f.bar = &bar.Bar{
			Description: fmt.Sprintf("Reading %s ... ", f.File.Name()),
		}

		f.bar.NewOption(int64(n), f.fileSize)
	}

	if err == nil {
		f.total += int64(n)
		f.bar.IncCur(f.total)
	}

	switch {
	case f.total == f.fileSize:
		f.finish()
	}

	return n, err
}