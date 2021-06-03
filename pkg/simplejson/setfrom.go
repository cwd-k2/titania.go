package simplejson

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

/**
 * file contents!
 */

// TODO: perfomance?
func (b *builder) writeEscapedStringFromBuffer(strbuf *bytes.Buffer) error {
	//b.buf.WriteString(strings.ReplaceAll(strconv.Quote(strbuf.String()), `\x`, `\u00`))
	b.buf.WriteByte('"')
	for {
		c, err := strbuf.ReadByte()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if c >= 0x20 && c != '\\' && c != '"' {
			b.buf.WriteByte(c)
			continue
		}
		switch c {
		case '\\', '"':
			b.buf.WriteByte('\\')
			b.buf.WriteByte(c)
		case '\n':
			b.buf.WriteByte('\\')
			b.buf.WriteByte('n')
		case '\f':
			b.buf.WriteByte('\\')
			b.buf.WriteByte('f')
		case '\b':
			b.buf.WriteByte('\\')
			b.buf.WriteByte('b')
		case '\r':
			b.buf.WriteByte('\\')
			b.buf.WriteByte('r')
		case '\t':
			b.buf.WriteByte('\\')
			b.buf.WriteByte('t')
		default:
			b.buf.WriteString(`\u00`)
			b.buf.WriteByte("0123456789abcdef"[c>>4])
			b.buf.WriteByte("0123456789abcdef"[c&0xF])
		}
	}
	b.buf.WriteByte('"')
	return nil
}

func (b *builder) SetStringFromFile(key, filename string) error {
	b.addKey(key)

	strbuf := bytes.NewBuffer([]byte{})

	fp, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fp.Close()

	bufio.NewReader(fp).WriteTo(strbuf)

	b.writeEscapedStringFromBuffer(strbuf)
	b.buf.WriteByte(',')

	return nil
}

func (b *builder) SetStringFromFiles(key string, filenames []string, delimiter string) error {
	b.addKey(key)

	strbuf := bytes.NewBuffer([]byte{})

	for i, filename := range filenames {
		fp, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer fp.Close()

		bufio.NewReader(fp).WriteTo(strbuf)

		if i < len(filenames)-1 {
			strbuf.WriteString(delimiter)
		}
	}

	b.writeEscapedStringFromBuffer(strbuf)
	b.buf.WriteByte(',')

	return nil
}

func (b *builder) SetStringFromReader(key string, reader io.Reader) {
	b.addKey(key)

	strbuf := bytes.NewBuffer([]byte{})
	strbuf.ReadFrom(reader)

	b.writeEscapedStringFromBuffer(strbuf)
	b.buf.WriteByte(',')
}

func (b *builder) SetStringFromReaders(key string, readers []io.Reader, delimiter string) {
	b.addKey(key)

	strbuf := bytes.NewBuffer([]byte{})

	for i, reader := range readers {
		strbuf.ReadFrom(reader)
		if i < len(readers)-1 {
			strbuf.WriteString(delimiter)
		}
	}

	b.writeEscapedStringFromBuffer(strbuf)
	b.buf.WriteByte(',')
}
