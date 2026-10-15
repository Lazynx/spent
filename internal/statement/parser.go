package statement

import "io"

type Parser interface {
	Parse(r io.Reader) (Statement, error)
}
