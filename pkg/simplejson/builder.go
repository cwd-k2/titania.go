package simplejson

import (
	"bufio"
	"strconv"
	"strings"
)

// hidden!
type builder struct {
	buf *bufio.Writer
	end byte
	fst bool
	// TODO: structure management (ok to set, add, etc)
	// (state: inside object/nestedobject/array/nestedarray...)
	// TODO: Error handling
}

func (b *builder) Flush() {
	b.buf.WriteByte(b.end)
	b.buf.Flush()
}

func (b *builder) addComma() {
	if !b.fst {
		b.buf.WriteByte(',')
	}
	b.fst = false
}

func (b *builder) addKey(key string) {
	b.addComma()
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
	b.buf.WriteString(strconv.Itoa(value))
}

func (b *builder) SetBool(key string, value bool) {
	b.addKey(key)
	if value {
		b.buf.WriteString("true")
	} else {
		b.buf.WriteString("false")
	}
}

func (b *builder) SetNull(key string) {
	b.addKey(key)
	b.buf.WriteString("null")
}

func (b *builder) SetString(key, value string) {
	b.addKey(key)
	b.writeEscapedString(value)
}

/**
 * Add values into array
 */

func (b *builder) AddInt(value int) {
	b.addComma()
	b.buf.WriteString(strconv.Itoa(value))
}

func (b *builder) AddBool(value bool) {
	b.addComma()
	if value {
		b.buf.WriteString("true")
	} else {
		b.buf.WriteString("false")
	}
}

func (b *builder) AddNull() {
	b.addComma()
	b.buf.WriteString("null")
}

func (b *builder) AddString(value string) {
	b.addComma()
	b.writeEscapedString(value)
}

/**
 * nested structures
 */

func (b *builder) SetObject(key string, value func(Object)) {
	b.addKey(key)
	b.fst = true
	b.buf.WriteByte('{')
	value(b)
	b.buf.WriteByte('}')
}

func (b *builder) SetArray(key string, value func(Array)) {
	b.addKey(key)
	b.fst = true
	b.buf.WriteByte('[')
	value(b)
	b.buf.WriteByte(']')
}

func (b *builder) AddObject(value func(Object)) {
	b.addComma()
	b.fst = true
	b.buf.WriteByte('{')
	value(b)
	b.buf.WriteByte('}')
}

func (b *builder) AddArray(value func(Array)) {
	b.addComma()
	b.fst = true
	b.buf.WriteByte('[')
	value(b)
	b.buf.WriteByte(']')
}
