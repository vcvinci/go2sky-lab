package common

import (
	"fmt"
	"os"
	"runtime"
)

const (
	SkyAddr = "127.0.0.1:11800"
)

func PanicError(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Fprintf(os.Stderr, "[%s:%d] %s", file, line, err)
		panic(err)
	}
}

type Params struct {
	Name string `json:"name"`
}
