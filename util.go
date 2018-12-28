/*
Package vaultlib is a lightweight Go library for reading Vault KV secrets.
Interacts with Vault server using HTTP API only.
First create a new *Config object using NewConfig()
Then create you Vault client using NewClient(*Config)
*/
package vaultlib

import (
	"runtime"
	"strconv"
)

func errInfo() (info string) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.Function + ":" + strconv.Itoa(frame.Line)
}
