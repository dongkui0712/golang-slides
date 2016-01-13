package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var (
	_ = time.Second
	_ = sync.Mutex{}
	_ = atomic.AddInt64
	_ = runtime.GOMAXPROCS
)

func main() {

	var a int64 = 100000

	for n := 0; n < 10; n++ {
		go func() {
			for {
				if a > 0 {
					a--
				} else {
					break
				}
			}
		}()
	}

	time.Sleep(time.Second)
	fmt.Println(a)
}
