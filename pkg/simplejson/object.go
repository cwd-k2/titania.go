package simplejson

import (
	"bytes"
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
	SetStringFromFiles(string, []string, string) error
	SetStringFromReader(string, io.Reader)
	SetStringFromReaders(string, []io.Reader, string)
}

type ObjectBuilder interface {
	// Once built, you shouldn't touch this.
	Build() io.Reader
	Object
}

func NewObjectBuilder() ObjectBuilder {
	return &builder{
		buf: bytes.NewBuffer([]byte{'{'}),
		end: '}',
	}
}
