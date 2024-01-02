// Package utils @Author Bing
// @Date 2024/1/9 10:35:00
// @Desc
package utils

import "reflect"

type Server struct {
	Name    string
	Type    reflect.Type
	Value   reflect.Value
	Methods map[string]*Method
}

type Method struct {
	Name        string
	Method      reflect.Method
	ValueMethod reflect.Value
	Args        reflect.Type
	Reply       reflect.Type
}

func RegisterServer(s any) *Server {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)
	name := t.Name()

	sv := &Server{
		Name:    name,
		Value:   v,
		Type:    t,
		Methods: make(map[string]*Method),
	}
	sv.registerMethod()
	return sv
}

func (s *Server) registerMethod() {
	for i := 0; i < s.Type.NumMethod(); i++ {
		valueMethod := s.Value.Method(i)
		method := s.Type.Method(i)
		name := method.Name
		args := method.Type.In(1)
		reply := method.Type.In(2)
		s.Methods[name] = &Method{Name: name, Method: method, ValueMethod: valueMethod, Args: args, Reply: reply}
	}
}

func (s *Server) Call(methodName string, args, reply any) {
	s.Value.MethodByName(methodName).Call([]reflect.Value{reflect.ValueOf(args), reflect.ValueOf(reply)})

}
