// Package client @Author Bing
// @Date 2023/12/25 11:19:00
// @Desc
package client

import (
	"encoding/json"
	"fmt"
	"github.com/learnselfs/GeeRPC/codec"
	"github.com/learnselfs/GeeRPC/utils"
	"net"
	"reflect"
	"sync"
)

type Client struct {
	remoteHost string
	remotePort int
	ipVersion  string
}

func NewClient(remoteHost string, remotePort int, ipVersion string) *Client {
	return &Client{remoteHost, remotePort, ipVersion}
}

func (c *Client) Start() {
	connect, err := c.connect()
	if err != nil {
		utils.ErrLog.Println(err)
		return
	}
	c.accept(connect)
}

func (c *Client) connect() (net.Conn, error) {
	return net.Dial(c.ipVersion, fmt.Sprintf("%s:%d", c.remoteHost, c.remotePort))
}

func (c *Client) accept(connect net.Conn) {
	pkg := codec.Package{Method: codec.GobCodecType, Sequence: uint(1)}
	err := json.NewEncoder(connect).Encode(pkg)
	if err != nil {
		utils.ErrLog.Println(err)
		return
	}
	ic := codec.CodecMap[codec.GobCodecType](connect)
	msg := &codec.Message{Wg: new(sync.WaitGroup), Mu: new(sync.Mutex), ICodec: ic, ClientClose: make(chan bool), ServiceClose: make(chan bool)}
	for i := 0; i < 10; i++ {
		msg.Head = &codec.Head{CodecType: codec.GobCodecType, Length: i}

		msg.Args = reflect.ValueOf("hello world")
		err := ic.Write(msg.Head, msg.Args.Interface())
		if err != nil {
			utils.ErrLog.Println(err)
			return
		}
		for j := 0; j < 2; j++ {
			request(msg)
		}
	}
	msg.Close()
}
func request(msg *codec.Message) {

	err := msg.ICodec.ReadHead(msg.Head)
	if err != nil {
		utils.ErrLog.Println(err)
	}
	arg := ""
	err = msg.ICodec.ReadBody(&arg)
	if err != nil {
		utils.ErrLog.Println(err)
	}
	utils.InfoLog.Printf("%#v, %s", msg.Head, arg)

}
