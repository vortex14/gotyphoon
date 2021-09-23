package progress_file

import (
	"fmt"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"io"
	"os"

	"github.com/vortex14/gotyphoon/elements/models/awaitable"
	"github.com/vortex14/gotyphoon/elements/models/singleton"
	"github.com/vortex14/gotyphoon/extensions/bar"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

type File struct {
	io.Reader
	*label.MetaInfo

	bar       *bar.Bar
	File      *os.File
	Path      string

	awaitable.Object
	singleton.Singleton
	LOG       interfaces.LoggerInterface

	total     int64
	fileSize  int64
	OnFinish  func(f *os.File)

}

func (f *File) init()  {
	if len(f.Path) > 0 {
		file, err := os.Open(f.Path)
		if err != nil { return }
		f.File = file
		f.Reader = file
	}
}

func (f *File) getFileSize() int64 {
	stat, _ := f.File.Stat()
	return stat.Size()
}

func (f *File) finish()  {
	f.Destruct(func() {
		f.bar.Finish()
		if f.OnFinish != nil { f.OnFinish(f.File) }
	})
}

func (f *File) Read(p []byte) (int, error) {
	if f.Reader == nil && len(f.Path) > 0 { f.init() }
	n, err := f.Reader.Read(p)

	if  f.total == 0 {
		f.fileSize = f.getFileSize()
		if f.LOG == nil { f.LOG = log.New(log.D{"from": f.File.Name()}) }
		f.LOG.Debug("start read")
		f.bar = &bar.Bar{
			Description: fmt.Sprintf("%s ...",f.Description),
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