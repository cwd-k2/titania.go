package simplejson

import (
	"bytes"
	"io"
	"strconv"
	"strings"
)

// hidden!
type builder struct {
	buf *bytes.Buffer
	end byte
	// TODO: structure management (ok to set, add, etc)
	// (state: inside object/nestedobject/array/nestedarray...)
	// TODO: Error handling
}

func (b *builder) Build() io.Reader {
	b.trailing(b.end)

	buf := b.buf
	b.buf = nil

	return buf
}

func (b *builder) trailing(token byte) {
	length := b.buf.Len()
	if length > 0 && b.buf.Bytes()[length-1] == ',' {
		b.buf.Bytes()[length-1] = token
	} else {
		b.buf.WriteByte(token)
	}
}

func (b *builder) addKey(key string) {
	b.writeEscapedString(key)
	b.buf.WriteByte(':')
}

func (b *builder) writeEscapedString(str string) {
	b.buf.WriteString(strings.ReplaceAll(strconv.Quote(str), `\x`, `\u00`))
}

/**
 * Set keys and values into object
 */

func (b *builder) SetInt(key string, value int) {
	b.addKey(key)
	b.AddInt(value)
}

func (b *builder) SetBool(key string, value bool) {
	b.addKey(key)
	b.AddBool(value)
}

func (b *builder) SetNull(key string) {
	b.addKey(key)
	b.AddNull()
}

func (b *builder) SetString(key, value string) {
	b.addKey(key)
	b.AddString(value)
}

/**
 * Add values into array
 */

func (b *builder) AddInt(value int) {
	b.buf.WriteString(strconv.Itoa(value))
	b.buf.WriteByte(',')
}

func (b *builder) AddBool(value bool) {
	if value {
		b.buf.WriteString("true")
	} else {
		b.buf.WriteString("false")
	}
	b.buf.WriteByte(',')
}

func (b *builder) AddNull() {
	b.buf.WriteString("null")
	b.buf.WriteByte(',')
}

func (b *builder) AddString(value string) {
	b.writeEscapedString(value)
	b.buf.WriteByte(',')
}

/**
 * nested structures
 */

func (b *builder) SetObject(key string, value func(Object)) {
	b.addKey(key)
	b.AddObject(value)
}

func (b *builder) SetArray(key string, value func(Array)) {
	b.addKey(key)
	b.AddArray(value)
}

func (b *builder) AddObject(value func(Object)) {
	b.buf.WriteByte('{')
	value(b)
	b.trailing('}')
	b.buf.WriteByte(',')
}

func (b *builder) AddArray(value func(Array)) {
	b.buf.WriteByte('[')
	value(b)
	b.trailing(']')
	b.buf.WriteByte(',')
}
