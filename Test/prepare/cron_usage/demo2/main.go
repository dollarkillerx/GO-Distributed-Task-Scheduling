/**
* Created by GoLand
* User: dollarkiller
* Date: 19-6-3
* Time: 下午4:01
* */
package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

type CronJob struct {
	cron *cronexpr.Expression
	tim time.Time
}

func main() {
	// 需要一个调度协程,它定时检查所有的Cron任务,谁过期就执行谁
	cron := cronexpr.MustParse("*/5 * * * * * *")
	nextTime := cron.Next(time.Now())
	jobs := make(map[string]*CronJob)
	jobs["job1"] = &CronJob{cron:cron,tim:nextTime}

	cron = cronexpr.MustParse("*/1 * * * * * *")
	nextTime = cron.Next(time.Now())
	jobs["job2"] = &CronJob{cron:cron,tim:nextTime}

	go func() {
		for {
			select {
			case <- time.NewTimer(300 * time.Millisecond).C:
				now := time.Now()
				for k,v := range jobs {
					go func(k string,v *CronJob) {
						if v.tim.Before(now) || v.tim.Equal(now) {
							fmt.Println("执行任务",k)
						}
						v.tim = v.cron.Next(now)
					}(k,v)
				}
			}
		}
	}()

	time.Sleep(10 * time.Second)
}
