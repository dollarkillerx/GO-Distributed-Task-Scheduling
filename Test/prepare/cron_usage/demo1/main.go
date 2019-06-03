/**
* Created by GoLand
* User: dollarkiller
* Date: 19-6-3
* Time: 下午3:41
* */
package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

func main() {
	// 解析cron表达式
	expression, e := cronexpr.Parse("*/5 * * * * * *")// 秒 分 时 日 月 周 年
	if e != nil {
		panic(e.Error())
	}
	now := time.Now()
	// 下一次调度时间
	next := expression.Next(now)

	// 等待定时器超时
	time.AfterFunc(next.Sub(now), func() {
		fmt.Println(next.Sub(now))
		fmt.Println("被调度了",next)
	})
	time.Sleep(10 * time.Second)
}