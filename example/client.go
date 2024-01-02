// Package service @Author Bing
// @Date 2023/12/29 15:37:00
// @Desc
package main

import "github.com/learnselfs/GeeRPC/client"

func clt() {
	c := client.NewClient("127.0.0.1", 3000, "tcp4")
	c.Start()
}

func main() {
	clt()
}
