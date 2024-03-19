// Package GeeRPC @Author Bing
// @Date 2024/3/7 13:56:00
// @Desc
package GeeRPC

import (
	"bufio"
	"encoding/gob"
	"io"
)

type GobCodec struct {
	io     io.ReadWriteCloser
	buf    *bufio.Writer
	encode *gob.Encoder
	decode *gob.Decoder
}

func (g *GobCodec) Write(v any) (err error) {
	defer g.buf.Flush()
	return g.encode.Encode(v)
}

func (g *GobCodec) Read(v any) (err error) {
	err = g.decode.Decode(v)
	return
}

func NewGobCode(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		io:     conn,
		buf:    buf,
		encode: gob.NewEncoder(buf),
		decode: gob.NewDecoder(conn),
	}
}

var _ Codec = (*GobCodec)(nil)
