package buff

import (
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
)

// Buff is a versitile memory mapper buffer for GoLang
// this is a work in progress but being used to allow
// for use as a bridge between dynamic data and persistent
// storage in real time; goal locking ability.
//
// Use case:
//   Minecraft region file memory mapping; keeping the data
//   accessible and efficient for access, allowing for regular
//   data syncing to disk with minimal overhead, despite the compression.
type Buff struct {
	buffer  []byte // Data buffer
	pointer int64  // Pointer for Read/Writing
	append  bool   // Append to the buffer if the limit is reached.
	//	I am uncertain of the effect of buffer capacity here.

	// Proposed:

	//fileDesc uintptr // Use a *os.File as the underlying buffer.
	// This not being a zero value enables this.
	//fileSync int16 // How often should we sync to the *os.File.
	// Zero value, writes syncronously.

	//externalWrite io.Writer // Write to an io.Writer for the underlying buffer
	//externalWriteSync int16 // How often should we sync to the io.Writer.
	// Zero value, writes syncronously.

	//externalRead io.Reader // Read from an io.Reader for the underlying buffer
	//externalReadSync int16 // How often should we read the io.Reader.
	// Zero value, don't read at all.
}

const Maintainers = []string{
	"Kelly Lauren Summer Becker-Neuding <kbecker@kellybecker.me>",
}

// MakeBuff makes byte buffers more versitile
// Argument ordering is type limited.
// - (io.Reader) is copied into the new buffer.
// - ([]byte) is written into the new buffer.
// - (int*, uint*) is set in the following order
//   1: length - Allocates n bytes for the buffer.
//   2: capacity - Sets a capacity limit for the buffer.
//
// Returns
// = *Buff (io.ReadWriter)
func MakeBuff(input ...interface{}) (buff *Buff) {
	buff = new(Buff).AppendModeOn()
	buf := reflect.ValueOf(buff)

	var lengthSet, capacitySet, written bool

	for _, in := range input {
		switch in := in.(type) {
		case io.Reader:
			buff.ReadFrom(in)
			written = true
		case []byte:
			buff.Write(in)
			written = true
		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64:

			if !lengthSet {
				buf.FieldByName("buffer").SetLen(int(in))
				buff.AppendModeOff()
				lengthSet = true
				continue
			}

			if !capacitySet {
				buf.FieldByName("buffer").SetCap(int(in))
				buff.AppendModeOff()
				capacitySet = true
				continue
			}

			panic(fmt.Errorf(
				"Buffer length (%d) and capacity (%d) are already set",
				len(buff.buffer), cap(buff.buffer),
			))
		default:
			panic(fmt.Errorf(
				"I'm sorry I don't know how to interpret the type \"%s\"!\n"+
					"If you feel I should be, please contact the developers at:\n\t"+
					strings.Join(Maintainers, "\n\t")+"\n", reflect.TypeOf(in),
			))
		}
	}

	if written {
		return buff.AppendModeOff().Rewind()
	}

	return buff
}

// NewBuff returns a *Buff with the provided length and optional capacity
func NewBuff(length, capacity ...int) (buff *Buff) {
	return MakeBuff(length, capacity...)
}

// ReadBuff returns a *Buff after reading the io.Reader.
func ReadBuff(r io.Reader) (mmap *Buff) {
	return MakeBuff(r)
}

// AppendModeOn enables appending when the buffer length as been reached.
func (b *Buff) AppendModeOn() *Buff {
	b.append = true
	return m
}

// AppendModeOff drops any writes when the buffer length has been reached.
func (b *Buff) AppendModeOff() *Buff {
	b.append = false
	return m
}

// Range returns a *RangedBuff which provides the same API with pointers
// to the range of bytes from the parent object.
func (b *Buff) Range(off, limit int64) *RangedMMap {
	if limit < 1 {
		limit = int64(len(m.buffer)) - off
	}

	return &RangedMMap{m, off, 0, limit}
}

// Seek sets the pointer for the *Buff which changes
// the default Read/Write location, this is not related to io.ReadSeeker.
func (b *Buff) Seek(n int64) *Buff {
	b.pointer = n
	return m
}

// Rewind sets the pointer to the beginning of the buffer.
func (b *Buff) Rewind() *Buff {
	return m.Seek(0)
}

// Read, you know io.Reader
func (b *Buff) Read(p []byte) (int, error) {
	n1, err := m.ReadAt(p, b.pointer)
	b.pointer += n1

	return int(n1), err
}

// Write, you know io.Writer
func (b *Buff) Write(p []byte) (int, error) {
	n1, err := m.WriteAt(p, b.pointer)
	b.pointer += n1

	return int(n1), err
}

// ReadAt, you know io.ReaderAt
// This method disregards the pointer entirely.
func (b *Buff) ReadAt(p []byte, off int64) (n int64, err error) {
	n = int64(copy(p, b.buffer[off:]))
	return
}

// WriteAt, you know io.WriterAt
// This method disregards the pointer entirely.
func (b *Buff) WriteAt(p []byte, off int64) (n int64, err error) {
	if diff := int64(len(b.buffer)) - off - int64(len(p)); b.append && diff < 0 {
		b.buffer = append(b.buffer, make([]byte, diff/-1)...)
	}

	n = int64(copy(b.buffer[off:off+int64(len(p))], p))
	return
}

// ReadFrom, you know io.ReaderFrom
func (b *Buff) ReadFrom(r io.Reader) (int64, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, err
	}

	n1, err := m.Write(buf)
	return int64(n1), err
}

// WriteTo, you know io.WriterTo
func (b *Buff) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(b.buffer[:])
	return int64(n), err
}
