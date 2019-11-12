package main

import (
	"flag"
	"fmt"
	"github.com/coolligang/go-test-report-parser/files"
	"github.com/jstemmer/go-junit-report/parser"
	"github.com/xiaosongluo/go-test-report-parser/formatter"
	"os"
	"strings"
)

var (
	packageName   string
	formatterName string
	goVersionFlag string
	setExitCode   bool
	logPath       string
	errSelect     bool
	csvReport     bool
)

func init() {
	flag.StringVar(&packageName, "package-name", "", "specify a package name (compiled test have no package name in output)")
	flag.StringVar(&formatterName, "formatter-name", "MarkdownFunctionFormatter", "specify a formatter name")
	flag.StringVar(&goVersionFlag, "go-version", "", "specify the value to use for the go.version property in the generated XML")
	flag.BoolVar(&setExitCode, "set-exit-code", false, "set exit code to 1 if tests failed")
	flag.StringVar(&logPath, "logs", "", "Absolute path of log")
	flag.BoolVar(&errSelect, "err", false, "Output error logs to files separately")
	flag.BoolVar(&csvReport, "csv", true, "Output CSV report")
}

//检查路径path路径是否以"/"结尾，如果不是则加上
func formatPath(path string) string {
	if !(strings.HasSuffix(path, "/") || strings.HasSuffix(path, "\\")) {
		path += "/"
	}
	path = strings.Replace(path, "\\", "/", -1)
	return path
}

func main() {
	flag.Parse()

	if flag.NArg() != 0 {
		fmt.Println("go-junit-report does not accept positional arguments")
		os.Exit(1)
	}

	// Read input
	report, err := parser.Parse(os.Stdin, packageName)
	if err != nil {
		fmt.Printf("Error reading input: %s\n", err)
		os.Exit(1)
	}

	// Output
	output := formatter.GetAllFormatter()[formatterName]
	if output != nil {
		err = output.Formatter(report, os.Stdout)
		if err != nil {
			fmt.Printf("Error Output: %s\n", err)
			os.Exit(1)
		}
	}

	// output logfile
	//如果logs
	if logPath != "" {
		logPath = formatPath(logPath)
		if err := files.Outputall(report, logPath); err != nil {
			fmt.Printf("OutputError: %s", err)
			os.Exit(1)
		}
	} else {
		path, err := os.Getwd()
		if err != nil {
			fmt.Printf("Get current path error: %s", err)
			os.Exit(1)
		}
		logPath = path + "/logs/"
	}

	if errSelect {
		if err := files.OutputError(report, logPath); err != nil {
			fmt.Printf("OutputError: %s", err)
			os.Exit(1)
		}
	}

	if csvReport {
		if err:=files.ReportCSV(report, logPath);err!=nil{
			fmt.Printf("ReportCSV: %s", err)
			os.Exit(1)
		}
	}

	if setExitCode && report.Failures() > 0 {
		os.Exit(1)
	}
}
