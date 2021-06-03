package simplejson

import (
	"bufio"
	"io"
)

type Array interface {
	AddInt(int)
	AddBool(bool)
	AddNull()
	AddString(string)
	AddObject(func(Object))
	AddArray(func(Array))
}

type ArrayBuilder interface {
	// Once built, you shouldn't touch this.
	Flush()
	Array
}

func NewArrayBuilder(w io.Writer) ArrayBuilder {
	buf := bufio.NewWriter(w)
	buf.WriteByte('[')

	return &builder{
		buf: buf,
		end: ']',
		fst: true,
	}
}
