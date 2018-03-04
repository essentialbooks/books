package main

import (
	"fmt"
	"runtime"
)

// :show start
type Foo struct {
	Is []int
}

func (fp *Foo) Panic() (err error) {
	defer PanicRecovery(&err)
	fp.Is[0] = 5
	return nil
}

func PanicRecovery(err *error) {

	if r := recover(); r != nil {
		if _, ok := r.(runtime.Error); ok {
			//fmt.Println("Panicing")
			//panic(r)
			*err = r.(error)
		} else {
			*err = r.(error)
		}
	}
}

func main() {
	fp := &Foo{}
	if err := fp.Panic(); err != nil {
		fmt.Printf("Error: %v", err)
	}
	fmt.Println("ok")
}

// :show end
