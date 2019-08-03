package main

import (
	"fmt"
	goe "github.com/crmathieu/goexcep/excep"
)

// runtime error (division by 0)
func divByZero() {
	a, b := 1, 0
	c := a / b
	fmt.Println(c)
}

// exception thrown
func letitthrow() {
	goe.Throw("let's throw an exception of type CUSTOM1", goe.EXCEP_CUSTOM1)
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
	if e2.TryAndCatch(letitthrow) {
		// catch code
		fmt.Printf("Caught in 'letitthrow' from inner try catch (%v)\n", e2.GetError())
		goe.Throw(fmt.Sprintf("Re-Throwning (%v)", e2.GetError()), e2.GetErrorCode())
	}
}

// indexRange
func indexRange() {
	x := []int{1,2} 

	for i:=0;i<5;i++ {
		fmt.Println(x[i])
	}
}


// deeper 
func deeper() {
    indexRange()
    fmt.Println("end")
}

// with subroutine
func withSubroutine() {	
	go func() {
		var e2 = goe.NewGoexcep()
		if e2.TryAndCatch(segViolation) {
			fmt.Printf("Caught in goroutine 'segViolation' (%v)\n", e2.GetError())
		}
	}()
	divByZero()
}

func main() {
	e := goe.NewGoexcep()

	// one way to do it
	e.Try(func() {
		indexRange()
		fmt.Println("end")
	})
	if e.Catch() {
		// catch code
		fmt.Printf("Caught in 'anonymous' (%v) - Code (%v)\n", e.GetError(), e.GetErrorCode())
	}

	// and the other way
	if e.TryAndCatch(withSubroutine) {
       	// catch code
        fmt.Printf("Caught in 'withSubroutine' (%v)\n", e.GetError())
 	}
	if e.TryAndCatch(divByZero) {
		// catch code
		fmt.Printf("Caught in 'divByZero' (%v) - code (%v)\n", e.GetError(), e.GetErrorCode())
	}
	if e.TryAndCatch(goodboy) {
		// catch code
		fmt.Printf("Caught in 'goodboy' (%v)\n", e.GetError())
	}
	if e.TryAndCatch(nestedProblems) {
		// catch code
		fmt.Printf("Caught in 'nestedProblems' (%v)\n", e.GetError())
	}
}
