package utils

import (
	"log"
	"os"
	"strings"
)

func GetCurrentDir() string {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	return path
}

func GetFirstDir(path string) string {
	paths := strings.Split(path, "/")
	for _, route := range paths {
		if strings.Contains(route, "..") { continue }
		return route
	}
	return ""
}


// Walker. ignoreDirs, ignoreFiles, matchFileExtensions,
// return: file, content,
func Walker()  {
	
}