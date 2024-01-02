// Package codec @Author Bing
// @Date 2023/12/27 10:11:00
// @Desc
package codec

import (
	"net"
	"reflect"
	"sync"
)

var (
	CodecMap map[string]NewCodecFunc
)

const (
	GobCodecType  = "application/gob"
	JsonCodecType = "application/json"
)

type NewCodecFunc func(connect net.Conn) ICodec

func init() {
	CodecMap = make(map[string]NewCodecFunc)
	CodecMap[GobCodecType] = NewCodecGob
	CodecMap[JsonCodecType] = NewCodecJson
}

type Head struct {
	CodecType string
	Length    int
}

type Package struct {
	Method   string
	Sequence uint
}
type Message struct {
	Head  *Head
	Args  reflect.Value
	Reply reflect.Value

	ClientClose  chan bool
	ServiceClose chan bool

	ICodec
	Wg *sync.WaitGroup
	Mu *sync.Mutex
}
