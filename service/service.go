// Package service @Author Bing
// @Date 2023/12/25 11:12:00
// @Desc
package service

import (
	"fmt"
	"github.com/learnselfs/GeeRPC/codec"
	"github.com/learnselfs/GeeRPC/utils"
	"net"
)

type Service struct {
	host         string
	port         string
	ipVersion    string
	sessionId    uint64 // sequence number
	sessionTasks map[uint64]*codec.Session
}

func NewService(host string, port string, ipVersion string) *Service {
	return &Service{host: host, port: port, ipVersion: ipVersion, sessionId: uint64(-1), sessionTasks: make(map[uint64]*codec.Session)}
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
	session := codec.NewSession(connect)
	s.sessionId++
	s.sessionTasks[s.sessionId] = session
	session.Sequence = s.sessionId

	ic := codec.CodecMap[codec.JsonCodecType](connect)
	err := ic.ReadMsg(session.Message)
	if err != nil {
		utils.ErrLog.Printf("Connecting error: %s", err)
	}
	if len(session.Message.CodecType) <= 0 {
		utils.ErrLog.Println("client codec type  error")
	}
	session.ICodec = codec.CodecMap[session.Message.CodecType](connect)
	for {
		go s.read(session)
		go s.write(session)
		select {
		case <-session.ClientClose:
			utils.ErrLog.Println("Client close connect")
			break
		case <-session.ServiceClose:
			utils.ErrLog.Println("service close connect")
			break

		}
	}
}

func (s *Service) read(session *codec.Session) {
	err := session.ICodec.ReadHead(session.Message.Head)
	if err != nil {
		utils.ErrLog.Println("decode client session.message.Head error")
		session.CClose()
	}

	err = session.ICodec.ReadBody(session.Message.Args)
	if err != nil {
		utils.ErrLog.Println("decode client session.message.Args error")
		session.CClose()
	}

	session.MsgChan <- session.Message
}

func (s *Service) write(session *codec.Session) {

	select {
	case msg := <-session.MsgChan:
		err := session.ICodec.Write(msg, session.Message.Head, session.Message.Reply)
		if err != nil {
			utils.ErrLog.Println("encode client session.message.Reply error")
			session.SClose()
		}
	}
}
