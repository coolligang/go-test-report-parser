package files

import (
	"bufio"
	"encoding/csv"
	"github.com/jstemmer/go-junit-report/parser"
	"os"
	"path/filepath"
	"strings"
	"fmt"
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

func writecsv(name string, pkg parser.Package) error {
	csvFile := createFile(name)
	defer csvFile.Close()
	csvFile.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	writer := csv.NewWriter(csvFile)
	// 写入文件头，及用例统计信息
	writer.Write([]string{"测试模块", "总用例数", "通过", "失败", "未执行"})
	pass, fail, skip := getTestCount(pkg)
	writer.Write([]string{pkg.Name, fmt.Sprintf("%d", pass+fail+skip), fmt.Sprintf("%d", pass), fmt.Sprintf("%d", fail), fmt.Sprintf("%d", skip)})
	writer.Write([]string{})
	writer.Flush()

	// 写入测试概要
	writer.Write([]string{"全部测试用例结果概要"})
	writer.Write([]string{"TESTNAME", "RESULT"})
	for _, test := range pkg.Tests {
		writer.Write([]string{test.Name, formatResult(test)})
	}
	writer.Write([]string{})
	writer.Flush()

	//写入测试用例错误详细
	writer.Write([]string{"错误测试用例详细"})
	writer.Write([]string{"TESTNAME", "RESULT", "OUTPUT"})
	for _, errtest := range pkg.Tests {
		if errtest.Result == parser.FAIL && len(errtest.Output) > 0 {
			output := strings.Join(errtest.Output, "\n")
			writer.Write([]string{errtest.Name, formatResult(errtest), output})
		}
	}
	return nil
}

func getTestCount(pkg parser.Package) (int, int, int) {
	var pass, fail, skip int
	for _, test := range pkg.Tests {
		if len(test.Output) > 0 {
			switch test.Result {
			case parser.PASS:
				pass ++
			case parser.FAIL:
				fail ++
			default:
				skip ++
			}
		}
	}
	return pass, fail, skip
}
