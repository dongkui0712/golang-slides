package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// msg_struct_start OMIT
type Msg struct {
	M      int
	Cancel chan struct{}
}

func (m *Msg) Do() {

	fmt.Println("doing:", m.M)
}

// msg_struct_end OMIT

// producer_struct_start OMIT
type Producer struct {
}

func (p *Producer) GenerateMsg() *Msg {

	c := make(chan struct{}, 1)
	go func() {
		time.Sleep(time.Second)
		close(c)
	}()
	return &Msg{rand.Int() % 100, c}
}

// producer_struct_end OMIT

// consumer_struct_start OMIT
type Consumer struct {
	MsgPool chan *Msg
}

func (c *Consumer) Run() {

	for {
		msg, ok := <-c.MsgPool
		if !ok {
			break
		}

		select {
		case <-time.After(time.Second / 2):
			msg.Do()
		case <-msg.Cancel:
			fmt.Println("msg is canceled")
		}
	}
}

// consumer_struct_end OMIT

func main() {
	var wg sync.WaitGroup
	msgPool := make(chan *Msg, 10)

	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func() {
			p := &Producer{}
			for n := 0; n < 5; n++ {
				msgPool <- p.GenerateMsg()
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(msgPool)
	}()
	for i := 0; i < 2; i++ {
		consumer := Consumer{msgPool}
		consumer.Run()
	}
	consumer := Consumer{msgPool}
	consumer.Run()
}
