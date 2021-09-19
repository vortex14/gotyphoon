package utils

import (
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"os"
)

func ReadFile(path string) string {
	at, err := ioutil.ReadFile(path)
	if err != nil {
		color.Red(err.Error())
	}
	return string(at)
}

func SaveData(path string, data string) error {

		f, err := os.Create(path)
		if err != nil {
			log.Println("create err", err)
			return err
		}
		_, errorWrite := f.WriteString(data)
		if errorWrite != nil {
			color.Red("Can't write %s", data )
			return errorWrite

		}
		_ = f.Close()

		return nil

}