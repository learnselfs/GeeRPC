// Package example @Author Bing
// @Date 2023/12/26 11:20:00
// @Desc
package main

import "github.com/learnselfs/GeeRPC/utils"

func log() {
	//utils.InfoLog.Println("info.....")
	utils.InfoLog.Printf("%s\n", "--------------------------------")
	utils.ErrLog.Printf("%s\n", "--------------------------------")
	utils.DebugLog.Printf("%s\n", "--------------------------------")
}

func main() {
	log()
}
