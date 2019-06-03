/**
* Created by GoLand
* User: dollarkiller
* Date: 19-6-3
* Time: 下午3:06
* */
package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

type rest struct {
	output []byte
	err error
}

var (
	outch = make(chan *rest)
)



func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	// 在一个协程里,执行一个cmd 让他执行2s
	go func(ctx context.Context) {
		commandContext := exec.CommandContext(ctx, "/bin/bash", "-c", "sleep 2;echo hello;")
		bytes, e := commandContext.CombinedOutput()

		outch <- &rest{output:bytes,err:e}
	}(ctx)
	// 1s的是否,我们杀死cmd
	time.Sleep(time.Second)
	cancel()
	i := <-outch
	if i.err != nil {
		fmt.Println(i.err.Error())
	}else {
		fmt.Println(string(i.output))
	}
}
