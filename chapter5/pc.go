package main

import (
	"fmt"
	"math/rand"
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

	msgPool := make(chan *Msg, 10)
	for i := 0; i < 5; i++ {
		go func() {
			p := &Producer{}
			for n := 0; n < 5; n++ {
				msgPool <- p.GenerateMsg()
			}
		}()
	}

	for i := 0; i < 3; i++ {
		consumer := Consumer{msgPool}
		go consumer.Run()
	}

	time.Sleep(20 * time.Second)
}
