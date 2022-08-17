package utils

import (
	"encoding/json"
	"fmt"

	"io/fs"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/fatih/color"
	"github.com/go-logfmt/logfmt"
	"github.com/gocarina/gocsv"
	"github.com/google/uuid"
	"github.com/olekukonko/tablewriter"
	"github.com/vortex14/gotyphoon/interfaces"
)

type Utils struct{}

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
	if NotNill(err) {
		log.Println("create file: ", err)
		status = false
	}

	err = tmpl.Execute(f, &goTemplate.Data)
	if NotNill(err) {
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

func (u *Utils) CopyDir(name string, dir fs.FS) error {

	errC := os.MkdirAll(name, 0755)
	if errC != nil {
		color.Red("CopyDir Error : %s", errC)
		panic(errC)
	}

	stat, _ := fs.Stat(dir, ".")

	color.Yellow("copy dir from: %+v", stat.Name())

	err := fs.WalkDir(dir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
			return err
		}

		pathDir := filepath.Dir(path)
		fileInfo, err := d.Info()

		if pathDir != "." {
			_ = os.Mkdir(name+"/"+pathDir, 0755)

		}

		if !fileInfo.IsDir() && !strings.Contains(fileInfo.Name(), "-.tml") {
			var exportPath string

			if fileInfo.Name() == "init.py" {

				exportPath = name + "/" + strings.ReplaceAll(path, "init.py", "__init__.py")
			} else {
				exportPath = name + "/" + path
			}

			f, err := os.Create(exportPath)
			if err != nil {
				log.Println("create err", err)
				panic(err)
			}

			file, err := fs.ReadFile(dir, path)
			if err != nil {

				return err
			}

			_, err = f.WriteString(string(file))
			if err != nil {
				color.Red("%s", err.Error())
				return err
			}
			_ = f.Close()

		}
		return nil
	})
	if err != nil {
		color.Red("%+v", err)
		return err
	}

	return nil
}

func (u *Utils) CopyFile(ExportPath string, object *interfaces.FileObject, dir fs.FS) error {

	f, err := os.Create(ExportPath)
	if err != nil {
		log.Println("create err", err)
	}

	file, err := fs.ReadFile(dir, object.GetPath())
	if err != nil {

		color.Red("%+v", err)
		os.Exit(0)
	}

	color.Red("%s: %d : %s", ExportPath, len(string(file)), string(file))

	_, errorWrite := f.WriteString(string(file))
	if errorWrite != nil {
		color.Red("Can't write %s", ExportPath)
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
		color.Red("Can't write %s", object.Path)
		os.Exit(0)

	}
	_ = f.Close()

	return nil

}

func walkDirEmbed(dir fs.FS) {
	_ = fs.WalkDir(dir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
			return err
		}

		pathDir := filepath.Dir(path)
		fileInfo, err := d.Info()

		println(pathDir, fileInfo.Name())

		return nil
	})
}

func (u *Utils) CopyFileAndReplaceLabel(name string, label *interfaces.ReplaceLabel, object *interfaces.FileObject, dir fs.FS) error {

	f, err := os.Create(name)
	if err != nil {
		log.Println("create err", err)
	}

	file, err := fs.ReadFile(dir, object.GetPath())
	if err != nil {

		color.Red("%+v", err)
		panic(err)
	}

	data := strings.ReplaceAll(string(file), label.Label, label.Value)
	_, err = f.WriteString(data)
	if err != nil {
		color.Red("%s", err.Error())
		return err
	}
	_ = f.Close()
	return nil
}

func (u *Utils) CopyFileAndReplaceLabelsFromHost(name string, labels []interfaces.ReplaceLabel, object *interfaces.FileObject) error {
	dat, err := ioutil.ReadFile(object.Name)
	if err != nil {
		return err
	}
	templateFile := string(dat)
	f, err := os.Create(name)
	defer func(f *os.File) {
		errG := f.Close()
		if errG != nil {
			color.Red("%s", errG.Error())
		}
	}(f)
	if err != nil {
		log.Println("create err", err)
	}
	data := templateFile
	for _, label := range labels {
		data = strings.ReplaceAll(data, label.Label, label.Value)
	}
	_, err = f.WriteString(data)
	if err != nil {
		color.Red("%s", err.Error())

	}
	_ = f.Close()

	return nil
}

