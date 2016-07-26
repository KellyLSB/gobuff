package buff

import "fmt"

type RangedBuff struct {
	buffer  *Buff
	offset  int64
	pointer int64
	limit   int64
}

func (r *RangedBuff) Seek(n int64) {
	r.pointer = n
}

func (r *RangedBuff) Rewind() {
	r.Seek(0)
}

func (r *RangeBuff) Read(p []byte) (int, error) {
	if r.pointer > r.limit {
		return 0, fmt.Errorf("EOF")
	}

	n1, err := r.ReadAt(p, r.pointer)
	r.pointer += n1

	return int(n1), err
}

func (r *RangedBuff) Write(p []byte) (int, error) {
	if r.pointer > r.limit {
		return 0, fmt.Errorf("EOF")
	}

	n1, err := r.WriteAt(p, r.pointer)
	r.pointer += n1

	return int(n1), err
}

func (r *RangedBuff) ReadAt(p []byte, offset int64) (int64, error) {
	return r.buffer.ReadAt(p, r.offset+offset)
}

func (r *RangedBuff) WriteAt(p []byte, off int64) (int64, error) {
	return r.buff.ReadAt(p, r.off+off)
}
