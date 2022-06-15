package utils

import (
	"embed"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/vortex14/gotyphoon/interfaces"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//go:embed test-dir
var testDir embed.FS

func TestFormatGoTemplateExport(t *testing.T) {
	Convey("format goTemplate and export", t, func() {

		tmp := `FROM {{.TYPHOON_IMAGE}}`

		goTemplateDocker := &interfaces.GoTemplate{
			Source:     tmp,
			ExportPath: "Dockerfile",
			Data: map[string]string{
				"TYPHOON_IMAGE": "typhoon-lite",
			},
		}

		u := Utils{}
		u.GoRunTemplate(goTemplateDocker)

		file, err := os.ReadFile("Dockerfile")
		if err != nil {
			return
		}

		So(string(file), ShouldEqual, "FROM typhoon-lite")

		_ = os.Remove("Dockerfile")
	})
}

func TestCopyFileAndReplaceLabelEmbed(t *testing.T) {
	Convey("copy file", t, func() {
		mainFolder := "test-dir"
		testFile := "test-tml-.tml"
		exportPath := "export-tml-file.tml"

		dir, errSub := fs.Sub(testDir, mainFolder)
		So(errSub, ShouldBeNil)

		u := Utils{}

		err := u.CopyFileAndReplaceLabel(exportPath, &interfaces.ReplaceLabel{Label: "{{-test-}}", Value: "::TEST::"}, &interfaces.FileObject{
			Path: ".",
			Name: testFile,
		}, dir)
		So(err, ShouldBeNil)

		Convey("check file", func() {
			readFile, err := os.ReadFile(exportPath)
			So(err, ShouldBeNil)
			So(string(readFile), ShouldEqual, "::TEST::")

			path, err := os.Getwd()
			So(err, ShouldBeNil)

			removeDir := filepath.Join(path, exportPath)
			println("remove: ", removeDir)
			errR := os.RemoveAll(removeDir)
			So(errR, ShouldBeNil)

		})
	})
}

func TestCopyDir(t *testing.T) {
	name := "test-dir"
	Convey("create directories with -.tml", t, func() {

		err := os.MkdirAll(name+"/test-example/test-subdir", 0755)

		So(err, ShouldEqual, nil)

		Convey("Create a new files", func() {
			f, err1 := os.Create(name + "/test-example/1.txt")
			So(err1, ShouldEqual, nil)

			_, err1 = f.WriteString("Hello 1.txt")
			So(err1, ShouldEqual, nil)

			f, err2 := os.Create(name + "/test-example/test-subdir/1-.tml")
			So(err2, ShouldEqual, nil)

			_, err2 = f.WriteString("Hello-.tml")
			So(err2, ShouldEqual, nil)

		})

		dir, errSub := fs.Sub(testDir, "test-dir")
		if errSub != nil {
			panic(err)
		}

		Convey("copy dir", func() {
			So(errSub, ShouldEqual, nil)

			errC := (&Utils{}).CopyDir("new_copy_dir", dir)
			So(errC, ShouldEqual, nil)
		})

		Convey("check -.tml", func() {
			_, errTml := os.Stat("new_copy_dir/test-example/test-subdir/1-.tml")
			status := strings.Contains(errTml.Error(), "no such file or directory")
			So(status, ShouldEqual, true)
		})

		Convey("check 1.txt", func() {
			_, err1txt := os.Stat("new_copy_dir/test-example/1.txt")
			So(err1txt, ShouldEqual, nil)
		})

		Convey("check init1.py", func() {
			_, errpy := os.Stat("new_copy_dir/test-example/init1.py")
			So(errpy, ShouldEqual, nil)
		})

		Convey("check __init__.py", func() {
			_, errinit := os.Stat("new_copy_dir/test-example/__init__.py")
			So(errinit, ShouldEqual, nil)
		})

		Convey("check init.py", func() {
			_, err := os.Stat("new_copy_dir/test-example/init.py")
			status := strings.Contains(err.Error(), "no such file or directory")
			So(status, ShouldEqual, true)
		})

		Convey("remove new_copy_dir", func() {

			path, err := os.Getwd()
			So(err, ShouldBeNil)

			removeDir := filepath.Join(path, "new_copy_dir")
			println(removeDir)
			errR := os.RemoveAll(removeDir)
			So(errR, ShouldBeNil)
		})

	})

}
