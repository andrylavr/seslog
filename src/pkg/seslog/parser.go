package seslog

import (
	"errors"
	"fmt"
	"strings"
)

type parsedLine map[string]string

type NginxLogParser struct {
	varnames   []string
	notEnougth error
}

func NewNginxLogParser(format string) *NginxLogParser {
	parser := &NginxLogParser{}

	format = strings.Replace(format, "$", "", -1)
	parser.varnames = strings.Split(format, "\t")
	parser.notEnougth = errors.New(fmt.Sprintf("Not enought variables in line (need %d)", len(parser.varnames)))

	return parser
}

func (this *NginxLogParser) parseString(line string) (parsedLine, error) {
	result := make(parsedLine)

	vars := strings.Split(line, "\t")
	if len(vars) != len(this.varnames) {
		return nil, this.notEnougth
	}
	for num, varname := range this.varnames {
		result[varname] = vars[num]
	}
	return result, nil
}
