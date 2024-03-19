// Package GeeRPC @Author Bing
// @Date 2024/3/7 14:52:00
// @Desc
package GeeRPC

import (
	"encoding/json"
	"errors"
	"github.com/learnselfs/wlog"
	"net"
	"strconv"
	"sync"
)

type Client struct {
	host  string
	port  int
	addr  string
	codec Codec
	mu    sync.Mutex
	tasks map[uint64]*Task
}

func NewClient(host string, port int) *Client {
	return &Client{host: host, port: port, addr: net.JoinHostPort(host, strconv.Itoa(port)), tasks: make(map[uint64]*Task)}
}

type Task struct {
	id            uint64
	ServiceMethod string
	Args          any
	Reply         any
	channel       chan *Task
	err           error
}

func (t *Task) done() {
	t.channel <- t
}

func newTask(serviceMethod string, args, reply any, channel chan *Task) *Task {
	return &Task{
		ServiceMethod: serviceMethod,
		Args:          args,
		Reply:         reply,
		channel:       channel,
	}
}

type Args struct {
	Num1 int
	Num2 int
}
type Msg struct {
	Id   int
	Data string
}

func (m *Msg) Add(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func (c *Client) Dial(codecType ...CodecType) {
	conn, _ := net.Dial("tcp", c.addr)
	var defaultCodecType CodecType
	if len(codecType) == 0 {
		defaultCodecType = Gob
	}
	if len(codecType) >= 1 {
		defaultCodecType = codecType[0]
	}
	json.NewEncoder(conn).Encode(defaultCodecType)
	c.codec = CodecMap[defaultCodecType](conn)
	go c.Reader()
}

func (c *Client) Reader() {
	var err error
	for err == nil {
		var head Header
		err = c.codec.Read(&head)
		if err != nil {
			break
		}
		task := c.removeTask(head.Id)
		switch {
		case task == nil:
			err = c.codec.Read(nil)
		default:
			err = c.codec.Read(task.Reply)
			if err != nil {
				wlog.Errorf("[client - read - raply]: %#v", err)
				break
			}
			task.done()
		}
	}
	c.terminateTask(err)
}

func (c *Client) Writer(task *Task) {
	var err error
	var id uint64
	id, err = c.registerTask(task)
	var head Header
	head.Id = id
	head.ServiceMethod = task.ServiceMethod
	err = c.codec.Write(head)
	if err != nil {
		wlog.Errorf("[client - head] write error: %#v", err)
		task.done()
	}
	err = c.codec.Write(task.Args)
	if err != nil {
		wlog.Errorf("[client - args] write error: %#v", err)
		task.done()
	}
}

func (c *Client) registerTask(task *Task) (uint64, error) {
	//c.mu.Lock()
	//defer c.mu.Unlock()
	id := uint64(len(c.tasks))
	task.id = id
	c.tasks[id] = task
	return id, nil
}

func (c *Client) removeTask(id uint64) *Task {
	if task, ok := c.tasks[id]; ok {
		delete(c.tasks, id)
		return task
	}
	return nil
}

func (c *Client) terminateTask(err error) {
	for _, task := range c.tasks {
		task.err = err
		task.done()
	}
}

func (c *Client) close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, task := range c.tasks {
		task.err = errors.New("connection closed")
		task.done()
	}
}

func (c *Client) Call(serviceMethod string, args, reply any) error {
	channel := make(chan *Task, 1)
	task := c.Go(serviceMethod, args, reply, channel)
	for {

		select {
		case t := <-task.channel:
			return t.err
		}
	}
}

func (c *Client) Go(serviceMethod string, args, reply any, channel chan *Task) *Task {
	task := newTask(serviceMethod, args, reply, channel)
	c.Writer(task)
	return task

}
