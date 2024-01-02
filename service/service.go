// Package service @Author Bing
// @Date 2023/12/25 11:12:00
// @Desc
package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/learnselfs/GeeRPC/codec"
	"github.com/learnselfs/GeeRPC/utils"
	"io"
	"net"
	"reflect"
	"sync"
)

type Service struct {
	host      string
	port      string
	ipVersion string
}

func NewService(host string, port string, ipVersion string) *Service {
	return &Service{host, port, ipVersion}
}
func (s *Service) Start() error {
	return s.connect()
}

func (s *Service) Stop() {

}

func (s *Service) connect() error {
	listen, err := net.Listen(s.ipVersion, fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		utils.ErrLog.Println(err)
		return err
	}
	utils.InfoLog.Printf("Connecting to %s:%s", s.host, s.port)
	for {
		connect, err := listen.Accept()
		if err != nil {
			return err
		}
		go s.accept(connect)
	}
}

func (s *Service) accept(connect net.Conn) {
	defer func() { _ = connect.Close() }()
	var pkg codec.Package
	err := json.NewDecoder(connect).Decode(&pkg)
	if err != nil {
		utils.ErrLog.Println(err)
		return
	}

	if len(pkg.Method) > 0 {
		utils.InfoLog.Printf("%#v", pkg)
	}
	//fmt.Printf("%#v\n, %#v\n", pkg, codec.CodecMap)
	ic := codec.CodecMap[pkg.Method](connect)
	s.serviceCodec(ic)
}

func (s *Service) serviceCodec(ic codec.ICodec) {
	msg := &codec.Message{Wg: new(sync.WaitGroup), Mu: new(sync.Mutex), ICodec: ic, ClientClose: make(chan bool), ServiceClose: make(chan bool)}
	for {
		err := s.request(msg)
		if err != nil {
			msg.ICodec.Close()
			return
		}

		err = s.response(msg)
		if err != nil {
			msg.ICodec.Close()
			utils.ErrLog.Println(err)
			return
		}

		msg.Wg.Add(1)
		go s.handle(msg)
	}
	msg.Wg.Wait()
	msg.ICodec.Close()

}

func (s *Service) request(msg *codec.Message) (err error) {
	msg.Head, err = s.readHead(msg.ICodec)
	if err != nil {
		return err
	}
	msg.Args, err = s.readBody(msg.ICodec)
	if err != nil {
		return err
	}
	return err

}

func (s *Service) response(msg *codec.Message) error {
	msg.Mu.Lock()
	defer msg.Mu.Unlock()
	err := msg.ICodec.Write(msg.Head, msg.Args.Interface())
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) handle(msg *codec.Message) {
	defer msg.Wg.Done()
	msg.Args = reflect.ValueOf(fmt.Sprintf("info: %s", msg.Head))
	err := s.response(msg)
	if err != nil {
		return
	}
}

func (s *Service) readHead(ic codec.ICodec) (*codec.Head, error) {
	h := codec.Head{}
	err := ic.ReadHead(&h)
	if err != nil {
		if err == io.EOF || errors.Is(err, errors.ErrUnsupported) {
			utils.InfoLog.Printf("client closed connection")
			return nil, err
		}
		utils.ErrLog.Println(err)
	}
	return &h, err
}

func (s *Service) readBody(ic codec.ICodec) (reflect.Value, error) {
	args := reflect.New(reflect.TypeOf(""))
	err := ic.ReadBody(args.Interface())
	if err != nil {
		utils.ErrLog.Println(err)
	}
	return args, err
}
