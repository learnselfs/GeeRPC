// Package client @Author Bing
// @Date 2023/12/25 11:19:00
// @Desc
package client

import (
	"fmt"
	"github.com/learnselfs/GeeRPC/codec"
	"github.com/learnselfs/GeeRPC/utils"
	"net"
	"sync"
)

type Client struct {
	remoteHost  string
	remotePort  int
	ipVersion   string
	sequence    uint64
	sessions    map[uint64]*codec.Session
	sessionTask chan *codec.Session
	conn        net.Conn
	mu          *sync.Mutex
}

func NewClient(remoteHost string, remotePort int, ipVersion string) *Client {

	client := &Client{
		remoteHost: remoteHost,
		remotePort: remotePort,
		ipVersion:  ipVersion,
		sequence:   uint64(0),
		sessions:   make(map[uint64]*codec.Session),
		mu:         new(sync.Mutex)}
	client.receive()
	return client
}

func (c *Client) registerSession(session *codec.Session) (uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var err error
	session.Sequence = c.sequence
	c.sessions[c.sequence] = session
	c.sequence++
	return session.Sequence, err
}

func (c *Client) removeSession(sequence uint64) *codec.Session {
	c.mu.Lock()
	defer c.mu.Unlock()
	session := c.sessions[sequence]
	delete(c.sessions, sequence)
	return session
}

func (c *Client) terminateSession(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, s := range c.sessions {
		s.Error = err
		s.ICodec.Close()
	}
}

func (c *Client) receive() {
	connect, err := c.connect()
	if err != nil {
		utils.ErrLog.Println(err)
		return
	}
	c.conn = connect
	c.accept(connect)
}

func (c *Client) connect() (net.Conn, error) {
	return net.Dial(c.ipVersion, fmt.Sprintf("%s:%d", c.remoteHost, c.remotePort))
}

func (c *Client) accept(connect net.Conn) {
	for {
		select {
		case session := <-c.sessionTask:
			err := session.ICodec.ReadMsg(session.Message)
			if err != nil {
				session.CClose()
				return
			}
			session.ICodec = codec.CodecMap[session.Message.CodecType](connect)
			c.send(session)
		}
	}
}
func (c *Client) read(session *codec.Session, ic codec.ICodec) {

	err := ic.ReadHead(session.Head)
	err = ic.ReadBody(session.Args)
	if err != nil {
		utils.ErrLog.Println(err)
		session.CClose()
		return
	}
}

//func (c *Client) write() {
//	c.mu.Lock()
//	defer c.mu.Unlock()
//
//	session := c.removeSession()
//	c.sessionTask <- session
//
//	session.Mu.Lock()
//	defer session.Mu.Unlock()
//
//	err := session.ICodec.Write(session.Message, session.Head, session.Args.Interface())
//	if err != nil {
//		utils.ErrLog.Println(err)
//		return
//	}
//}

func (c *Client) sendDefaultCodec(session *codec.Session) {
	ic := codec.CodecMap[codec.GobCodecType](c.conn)
	c.sendCodec(session, ic)
}

func (c *Client) send(session *codec.Session) {
	c.sendCodec(session, session.ICodec)
}

func (c *Client) sendCodec(session *codec.Session, ic codec.ICodec) {
	session.Mu.Lock()
	defer session.Mu.Unlock()

	err := ic.ReadHead(session.Head)
	if err != nil {
		utils.ErrLog.Println(err)
	}

	err = ic.ReadBody(session.Args.Elem().Addr())
	if err != nil {
		utils.ErrLog.Println(err)
	}
	c.sessionTask <- session
	utils.InfoLog.Printf("%#v, %s", session.Head, session.Args)
}

func (c *Client) Go(method string, args, reply any) {
	session := c.sendSession(method, args, reply, codec.GobCodecType)
	_, _ = c.registerSession(session)
	c.send(session)
}

func (c *Client) Call(method string, args, repay any) {

}

//func defaultSession(codecType string) *codec.Session {
//	head := &codec.Head{CodecType: codecType}
//	message := &codec.Message{Head: head}
//	session := &codec.Session{Message: message, Wg: new(sync.WaitGroup), Mu: new(sync.Mutex), ClientClose: make(chan bool), ServiceClose: make(chan bool)}
//	return session
//}
