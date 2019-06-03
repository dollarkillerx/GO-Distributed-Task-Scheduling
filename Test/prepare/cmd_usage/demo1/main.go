/**
* Created by GoLand
* User: dollarkiller
* Date: 19-6-3
* Time: 下午2:58
* */
package main

import (
	"fmt"
	"os/exec"
)

func main() {
	// 生成cmd
	command := exec.Command("/bin/bash", "-c", "sleep 5;ls -l")

	// 执行命令,捕获了子进程的输出(pipe)
	bytes, e := command.CombinedOutput()
	if e != nil {
		panic(e.Error())
	}

	fmt.Println((string(bytes)))
}
