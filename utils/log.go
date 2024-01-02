// Package utils @Author Bing
// @Date 2023/12/26 10:57:00
// @Desc
package utils

import (
	"log"
	"os"
)

var (
	InfoLog  *log.Logger
	ErrLog   *log.Logger
	DebugLog *log.Logger
)

func init() {
	InfoLog = log.New(os.Stdout, "\x1b[;34m[INFO]\t\x1b[0m", log.Ldate|log.Ltime|log.Llongfile)
	ErrLog = log.New(os.Stdout, "\x1b[;31m[Error]\t\x1b[0m", log.Ldate|log.Ltime|log.Llongfile)
	DebugLog = log.New(os.Stdout, "\x1b[;33m[Debug]\t\x1b[0m", log.Ldate|log.Ltime|log.Llongfile)

}
