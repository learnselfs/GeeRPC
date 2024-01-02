// Package main @Author Bing
// @Date 2024/1/10 15:07:00
// @Desc
package main

import (
	"fmt"
	"github.com/learnselfs/GeeRPC/utils"
)

type info struct {
	id int
	a  int
	b  int
}

func (i *info) Add(a, b int) (int, error) {
	c := i.a + i.b
	return c, nil
}

func main() {
	var c int
	var err error
	i := utils.RegisterServer(&info{id: 1, a: 5, b: 10})
	i.Call("Add")
	fmt.Println(c, err)

}
