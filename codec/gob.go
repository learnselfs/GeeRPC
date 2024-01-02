// Package codec @Author Bing
// @Date 2023/12/27 10:10:00
// @Desc
package codec

import (
	"bufio"
	"encoding/gob"
	"github.com/learnselfs/GeeRPC/utils"
	"io"
	"net"
)

type GobCodec struct {
	connect io.ReadWriteCloser
	buf     *bufio.Writer
	encode  *gob.Encoder
	decode  *gob.Decoder
}

func (g *GobCodec) ReadHead(head *Head) error {
	return g.decode.Decode(head)
}

func (g *GobCodec) ReadBody(body any) error {
	return g.decode.Decode(body)
}

func (g *GobCodec) Write(head *Head, body any) error {
	defer func() {
		err := g.buf.Flush()
		if err != nil {
			utils.ErrLog.Println(err)
		}
	}()
	err := g.encode.Encode(head)
	if err != nil {
		return err
	}
	err = g.encode.Encode(body)
	if err != nil {
		return err
	}
	return nil
}

func (g *GobCodec) Close() error {
	return g.connect.Close()
}

func NewCodecGob(connect net.Conn) ICodec {
	buf := bufio.NewWriter(connect)
	return &GobCodec{
		connect: connect,
		buf:     buf,
		encode:  gob.NewEncoder(buf),
		decode:  gob.NewDecoder(connect),
	}
}

var _ ICodec = (*GobCodec)(nil)
