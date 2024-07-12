package xchan

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)

		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			m := len(channels) / 2
			select {
			case <-or(channels[:m]...):
			case <-or(channels[m:]...):
			}
		}
	}()

	return orDone
}

func orWithSelect(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)
		var cases []reflect.SelectCase
		for _, c := range channels {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}

		_, recv, _ := reflect.Select(cases)
		if recv.IsValid() {
			orDone <- recv.Interface()
		}
	}()

	return orDone
}

func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()

	return c
}

func TestOrDone(t *testing.T) {
	start := time.Now()
	<-or(
		sig(1*time.Second),
		sig(2*time.Second),
		sig(4*time.Second),
		sig(8*time.Second),
		sig(16*time.Second),
		sig(32*time.Second),
		sig(64*time.Second),
		sig(128*time.Second),
		sig(256*time.Second),
		sig(512*time.Second),
		sig(1024*time.Second),
		sig(2048*time.Second),
		sig(4096*time.Second),
		sig(8192*time.Second),
		sig(16384*time.Second),
		sig(32*time.Second),
	)

	fmt.Printf("done after %v\n", time.Since(start))
}

func TestOrDoneWithSelect(t *testing.T) {
	start := time.Now()
	<-orWithSelect(
		sig(1*time.Second),
		sig(2*time.Second),
		sig(4*time.Second),
		sig(8*time.Second),
		sig(16*time.Second),
		sig(32*time.Second),
		sig(64*time.Second),
		sig(128*time.Second),
		sig(256*time.Second),
		sig(512*time.Second),
		sig(1024*time.Second),
		sig(2048*time.Second),
		sig(4096*time.Second),
		sig(8192*time.Second),
		sig(16384*time.Second),
		sig(32*time.Second),
	)

	fmt.Printf("done after %v\n", time.Since(start))
}
