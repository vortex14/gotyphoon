package utils

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"path/filepath"
	"strings"
)

func UnCommentCode(marker string, code string) string {
	lines := strings.Split(code, "\n")
	var unCommentedLineList []string
	for _, line := range lines {
		if strings.Contains(line, marker) { unCommentedLineList = append(unCommentedLineList, line);continue }
		if strings.Contains(line, "*/") { unCommentedLineList = append(unCommentedLineList, line); continue}
		uncommentLine := strings.Replace(line, "//", "", 1)
		unCommentedLineList = append(unCommentedLineList, uncommentLine)

	}
	unCommentedLines := strings.Join(unCommentedLineList, "\n")
	return unCommentedLines
}

func CommentCode(marker string, code string) string {
	lines := strings.Split(code, "\n")
	//var stopMarker bool
	var CommentedLineList []string

	isComment := false

	for _, line := range lines {
		if IsStrContain(line, "package") { CommentedLineList = append(CommentedLineList, line); continue }

		firstSliceStr := ""

		if len(line) >= 5 { firstSliceStr = line[:5] }

		if strings.Contains(line, " */") {
			isComment = false
			CommentedLineList = append(CommentedLineList, line)
		} else if strings. Contains(line, marker) {
			isComment = true
			CommentedLineList = append(CommentedLineList, line)
		} else if strings.Contains(line, fmt.Sprintf("// %s", marker)) {
			isComment = true
			CommentedLineList = append(CommentedLineList, line)
		} else if strings.Contains(firstSliceStr, "//") {
			CommentedLineList = append(CommentedLineList, line)

		} else if isComment || strings.Contains(line, "//") {
			commentLine := fmt.Sprintf("//%s", line)
			CommentedLineList = append(CommentedLineList, commentLine)
		} else {
			CommentedLineList = append(CommentedLineList, line)
		}

		//else if strings.Contains(line, marker) {
		//	isComment = true
		//	CommentedLineList = append(CommentedLineList, line)
		//} else if strings.Contains(line, "*/") {
		//	isComment = false
		//	CommentedLineList = append(CommentedLineList, line)
		//} else {
		//	isComment = false
		//}



		//if firstLine || !stopMarker && isComment {
		//	commentLine := fmt.Sprintf("//%s", line)
		//	CommentedLineList = append(CommentedLineList, commentLine)
		//} else if !firstLine {
		//	if strings.Contains(line, marker) { stopMarker = false; isComment = true }
		//	CommentedLineList = append(CommentedLineList, line)
		//} else if strings.Contains(line, marker) { stopMarker = false; isComment = true; continue }
		//println(line)
		//firstLine = false




	}
	CommentedLines := strings.Join(CommentedLineList, "\n")
	return CommentedLines
}

func UncommentDir(startDir string, matchCode string, excludeDirs map[string]bool)  {
	_ = filepath.Walk(startDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil { return err }
			firstDir := GetFirstDir(path)
			if _, ok := excludeDirs[firstDir]; ok { return nil }
			if info.IsDir() { return nil}
			contentFileCode := ReadFile(path)
			marker := fmt.Sprintf("/* %s", matchCode)
			if strings.Contains(contentFileCode, marker) {
				//println(contentFileCode, matchCode)
				unCommentCode := UnCommentCode(marker, contentFileCode)
				errUn := SaveData(path, unCommentCode)
				if errUn != nil { color.Red(errUn.Error()) }
			}
			return nil
		})
}

func CommentDir(startDir string, matchCode string, excludeDirs map[string]bool)  {
	println("CommentDir ... ")
	_ = filepath.Walk(startDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil { return err }
			firstDir := GetFirstDir(path)
			if _, ok := excludeDirs[firstDir]; ok { return nil }
			if info.IsDir() { return nil}
			contentFileCode := ReadFile(path)
			marker := fmt.Sprintf("/* %s", matchCode)
			if strings.Contains(contentFileCode, marker) {
				//println(contentFileCode, matchCode)
				commentedCode := CommentCode(marker, contentFileCode)
				//println(commentedCode)
				errUn := SaveData(path, commentedCode)
				if errUn != nil { color.Red(errUn.Error()) }
			}
			return nil
		})
}
