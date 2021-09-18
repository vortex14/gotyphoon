package utils

import (
	"log"
	"os"
)

func GetCurrentDir() string {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	return path
}
