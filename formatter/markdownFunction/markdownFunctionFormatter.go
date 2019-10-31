package markdownFunction

import (
	"bufio"
	"fmt"
	"github.com/jstemmer/go-junit-report/parser"
	"github.com/xiaosongluo/go-test-report-parser/formatter"
	"io"
)

func init() {
	item := new(MarkdownFunctionFormatter)
	formatter.RegisterFormatter(item)
}

type MarkdownFunctionFormatter struct {
}

func (formatter MarkdownFunctionFormatter) GetName() string {
	return "MarkdownFunctionFormatter"
}

func (formatter MarkdownFunctionFormatter) Formatter(report *parser.Report, w io.Writer) error {
	writer := bufio.NewWriter(w)
	summary := ""
	detail := ""
	error_case := ""

	// convert Report to JUnit test suites
	for _, pkg := range report.Packages {
		classname := pkg.Name

		pass := 0
		fail := 0
		skip := 0

		// individual test cases
		for _, test := range pkg.Tests {

			var result string
			switch test.Result {
			case parser.PASS:
				result = "通过"
				pass += 1
			case parser.FAIL:
				result = "不通过"
				fail += 1
			default:
				result = "未测试"
				skip += 1
			}
			item := fmt.Sprintf("|API测试|%s|%s|%s|\n", classname, test.Name, result)
			detail = detail + item
			if result == "不通过" {
				error_case = error_case + item
			}
		}

		item := fmt.Sprintf("|%s|%d|%d|%d|%d|\n", classname, pass+fail+skip, pass, fail, skip)
		summary = summary + item
	}
	_, _ = writer.WriteString(summaryHeader())
	_, _ = writer.WriteString(summary)
	_ = writer.WriteByte('\n')

	_, _ = writer.WriteString(detailHeader())
	_, _ = writer.WriteString(detail)
	_ = writer.WriteByte('\n')
	_ = writer.Flush()

	_, _ = writer.WriteString(errorHeader())
	_, _ = writer.WriteString(error_case)
	_ = writer.WriteByte('\n')
	_ = writer.Flush()

	return nil
}

func detailHeader() string {
	return "###详细测试结果\n|测试内容|测试模块|子测试项|测试结果|\n|--------|---------|--------|--------|\n"
}

func summaryHeader() string {
	return "###汇总测试结果\n|测试模块|总用例数|通过用例数|未通过用例数|跳过用例数|\n|--------|---------|--------|--------|--------|\n"
}

func errorHeader() string {
	return "###失败的详细测试结果\n|测试内容|测试模块|子测试项|测试结果|\n|--------|---------|--------|--------|\n"
}

