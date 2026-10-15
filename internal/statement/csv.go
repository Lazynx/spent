package statement

import "io"

type CSVParser struct{}

func (CSVParser) Parse(r io.Reader) (Statement, error) {
	return Statement{}, errNotImplemented
}
