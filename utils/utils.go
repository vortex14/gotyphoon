package utils

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-logfmt/logfmt"
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr"
	"github.com/olekukonko/tablewriter"
	"github.com/vortex14/gotyphoon/interfaces"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Utils struct {}

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

func (u *Utils) GoRunTemplate(goTemplate *interfaces.GoTemplate) bool {
	tmpl, _ := template.New("new").Parse(goTemplate.Source)
	status := true
	f, err := os.Create(goTemplate.ExportPath)
	if err != nil {
		log.Println("create file: ", err)
		status = false
	}

	err = tmpl.Execute(f, &goTemplate.Data)
	if err != nil {
		log.Print("execute: ", err)
		status = false
	}
	f.Close()

	return status
}

func (u *Utils) ParseLog(object *interfaces.FileObject) error {
	currentPath, _ := os.Getwd()
	logPath := fmt.Sprintf("%s/%s", currentPath, object.Path)
	dat, err := ioutil.ReadFile(logPath)

	color.Red("Log file path: %s", logPath)
	if err != nil {

		color.Red("Log file not found")
		os.Exit(0)


	}

	logDataMap := logfmt.NewDecoder(strings.NewReader(string(dat)))
	for logDataMap.ScanRecord() {
		for logDataMap.ScanKeyval() {
			color.Yellow("%s = %s", logDataMap.Key(), logDataMap.Value())
		}
	}



	return nil
}

func (u *Utils) CopyDir(name string, object *interfaces.FileObject) error {
	box := packr.NewBox(object.Path)
	errC := os.MkdirAll(name, 0755)
	if errC != nil{
		color.Red("CopyDir Error : %s",errC)
	}
	box.Walk(func(s string, file packd.File) error {
		//color.Yellow("s: %s, file: %+f \n", s, file)
		pathDir := filepath.Dir(s)
		//color.Yellow("data : %s, info: %s",s, pathDir)
		if pathDir != "." {
			errD := os.Mkdir(name+"/"+pathDir, 0755)
			if errD != nil{
				//color.Red("%s",errD)
			}

		}

		f, err := os.Create(name+"/"+s)
		if err != nil {
			log.Println("create err", err)
		}
		f.WriteString(file.String())
		_ = f.Close()

		return nil
	})
	return nil
}

func (u *Utils) CopyFile(ExportPath string, object *interfaces.FileObject) error {
	box := packr.NewBox(object.Path)
	f, err := os.Create(ExportPath)
	if err != nil {
		log.Println("create err", err)
	}

	dat, err := box.FindString(object.Name)

	if err != nil {

		color.Red("Log file not found")
		os.Exit(0)


	}

	_, errorWrite := f.WriteString(dat)
	if errorWrite != nil {
		color.Red("Can't write %s", ExportPath )
		os.Exit(0)

	}
	_ = f.Close()

	return nil


}

func (u *Utils) DumpToFile(object *interfaces.FileObject) error {
	f, err := os.Create(object.Path)
	if err != nil {
		log.Println("create err", err)
	}


	_, errorWrite := f.WriteString(object.Data)
	if errorWrite != nil {
		color.Red("Can't write %s", object.Path )
		os.Exit(0)

	}
	_ = f.Close()

	return nil


}

func (u *Utils) CopyFileAndReplaceLabel(name string, label *interfaces.ReplaceLabel, object *interfaces.FileObject) error {
	box := packr.NewBox(object.Path)
	templateFile, _ := box.FindString(object.Name)
	f, err := os.Create(name)
	if err != nil {
		log.Println("create err", err)
	}
	data := strings.ReplaceAll(templateFile, label.Label, label.Value)
	f.WriteString(data)
	_ = f.Close()
	return nil
}

func (u *Utils) CopyFileAndReplaceLabelsFromHost(name string, labels []interfaces.ReplaceLabel, object *interfaces.FileObject) error {
	dat, err := ioutil.ReadFile(object.Name)
	templateFile := string(dat)
	f, err := os.Create(name)
	if err != nil {
		log.Println("create err", err)
	}
	data := templateFile
	for _, label := range labels {
		data = strings.ReplaceAll(data, label.Label, label.Value)
	}
	f.WriteString(data)
	_ = f.Close()

	return nil
}

func (u *Utils) CopyDirAndReplaceLabel(name string, label *interfaces.ReplaceLabel,object *interfaces.FileObject) error {
	box := packr.NewBox(object.Path)
	errC := os.MkdirAll(name, 0755)
	if errC != nil{
		color.Red("%s",errC)
	}
	box.Walk(func(s string, file packd.File) error {
		//color.Yellow("s: %s, file: %+f \n", s, file)
		pathDir := filepath.Dir(s)
		//color.Yellow("data : %s, info: %s",s, pathDir)
		if pathDir != "." {
			errD := os.Mkdir(name+"/"+pathDir, 0755)
			if errD != nil{
				//color.Red("%s",errD)
			}

		}

		f, err := os.Create(name+"/"+s)
		if err != nil {
			log.Println("create err", err)
		}
		data := file.String()
		data = strings.ReplaceAll(data, label.Label, label.Value)
		f.WriteString(data)
		_ = f.Close()

		return nil
	})
	return nil
}

func (u *Utils) GetGoTemplate(object *interfaces.FileObject) (error error, data string)  {
	box := packr.NewBox(object.Path)

	data, err := box.FindString(object.Name)

	if err != nil {
		log.Fatal(err)
	}

	return err, data
}

func (u *Utils) CheckSlice(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func (u *Utils) RemoveFiles(paths []string)  {
	for _, path := range paths {
		_ = os.RemoveAll(path)
	}
}

func (u *Utils) RenderTableOutput(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.AppendBulk(data)
	table.Render()
}


func (u *Utils) print(rd io.Reader) error {
	var lastLine string

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		lastLine = scanner.Text()
		color.Yellow("%s", lastLine)
	}

	errLine := &ErrorLine{}
	json.Unmarshal([]byte(lastLine), errLine)
	if errLine.Error != "" {
		return errors.New(errLine.Error)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

