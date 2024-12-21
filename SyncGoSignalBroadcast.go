package main

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"
)

/*
func (c *Cond) Wait(unique string, getTime func() string, printLn func(a ...interface{}), Sleep func(int)) {
	c.checker.check()
	printLn(unique, "runtime_notifyListAdd start", getTime())
	t := runtime_notifyListAdd(&c.notify)
	printLn(unique, "runtime_notifyListAdd end", getTime())
	printLn(unique, "c.L.Unlock() start", getTime())
	c.L.Unlock()
	printLn(unique, "c.L.Unlock() end", getTime())
	printLn(unique, "runtime_notifyListWait start", getTime())
	runtime_notifyListWait(&c.notify, t)
	printLn(unique, "runtime_notifyListWait end", getTime())
	Sleep(50)
	printLn(unique, "c.L.Lock() start", getTime())
	Sleep(50)
	c.L.Lock()
	printLn(unique, "c.L.Lock() end", getTime())
}
*/

func main() {
	fmt.Println(runtime.GOMAXPROCS(0))
	c := sync.NewCond(&sync.Mutex{})
	v := new(int)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go consume(c, v, "one", wg)
	wg.Add(1)
	go consume(c, v, "two", wg)
	wg.Add(1)
	go produce(c, v, wg)
	wg.Wait()
}

func GetTime() string {
	return strconv.Itoa(int(time.Now().Second()))
}

func Println(a ...interface{}) {
	fmt.Println(a...)
}

func Sleep(milli int) {
	time.Sleep(time.Duration(milli) * time.Millisecond)
}

func consume(c *sync.Cond, v *int, unique string, wg *sync.WaitGroup) {
	fmt.Println(unique, "started")
	fmt.Println(unique, "c.L.Lock() base start", *v, GetTime())
	time.Sleep(50 * time.Millisecond)
	c.L.Lock()
	fmt.Println(unique, "c.L.Lock() base end", *v, GetTime())
	for {
		fmt.Println(unique, "waiting", *v, GetTime())
		c.Wait(unique, GetTime, Println, Sleep)
		fmt.Println(unique, "waited", *v, GetTime())
	}
	fmt.Println("V is:", *v)
	c.L.Unlock()
	wg.Done()
}

func produce(c *sync.Cond, v *int, wg *sync.WaitGroup) {
	for {
		fmt.Println("====================")
		*v += 1
		time.Sleep(1000 * time.Millisecond)
		c.Broadcast()
	}
	wg.Done()
}