func (u *Utils) CopyDirAndReplaceLabel(name string, label *interfaces.ReplaceLabel, object *interfaces.FileObject, dir fs.FS) error {
	errC := os.MkdirAll(name+"/templates", 0755)
	if errC != nil {
		color.Red("%s", errC)
		os.Exit(0)
	}

	err := fs.WalkDir(dir, ".", func(path string, d fs.DirEntry, err error) error {
		//color.Yellow("s: %s, file: %+f \n", s, file)
		pathDir := filepath.Dir(path)
		//color.Yellow("data : %s, info: %s",s, pathDir)

		_ = os.MkdirAll(name+"/"+pathDir, 0755)

		color.Yellow(">>>> %s", name+"/"+pathDir)

		fileInfo, _ := d.Info()

		if !fileInfo.IsDir() {

			f, err := os.Create(name + "/" + path)
			defer func(f *os.File) {
				errD := f.Close()
				if errD != nil {
					color.Red("%s", errD.Error())
				}
			}(f)
			if err != nil {
				log.Println("create err", err)
			}

			file, err := fs.ReadFile(dir, path)
			if err != nil {

				return err
			}
			data := strings.ReplaceAll(string(file), label.Label, label.Value)
			_, err = f.WriteString(data)
			if err != nil {
				return err
			}

		}

		return nil
	})
	return err
}

func (u *Utils) GetGoTemplate(object *interfaces.FileObject, dir fs.FS) (error error, data string) {

	file, err := fs.ReadFile(dir, object.GetPath())
	if err != nil {

		color.Red("%+v", err)
		os.Exit(0)
	}

	return err, string(file)
}

func (u *Utils) CheckSlice(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func (u *Utils) RemoveFiles(paths []string) {
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

//func (u *Utils) print(rd io.Reader) error {
//	var lastLine string
//
//	scanner := bufio.NewScanner(rd)
//	for scanner.Scan() {
//		lastLine = scanner.Text()
//		color.Yellow("%s", lastLine)
//	}
//
//	errLine := &ErrorLine{}
//	err := json.Unmarshal([]byte(lastLine), errLine)
//	if err != nil {
//		return err
//	}
//	if errLine.Error != "" {
//		return errors.New(errLine.Error)
//	}
//
//	if err := scanner.Err(); err != nil {
//		return err
//	}
//
//	return nil
//}

func (u *Utils) PrintPrettyJson(f interface{}) string {
	dump, err := json.MarshalIndent(f, "  ", "  ")
	if err != nil {
		color.Red("%s", err.Error())
	}
	return string(dump)
}

func (u *Utils) ReadCSV(object *interfaces.FileObject, bindings interface{}) {
	csvData, err := os.OpenFile(object.Path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer func(csvData *os.File) {
		_ = csvData.Close()
	}(csvData)

	if err := gocsv.UnmarshalFile(csvData, bindings); err != nil { // Load clients from file
		panic(err)
	}
}

func GetRandomIntRange(max int, min int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min

}

func GetRandomFromSlice(slice []string) string {
	rand.Seed(time.Now().UnixNano())
	return slice[rand.Intn(len(slice))]
}

func (u *Utils) GetRandomFromSlice(slice []string) string {
	rand.Seed(time.Now().UnixNano())
	return slice[rand.Intn(len(slice))]
}

func (u *Utils) GetRandomString(length int, sequence string) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune(sequence)
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func GetUUID() string {
	return uuid.NewString()
}

func (u *Utils) GetUUID() string {
	return uuid.NewString()
}

func (u *Utils) GetRandomFloat() float64 {

	rand.Seed(time.Now().UnixNano())

	roundValue := math.Floor(rand.Float64()*10000) / 100
	return roundValue
}

func (u *Utils) ConvertStringListToIntList(input []string) []int {
	var output []int

	for _, i := range input {
		j, err := strconv.Atoi(i)
		if err != nil {
			color.Red("%s", err.Error())
			continue
		}
		output = append(output, j)
	}
	return output
}
