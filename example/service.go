// Package main @Author Bing
// @Date 2023/12/29 10:43:00
// @Desc
package main

import "github.com/learnselfs/GeeRPC/service"

func server() {
	s := service.NewService("127.0.0.1", "3000", "tcp4")
	err := s.Start()
	if err != nil {
		return
	}
}

func main() {
	server()
}
