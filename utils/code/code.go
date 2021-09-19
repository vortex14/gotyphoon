package code

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/vortex14/gotyphoon/utils"
	"os"
	"path/filepath"
	"strings"
)

func UnCommentCode(marker string, code string) string {
	lines := strings.Split(code, "\n")
	var unCommentedLineList []string
	isUncomment := false
	for _, line := range lines {

		if strings.Contains(line, marker){
			unCommentedLineList = append(unCommentedLineList, line)
			isUncomment = true
			continue
		} else if !isUncomment {
			unCommentedLineList = append(unCommentedLineList, line)
			continue
		} else if len(line) <= 1 {
			unCommentedLineList = append(unCommentedLineList, line)
			continue
		} else if len(line) <= 3 && strings.Contains(line, "//") && isUncomment{
			nline := strings.ReplaceAll(line, "//", "")
			unCommentedLineList = append(unCommentedLineList, nline)
			continue
		} else if strings.Contains(line, "*/") {
			unCommentedLineList = append(unCommentedLineList, line)
			isUncomment = false
			continue
		} else if strings.Contains(line, "//") {
			uncommentLine := strings.Replace(line, "//", "", 1)
			unCommentedLineList = append(unCommentedLineList, uncommentLine)
		}

	}
	unCommentedLines := strings.Join(unCommentedLineList, "\n")
	return unCommentedLines
}

func CommentCode(marker string, code string) string {
	lines := strings.Split(code, "\n")
	var CommentedLineList []string

	isComment := false

	for _, line := range lines {
		if utils.IsStrContain(line, "package") { CommentedLineList = append(CommentedLineList, line); continue }

		firstSliceStr := ""

		if len(line) >= 5 { firstSliceStr = line[:5] }

		if strings.Contains(line, "*/") {
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

		} else if isComment {
			commentLine := fmt.Sprintf("//%s", line)
			CommentedLineList = append(CommentedLineList, commentLine)
		} else {
			CommentedLineList = append(CommentedLineList, line)
		}

	}

	CommentedLines := strings.Join(CommentedLineList, "\n")
	return CommentedLines
}

func UncommentDir(startDir string, matchCode string, excludeDirs map[string]bool)  {
	_ = filepath.Walk(startDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil { return err }
			firstDir := utils.GetFirstDir(path)
			if _, ok := excludeDirs[firstDir]; ok { return nil }
			if info.IsDir() { return nil}
			if strings.Contains(path, "_test.go") { return nil}
			contentFileCode := utils.ReadFile(path)
			marker := fmt.Sprintf("/* %s", matchCode)
			if strings.Contains(contentFileCode, marker) {
				unCommentCode := UnCommentCode(marker, contentFileCode)
				errUn := utils.SaveData(path, unCommentCode)
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
			firstDir := utils.GetFirstDir(path)
			if _, ok := excludeDirs[firstDir]; ok { return nil }
			if info.IsDir() { return nil}
			if strings.Contains(path, "_test.go") { return nil}
			contentFileCode := utils.ReadFile(path)
			marker := fmt.Sprintf("/* %s", matchCode)
			if strings.Contains(contentFileCode, marker) {
				commentedCode := CommentCode(marker, contentFileCode)
				errUn := utils.SaveData(path, commentedCode)
				if errUn != nil { color.Red(errUn.Error()) }
			}
			return nil
		})
}
