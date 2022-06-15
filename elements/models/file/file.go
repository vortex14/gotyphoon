package file

import (
	"io"
	"os"

	"github.com/vortex14/gotyphoon/elements/models/awaitabler"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/elements/models/singleton"
	"github.com/vortex14/gotyphoon/interfaces"
)

type File struct {
	io.Reader
	*label.MetaInfo

	Link      *os.File
	Path      string

	awaitabler.Object
	singleton.Singleton
	LOG       interfaces.LoggerInterface

	OnFinish  func(f *os.File)
}

func (f *File) Init()  {
	if len(f.Path) > 0 {
		file, err := os.Open(f.Path)
		if err != nil { return }
		f.Link = file
		f.Reader = file
	}
}

func (f *File) GetFileSize() int64 {
	stat, _ := f.Link.Stat()
	return stat.Size()
}



func (f *File) finish()  {
	f.Destruct(func() {
		if f.OnFinish != nil { f.OnFinish(f.Link) }
	})
}
