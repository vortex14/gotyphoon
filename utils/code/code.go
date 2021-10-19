package code

import (
	"fmt"
	"os"
	"strings"

	"path/filepath"

	"github.com/fatih/color"

	"github.com/vortex14/gotyphoon/utils"
)

func UnCommentCode(marker string, code string) string {
	lines := strings.Split(code, "\n")
	var unCommentedLineList []string
	isUncomment := false
	for _, line := range lines {

		if strings.Contains(line, marker) {
			unCommentedLineList = append(unCommentedLineList, line)
			isUncomment = true
			continue
		} else if !isUncomment || len(line) <= 1 {
			unCommentedLineList = append(unCommentedLineList, line)
			continue
		} else if len(line) <= 3 && strings.Contains(line, "//") && isUncomment {
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
		} else {
			unCommentedLineList = append(unCommentedLineList, line)
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
		if utils.IsStrContain(line, "package") {
			CommentedLineList = append(CommentedLineList, line)
			continue
		}

		firstSliceStr := ""

		if len(line) >= 2 {
			firstSliceStr = line[:2]
		}

		if strings.Contains(line, "*/") {
			isComment = false
			CommentedLineList = append(CommentedLineList, line)
		} else if strings.Contains(line, marker) {
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

func getSourceCode(
	path string,
	matchCode string,
	info os.FileInfo,
	excludeDirs map[string]bool,
	callback func(marker string, source string)) error {

	firstDir := utils.GetFirstDir(path)

	switch {
	case info.IsDir() || strings.Contains(path, "_test.go") || excludeDirs[firstDir]:
		return nil
	}

	contentFileCode := utils.ReadFile(path)
	marker := fmt.Sprintf("/* %s", matchCode)
	if strings.Contains(contentFileCode, marker) {
		callback(marker, contentFileCode)
	}
	return nil
}

func UncommentDir(startDir string, matchCode string, excludeDirs map[string]bool)  {
	_ = filepath.Walk(startDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil { return err }
			return getSourceCode(path, matchCode, info, excludeDirs, func(marker string, source string) {
				unCommentCode := UnCommentCode(marker, source)
				errUn := utils.SaveData(path, unCommentCode)
				if errUn != nil { color.Red(errUn.Error()) }
			})
		})
}

func CommentDir(startDir string, matchCode string, excludeDirs map[string]bool)  {
	println("CommentDir ... ")
	_ = filepath.Walk(startDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil { return err }
			return getSourceCode(path, matchCode, info, excludeDirs, func(marker string, source string) {
				commentedCode := CommentCode(marker, source)
				errUn := utils.SaveData(path, commentedCode)
				if errUn != nil { color.Red(errUn.Error()) }
			})
		})
}
