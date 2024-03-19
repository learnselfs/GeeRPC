// Package GeeRPC @Author Bing
// @Date 2024/3/7 13:59:00
// @Desc
package GeeRPC

import (
	"context"
	"fmt"
	"github.com/learnselfs/wlog"
	"sync"
	"testing"
	"time"
)

const (
	host = "localhost"
	port = 80
)

var wg sync.WaitGroup

func server(ctx context.Context) {

	s := New(host, port)
	s.Registers(&Header{}, &Msg{})
	s.Serve()
	//select {
	//case <-ctx.Done():
	//	wg.Done()
	//	return
	//}
}

func client(cancel context.CancelFunc) {
	c := NewClient(host, port)
	c.Dial()
	var result int
	err := c.Call("Msg.Add", Args{Num1: 1, Num2: 2}, &result)
	if err != nil {
		wlog.Error(fmt.Sprintf("[client-end]: %#v", err))
	}
	wlog.Infof("%d", result)
	//cancel()
	//wg.Done()
}

func TestGeeRPC(t *testing.T) {
	wg.Add(2)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	go server(ctx)
	go client(cancel)
	wg.Wait()
}
func TestRpcServer(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	server(ctx)
}
func TestRpcClient(t *testing.T) {
	_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	client(cancel)
}
func TestCall(t *testing.T) {
	//var msg Msg
	//s := reflect2.newStruct()
	//s.registerRpc(&msg)
	//
	//m, _ := s.methods["Add"]
	//arg := m.newArgs()
	//arg.Set(reflect.ValueOf(Args{1, 2}))
	//reply := m.newReply()
	//s := NewStruct(&msg)
	//args := Args{1, 2}
	//
	//reply := s.Call("Add", args)
	//wlog.Infof("%#v", reply)
}
