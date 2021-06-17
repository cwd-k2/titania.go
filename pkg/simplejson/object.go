package simplejson

import (
	"bufio"
	"io"
)

type Object interface {
	SetInt(string, int)
	SetBool(string, bool)
	SetNull(string)
	SetString(string, string)
	SetObject(string, func(Object))
	SetArray(string, func(Array))
	SetStringFromFile(string, string) error
	SetStringFromReader(string, io.Reader)
}

type ObjectBuilder interface {
	// Once built, you shouldn't touch this.
	Flush()
	Object
}

func NewObjectBuilder(w io.Writer) ObjectBuilder {
	buf := bufio.NewWriter(w)
	buf.WriteByte('{')

	return &builder{
		buf: buf,
		end: '}',
		fst: true,
	}
}
