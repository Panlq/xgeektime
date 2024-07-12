package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func TestGoroutineLeak() int64 {
	done := make(chan int64)
	t := time.NewTicker(5 * time.Second)
	go func() {
		loopCount := 0
		for {
			c := rand.Int63n(10)
			if c == 5 {
				fmt.Printf("循环%d次命中结果[%d]\n", loopCount, c)
				done <- c
				// 如果time ticker 在5秒后退出，channel赋值将会被阻塞, 解决方法，使用带缓冲带channel
				fmt.Println("退出go for")
				break
			}

			loopCount += 1
			fmt.Printf("goroutine loop at: %d, val: %d\n", loopCount, c)
			time.Sleep(1 * time.Second)
		}
	}()

	select {
	case <-t.C:
		fmt.Println("时间到，退出")
		return 0
	case v := <-done:
		fmt.Println("执行结束")
		return v
	}
}

func main() {
	val := TestGoroutineLeak()
	fmt.Println("执行结束：", val)
	// 如果main启动的是一个常驻的服务，上面那个goroutine就会泄漏
	time.Sleep(20 * time.Second)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world!")
	})

	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	select {}
}
