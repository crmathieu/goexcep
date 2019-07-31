package main

import (
	"fmt"
	goe "github.com/crmathieu/goexcep/excep"
)

// runtime error (division by 0)
func runtime() {
	a, b := 1, 0
	c := a / b
	fmt.Println(c)
}

// exception thrown
func letitthrow() {
	goe.Throw("let's throw an exception")
}

// nicely behaving function
func goodboy() {
	fmt.Println("It's all good...")
}

// segment violation
func segViolation() {
	var p *int
	*p = 1
}

// nested exception
func nestedProblems() {
	var e2 = goe.NewGoexcep()
	if err := e2.TryAndCatch(letitthrow); err != nil {
		// catch code
		fmt.Printf("Caught from inner try catch (%v)\n", err.Error())
		goe.Throw(fmt.Sprintf("Re-Throwning (%v)", err.Error()))
	}

}

func main() {
	e := goe.NewGoexcep()
	if err := e.TryAndCatch(runtime); err != nil {
		// catch code
		fmt.Printf("Caught (%v)\n", err.Error())
	}
	if err := e.TryAndCatch(goodboy); err != nil {
		// catch code
		fmt.Printf("Caught (%v)\n", err.Error())
	}
	if err := e.TryAndCatch(segViolation); err != nil {
		// catch code
		fmt.Printf("Caught (%v)\n", err.Error())
	}

	if err := e.TryAndCatch(nestedProblems); err != nil {
		// catch code
		fmt.Printf("Caught (%v)\n", err.Error())
	}

}
