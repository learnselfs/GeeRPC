// Package GeeRPC @Author Bing
// @Date 2024/3/11 10:05:00
// @Desc
package GeeRPC

import (
	"bufio"
	"encoding/json"
	"io"
)

type JsonCodec struct {
	conn   io.ReadWriteCloser
	buf    *bufio.Writer
	encode *json.Encoder
	decode *json.Decoder
}

func (j *JsonCodec) Read(v any) (err error) {
	return j.decode.Decode(v)
}
func (j *JsonCodec) Write(v any) (err error) {
	defer j.buf.Flush()
	return j.encode.Encode(v)
}

func NewJsonCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &JsonCodec{
		conn:   conn,
		buf:    buf,
		encode: json.NewEncoder(buf),
		decode: json.NewDecoder(conn),
	}
}

var _ Codec = (*JsonCodec)(nil)
