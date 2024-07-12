package main

import "fmt"

type mapKey struct {
	key int
}

func main() {
	m := make(map[*mapKey]string)
	key := &mapKey{key: 10}
	m[key] = "hello world"
	fmt.Printf("map[key]=%v\n", m[key])

	key.key = 100
	fmt.Printf("map[key]=%v\n", m[key])
}
