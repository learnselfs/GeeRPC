// Package GeeRPC @Author Bing
// @Date 2024/3/10 21:19:00
// @Desc struct remote procedure call
package GeeRPC

import (
	"reflect"
)

type Struct struct {
	name    string
	type_   reflect.Type
	value_  reflect.Value
	methods map[string]*Method
}

func (s *Struct) parseMethod() {
	for i := 0; i < s.type_.NumMethod(); i++ {
		method := s.type_.Method(i)
		name := method.Name
		methodNumIn := method.Type.NumIn()
		if methodNumIn == 3 {
			args := method.Type.In(1)
			reply := method.Type.In(2)
			s.methods[name] = NewMethod(args, reply, method)
		}
	}

}

func NewStruct(value any) *Struct {
	v := reflect.ValueOf(value)
	t := reflect.TypeOf(value)
	name := reflect.TypeOf(v.Elem().Interface()).Name()
	s := &Struct{name: name, value_: v, type_: t, methods: make(map[string]*Method)}
	s.parseMethod()
	return s
}

func (s *Struct) Call(method string, args, r reflect.Value) {
	m := s.methods[method]
	a := m.NewArgs()
	a.Set(args)
	m.method.Func.Call([]reflect.Value{s.value_, a, r})
	//reply = r.Elem().Interface()
	//wlog.Info(fmt.Sprintf("%s", reply))
	return
}

type Method struct {
	args   reflect.Type
	reply  reflect.Type
	method reflect.Method
}

func (m *Method) NewArgs() reflect.Value {
	var args reflect.Value
	if m.args.Kind() == reflect.Pointer {
		args = reflect.New(m.args.Elem())
	} else {
		args = reflect.New(m.args).Elem()
	}
	return args
}

func (m *Method) NewReply() reflect.Value {
	reply := reflect.New(m.reply.Elem())
	if m.reply.Kind() == reflect.Map {
		reply = reflect.MakeMap(m.reply)
	}
	if m.reply.Kind() == reflect.Slice {
		reply = reflect.MakeSlice(m.reply, 0, 0)
	}
	return reply
}

func NewMethod(args, reply reflect.Type, method reflect.Method) *Method {
	return &Method{args: args, reply: reply, method: method}
}
