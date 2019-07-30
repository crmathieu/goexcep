package main

import (
	"fmt"
	goe "github.com/crmathieu/goexcep/excep"
)

// f1 will trigger a runtime error (division by 0)
func f1() {
	a, b := 1, 0
	c := a/b
	fmt.Println(c)
}

// trigger a throw for a particular reason
func f2() {
	goe.Throw("something triggering a throw")
}

// nothing happened
func f3() {
	fmt.Println("It's all good...")
}

func main() {
	e := goe.NewGoexcep()
	if err := e.TryAndCatch(f1); err != nil {
		// catch code
		fmt.Println(err.Error())
	}
	if err := e.TryAndCatch(f2); err != nil {
		// catch code
		fmt.Println(err.Error())
	}
	if err := e.TryAndCatch(f3); err != nil {
		// catch code
		fmt.Println(err.Error())
	}

}

