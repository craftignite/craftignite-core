package minecraft

import (
	"bytes"
	"encoding/binary"
)

type Buffer struct {
	data   []byte
	offset uint
}

func (buf *Buffer) ReadLong() uint64 {
	var value uint64

	readBuf := bytes.NewBuffer(buf.data[buf.offset : buf.offset+8])
	err := binary.Read(readBuf, binary.LittleEndian, &value)

	if err != nil {
		return 0
	}
	return value
}

func (buf *Buffer) ReadByte() byte {
	val := buf.data[buf.offset]
	buf.offset++
	return val
}

func (buf *Buffer) ReadVarInt() int {
	numRead, result := 0, 0

	for {
		read := buf.ReadByte()
		val := int(read & 0b01111111)
		result |= val << (7 * numRead)
		numRead++

		if read&0b10000000 == 0 {
			break
		}
	}

	return result
}

func (buf *Buffer) WriteLong(val uint64) {
	binary.LittleEndian.PutUint64(buf.data[buf.offset:buf.offset+8], val)
	buf.offset += 8
}

func (buf *Buffer) WriteByte(val byte) {
	buf.data[buf.offset] = val
	buf.offset++
}

func (buf *Buffer) WriteVarInt(val int) {
	for {
		b := byte(val & 0b01111111)
		val >>= 7

		if val != 0 {
			b |= 0b10000000
		}
		buf.WriteByte(b)

		if val == 0 {
			break
		}
	}
}

func (buf *Buffer) WriteBytes(data []byte) {
	for _, val := range data {
		buf.WriteByte(val)
	}
}

func (buf *Buffer) WriteString(data string) {
	binaryData := []byte(data)
	buf.WriteVarInt(len(binaryData))
	for _, b := range binaryData {
		buf.WriteByte(b)
	}
}

func (buf *Buffer) Skip(num int) {
	buf.offset += uint(num)
}