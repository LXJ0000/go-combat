package main

import (
	"strconv"
	"time"

	"log"
)

func main() {
	testContinue()
}

func testContinue() {
	in := make(chan *Content, 20)
	audit := make(chan *Content, 20)
	streamTextPreProcessStop := make(chan struct{})
	// 向in协程无脑放2000个数据
	go func() {
		for i := 0; i < 2000; i++ {
			in <- &Content{
				i: i,
			}
			log.Printf("put in content = %s", strconv.Itoa(i))
		}
	}()
	// 异步审核流程，在第三十条的时候触发审核失败
	go func() {
		for {
			select {
			case content, ok := <-audit:
				if !ok {
					log.Printf("audit get in not ok")
				}
				time.Sleep(30 * time.Millisecond) // 等待的时候就有问题了 因为 in 一直在写 主循环从 in 读 往 audit 写 audit 写满了 阻塞 导致主循环阻塞在2 没法进入streamTextPreProcessStop读的分支 第二个go阻塞在1 问题就来啦～
				if content.i == 30 {
					log.Printf("audit streamTextPreProcessStop before")
					streamTextPreProcessStop <- struct{}{} // 写数据 阻塞点1 
					log.Printf("audit streamTextPreProcessStop after")
				}
			}
		}
	}()

	for {
		select {
		case <-streamTextPreProcessStop:
			log.Printf("get streamTextPreProcessStop")
			waitTimes := 0
			for {
				if waitTimes > 50 {
					break
				}
				waitTimes++
				time.Sleep(100 * time.Millisecond)
			}
			continue
		case content, ok := <-in:
			if !ok {
				log.Printf("get in not ok")
			}
			log.Printf("get in content = %s", strconv.Itoa(content.i))
			log.Printf("audit in before content = %s", strconv.Itoa(content.i))
			audit <- content // 阻塞点2 
			log.Printf("audit in after content = %s", strconv.Itoa(content.i))
		}
	}
}

type Content struct {
	i int
}
