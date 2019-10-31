package files

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

//判断文件是否存在
func isExist(path string) (bool) {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	} else if err != nil {
		panic(err)
	}
	return true
}

// 创建多级目录
func mkdirall(path string) {
	e := os.MkdirAll(path, os.ModePerm)
	if e != nil {
		panic(e)
	}
}

// 要把文件名(带绝对路径)创建文件名
func createFile(path string) *os.File {
	dirpath, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		panic(err)
	}
	if !isExist(dirpath) {
		mkdirall(dirpath)
	}
	if isExist(path) {
		if err := os.Remove(path); err != nil {
			panic(err)
		}
	}
	fileObj, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	return fileObj
}

//name:文件名(带绝对路径)
func write(name string, info map[string][]string) error {
	// 如果存在name文件，则会先删除后再新建该文件
	logFile := createFile(name)
	defer logFile.Close()
	writer := bufio.NewWriter(logFile)
	for name, content := range info {
		var err1, err2 error
		if len(content) > 0 {
			_, err1 = writer.WriteString("### " + name + "\n")
			_, err2 = writer.WriteString(strings.Join(content, "\n"))
		} else {
			_, err1 = writer.WriteString("\n## " + name)
		}
		if err1 != nil {
			return err1
		}
		if err2 != nil {
			return err2
		}
		writer.WriteString("\n")
	}
	writer.Flush()
	return nil
}
