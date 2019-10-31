package formatter

import (
	"github.com/jstemmer/go-junit-report/parser"
	"io"
)

type IFormatter interface {
	GetName() string
	Formatter(*parser.Report, io.Writer) error
}
