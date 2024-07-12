package main

import (
	"fmt"
	"sync"
)

type Counter struct {
	sync.Mutex
	Count int
}

func (c *Counter) add() int {
	c.Lock()
	defer c.Unlock()
	c.Count++

	return c.getc()
}

func (c *Counter) getc() int {
	c.Lock()
	defer c.Unlock()

	return c.Count
}

func foo(c *Counter) int {
	c.Lock()
	defer c.Unlock()
	fmt.Println("in foo")
	return 1
}

func main() {
	var c Counter
	fmt.Println(c.add())
}
