package main

import (
	"fmt"
	goe "github.com/crmathieu/goexcep/excep"
)

// triggers a runtime error (division by 0)
func runtime() {
	a, b := 1, 0
	c := a / b
	fmt.Println(c)
}

// triggers a throw for a particular reason
func letitthrow() {
	goe.Throw("let's throw an exception")
}

// nicely behaving function
func goodboy() {
	fmt.Println("It's all good...")
}

// nested exception
func nestedProblems() {
	var e2 = goe.NewGoexcep()
	if err := e2.TryAndCatch(letitthrow); err != nil {
		// catch code
		fmt.Println("Caught from inner try catch:",err.Error())
		goe.Throw(fmt.Sprintf("Re-Throwning %v", err.Error()))
	}
	
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

	if err := e.TryAndCatch(nestedProblems); err != nil {
		// catch code
		fmt.Println("Caught:",err.Error())
	}

}
