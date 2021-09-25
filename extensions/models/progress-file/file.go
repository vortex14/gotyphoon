package progress_file

import (
	"fmt"

	"github.com/vortex14/gotyphoon/elements/models/bar"
	"github.com/vortex14/gotyphoon/elements/models/file"
	"github.com/vortex14/gotyphoon/log"
)

type File struct {
	file.File

	bar       *bar.Bar
	total     int64
}


func (f *File) finish()  {
	f.Destruct(func() {
		f.bar.Finish()
		if f.OnFinish != nil { f.OnFinish(f.Link) }
	})
}

func (f *File) Read(p []byte) (int, error) {
	if f.Reader == nil && len(f.Path) > 0 { f.Init() }
	n, err := f.Reader.Read(p)

	if  f.total == 0 {
		if f.LOG == nil { f.LOG = log.New(log.D{"from": f.Link.Name()}) }
		f.LOG.Debug("start read")
		f.bar = &bar.Bar{
			Description: fmt.Sprintf("%s ...",f.Description),
		}
		f.bar.NewOption(int64(n), f.GetFileSize())
	}

	if err == nil {
		f.total += int64(n)
		f.bar.IncCur(f.total)
	}

	switch {
	case f.total == f.GetFileSize():
		f.finish()
	}

	return n, err
}