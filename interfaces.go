// Package GeeRPC @Author Bing
// @Date 2024/3/7 13:57:00
// @Desc
package GeeRPC

import "io"

type Codec interface {
	Write(any) error
	Read(any) error
}
type (
	CodecType string
	CodecFunc func(closer io.ReadWriteCloser) Codec
)

const (
	Json CodecType = "application/json"
	Gob  CodecType = "application/gob"
)

var (
	CodecMap map[CodecType]CodecFunc
)

func init() {
	CodecMap = make(map[CodecType]CodecFunc)
	CodecMap[Json] = NewGobCode
	CodecMap[Gob] = NewGobCode
}

type Header struct {
	Id            uint64
	ServiceMethod string
	Err           error
}
