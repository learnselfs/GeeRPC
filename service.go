// Package GeeRPC @Author Bing
// @Date 2024/3/7 13:54:00
// @Desc
package GeeRPC

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/learnselfs/wlog"
	"io"
	"net"
	"reflect"
	"strconv"
	"sync"
)

type Rpc struct {
	host    string
	port    int
	addr    string
	structs map[string]*Struct
}

func (r *Rpc) Serve() {
	r.listen()
}

func (r *Rpc) listen() {
	listen, _ := net.ResolveTCPAddr("tcp", r.addr)
	tcp, _ := net.ListenTCP("tcp", listen)

	for {
		con, _ := tcp.Accept()
		go r.Conn(con)
	}
}

func (r *Rpc) Conn(con net.Conn) {
	defer con.Close()
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	var codecType CodecType
	err := json.NewDecoder(con).Decode(&codecType)
	codec := CodecMap[codecType](con)
	for {
		var head Header
		//var m Msg
		err = r.Read(codec, &head)
		if err != nil {
			//wlog.Errorf("[server - head]: %#v", err)
			break
		}
		service, method := parseServiceMethod(head.ServiceMethod)
		s := r.structs[service]
		m := s.methods[method]

		args := m.NewArgs()
		argsI := args.Interface()
		if args.Type().Kind() != reflect.Ptr {
			argsI = args.Addr().Interface()
		}
		err = r.Read(codec, argsI)

		if !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
			if len(head.ServiceMethod) > 0 {
				reply := m.NewReply()
				s.Call(method, args, reply)
				wlog.Info(fmt.Sprintf("%#v, %#v, %#v", head, args, reply))
				r.Write(codec, head, reply.Interface())
			}
		}
	}
}

func (r *Rpc) Register(structDetail any) {
	s := NewStruct(structDetail)
	r.structs[s.name] = s
}

func (r *Rpc) Registers(structs ...any) {
	for _, s := range structs {
		r.Register(s)
	}
}

func (r *Rpc) Read(codec Codec, head any) (err error) {
	if err = codec.Read(head); err != nil {
		wlog.Errorf("[server - read]: %#v", err)
	}
	return
}

func (r *Rpc) Write(codec Codec, head, reply any) {
	err := codec.Write(head)
	if err != nil {
		wlog.Errorf("[server - write - head]: %#v", err)
		return
	}
	err = codec.Write(reply)
	if err != nil {
		wlog.Errorf("[server - write - head]: %#v", err)
		return
	}

}

func New(host string, port int) *Rpc {
	return &Rpc{host: host, port: port, addr: net.JoinHostPort(host, strconv.Itoa(port)), structs: make(map[string]*Struct)}
}
