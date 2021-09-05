package utils

import (
	"encoding/json"
	"fmt"
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
	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr"
	"github.com/gocarina/gocsv"
	"github.com/google/uuid"
	"github.com/olekukonko/tablewriter"
	"github.com/vortex14/gotyphoon/interfaces"
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

func (u *Utils) CopyDir(name string, object *interfaces.FileObject) error {
	box := packr.NewBox(object.Path)
	errC := os.MkdirAll(name, 0755)
	if errC != nil{
		color.Red("CopyDir Error : %s",errC)
	}
	err := box.Walk(func(s string, file packd.File) error {
		//color.Yellow("s: %s, file: %+f \n", s, file)
		pathDir := filepath.Dir(s)
		//color.Yellow("data : %s, info: %s",s, pathDir)
		if pathDir != "." {
			_ = os.Mkdir(name+"/"+pathDir, 0755)

		}

		f, err := os.Create(name + "/" + s)
		if err != nil {
			log.Println("create err", err)
		}
		_, err = f.WriteString(file.String())
		if err != nil {
			color.Red("%s", err.Error())
			return err
		}
		_ = f.Close()

		return nil
	})
	return err
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

func (u *Utils) CopyDirAndReplaceLabel(name string, label *interfaces.ReplaceLabel,object *interfaces.FileObject) error {
	box := packr.NewBox(object.Path)
	errC := os.MkdirAll(name, 0755)
	if errC != nil{
		color.Red("%s",errC)
	}
	err := box.Walk(func(s string, file packd.File) error {
		//color.Yellow("s: %s, file: %+f \n", s, file)
		pathDir := filepath.Dir(s)
		//color.Yellow("data : %s, info: %s",s, pathDir)
		if pathDir != "." {
			_ = os.Mkdir(name+"/"+pathDir, 0755)

		}

		f, err := os.Create(name + "/" + s)
		defer func(f *os.File) {
			errD := f.Close()
			if errD != nil {
				color.Red("%s", errD.Error())
			}
		}(f)
		if err != nil {
			log.Println("create err", err)
		}
		data := file.String()
		data = strings.ReplaceAll(data, label.Label, label.Value)
		_, err = f.WriteString(data)
		if err != nil {
			return err
		}

		return nil
	})
	return err
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

func (u *Utils) ReadCSV(object *interfaces.FileObject, bindings interface{})  {
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
	return rand.Intn(max - min) + min

}

func (u *Utils) GetRandomFromSlice(slice []string) string {
	rand.Seed(time.Now().UnixNano())
	return slice[rand.Intn(len(slice))]
}

func (u *Utils) GetRandomString(length int, sequence string)  string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune(sequence)
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func (u *Utils) GetUUID() string {
	return uuid.NewString()
}

func (u *Utils) GetRandomFloat() float64 {

	rand.Seed(time.Now().UnixNano())

	roundValue := math.Floor(rand.Float64() * 10000) / 100
	return roundValue
}

func (u *Utils) ConvertStringListToIntList(input []string) []int {
	var output []int

	for _, i := range input {
		j, err := strconv.Atoi(i)
		if err != nil {
			color.Red("%s",err.Error())
			continue
		}
		output = append(output, j)
	}
	return output
}
