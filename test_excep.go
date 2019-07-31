package main

import (
	"fmt"
	goe "github.com/crmathieu/goexcep/excep"
)

// f1 will trigger a runtime error (division by 0)
func runtime() {
	a, b := 1, 0
	c := a / b
	fmt.Println(c)
}

// f2 will trigger a throw for a particular reason
func letitthrow() {
	goe.Throw("let's throw an exception")
}

// nothing happened
func goodboy() {
	fmt.Println("It's all good...")
}

func complexStuff() {
	var e2 = goe.NewGoexcep()
	if err := e2.TryAndCatch(letitthrow); err != nil {
		// catch code
		fmt.Println("Caught:",err.Error())
	}
	goe.Throw("Thrown from ComplexStuff")
}

func main() {
	e := goe.NewGoexcep()
	if err := e.TryAndCatch(runtime); err != nil {
		// catch code
		fmt.Println("Caught:",err.Error())
	}
	if err := e.TryAndCatch(goodboy); err != nil {
		// catch code
		fmt.Println("Caught:",err.Error())
	}
	if err := e.TryAndCatch(letitthrow); err != nil {
		// catch code
		fmt.Println("Caught:",err.Error())
	}

	if err := e.TryAndCatch(complexStuff); err != nil {
		// catch code
		fmt.Println("Caught:",err.Error())
	}

}
