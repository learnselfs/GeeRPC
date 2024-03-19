// Package GeeRPC @Author Bing
// @Date 2024/3/13 16:16:00
// @Desc
package GeeRPC

import "strings"

func parseServiceMethod(serviceMethod string) (string, string) {
	serviceMethodSlice := strings.SplitN(serviceMethod, ".", 2)
	return serviceMethodSlice[0], serviceMethodSlice[1]

}
