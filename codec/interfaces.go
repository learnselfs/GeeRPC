// Package codec @Author Bing
// @Date 2023/12/27 10:02:00
// @Desc
package codec

type ICodec interface {
	ReadMsg(msg *Message) error
	ReadHead(head *Head) error
	ReadBody(body any) error
	Write(msg *Message, head *Head, body any) (err error)
	Close() error
}
