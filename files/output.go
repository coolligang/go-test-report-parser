package files

import (
	"fmt"
	"github.com/jstemmer/go-junit-report/parser"
	"regexp"
	"strings"
)

// 在line中查找符合str规则的字符串，如果有匹配项返回true,否则false
func isRegular(line string, str string) bool {
	reg := regexp.MustCompile(str)
	return reg.FindString(string(line)) != ""
}

func formatResult(test *parser.Test) string {
	var result string
	switch test.Result {
	case parser.PASS:
		result = "PASS"
	case parser.FAIL:
		result = "FAIL"
	default:
		result = "SKIP"
	}
	return result
}

//根据输出，获取开始时间
func getStartTime(output []string) string {
	if len(output) == 0 {
		return ""
	}
	starttime := ""
	for _, line := range output {
		line = strings.Replace(line, " ", "", -1)
		line = strings.Replace(line, "\t", "", -1)
		if isRegular(line, "^Date:.*?GMT") {
			starttime = line[5 : len(line)-3]
			break
		}
	}
	return strings.Replace(starttime, ":", "", -1)
}

//格式化文件名称，替换字符串中的 "/" "\" 为 "_"
func formatname(name string) string {
	name = strings.Replace(name, "/", "-", -1)
	return strings.Replace(name, "\\", "-", -1)
}

//组装用例头
func getHeader(test *parser.Test) string {
	result := formatResult(test)
	return fmt.Sprintf("%s %s duration: %s", result, test.Name, test.Duration)
}

// 将所有日志以package为单位输出到文件
func Outputall(report *parser.Report, path string) error {
	for _, pkg := range report.Packages {
		info := make(map[string][]string)
		finename := formatname(pkg.Name)
		for _, test := range pkg.Tests {
			header := getHeader(test)
			info[header] = test.Output
		}
		if err := write(path+finename+".md", info); err != nil {
			return err
		}
	}
	return nil
}

//将错误用例的日志分别输出到不同的文件
func OutputError(report *parser.Report, path string) error {
	for _, pkg := range report.Packages {
		info := make(map[string][]string)
		for _, test := range pkg.Tests {
			if test.Result == parser.FAIL && len(test.Output) > 0 {
				header := getHeader(test)
				startTime := getStartTime(test.Output)
				info[header] = test.Output
				if err := write(path+"errors/"+test.Name+"_"+startTime+".md", info); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
