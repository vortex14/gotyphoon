package folder

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