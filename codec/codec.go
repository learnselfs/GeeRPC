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
	Method    string
	Sequence  uint64
}

type Session struct {
	Sequence uint64

	MsgChan chan *Message

	ClientClose  chan bool
	ServiceClose chan bool

	ICodec
	Wg *sync.WaitGroup
	Mu *sync.Mutex

	Error error

	*Message
}

type Message struct {
	CodecType string
	Head      *Head
	Args      reflect.Value
	Reply     reflect.Value
}

func (s *Session) CClose() {
	s.ClientClose <- true
	close(s.ServiceClose)
	close(s.ClientClose)
	close(s.MsgChan)
	s.Close()
}

func (s *Session) SClose() {
	s.ServiceClose <- true
	close(s.ServiceClose)
	close(s.ClientClose)
	close(s.MsgChan)
	s.Close()
}

func NewSession(conn net.Conn) *Session {
	//ic := CodecMap[codecType](conn)
	//head := &Head{Method: method, CodecType: codecType}
	//message := &Message{CodecType: codecType, Args: reflect.ValueOf(args), Reply: reflect.ValueOf(reply), Head: head}
	session := &Session{Wg: new(sync.WaitGroup), Mu: new(sync.Mutex), ClientClose: make(chan bool), ServiceClose: make(chan bool)}

	return session
}
