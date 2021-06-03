package simplejson

import (
	"bytes"
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
	Build() io.Reader
	Array
}

func NewArrayBuilder() ArrayBuilder {
	return &builder{
		buf: bytes.NewBuffer([]byte{'['}),
		end: ']',
	}
}
