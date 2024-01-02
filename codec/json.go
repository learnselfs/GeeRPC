// Package codec @Author Bing
// @Date 2023/12/27 10:27:00
// @Desc
package codec

import (
	"bufio"
	"encoding/json"
	"github.com/learnselfs/GeeRPC/utils"
	"io"
	"net"
)

type JsonCodec struct {
	connect io.ReadWriteCloser
	buf     *bufio.Writer
	encode  *json.Encoder
	decode  *json.Decoder
}

func (j *JsonCodec) ReadMsg(msg *Message) error {
	return j.decode.Decode(msg)
}
func (j *JsonCodec) ReadHead(head *Head) error {
	return j.decode.Decode(head)
}
func (j *JsonCodec) ReadBody(body any) error {
	return j.decode.Decode(body)
}

func (j *JsonCodec) Write(msg *Message, head *Head, body any) error {
	defer func() {
		err := j.Close()
		if err != nil {
			utils.ErrLog.Println(err)
		}
	}()
	err := j.encode.Encode(msg)
	if err != nil {
		return err
	}
	err = j.encode.Encode(head)
	if err != nil {
		return err
	}
	err = j.encode.Encode(body)
	if err != nil {
		return err
	}
	return nil
}

func (j *JsonCodec) Close() error {
	return j.connect.Close()
}

func NewCodecJson(connect net.Conn) ICodec {
	buf := bufio.NewWriter(connect)
	return &JsonCodec{
		connect: connect,
		buf:     buf,
		encode:  json.NewEncoder(buf),
		decode:  json.NewDecoder(connect),
	}
}

var _ ICodec = (*JsonCodec)(nil)
