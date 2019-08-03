package main

import (
	"fmt"
	goe "github.com/crmathieu/Goexcep/excep"
)

// runtime error (division by 0)
func divByZero(e *goe.Goexcep) {
	a, b := 1, 0
	c := a / b
	fmt.Println(c)
}

// exception thrown
func letitthrow(e *goe.Goexcep) {
	goe.Throw("let's throw an exception of type CUSTOM1", goe.EXCEP_CUSTOM1)
}

// nicely behaving function
func goodboy(e *goe.Goexcep) {
	fmt.Println("It's all good...")
}

// segment violation
func segViolation(e *goe.Goexcep) {
	var p *int
	*p = 1
}

// nested exception
func nestedProblems(e *goe.Goexcep) {
	if e.TryAndCatch(letitthrow) {
		// catch code
		fmt.Printf("Caught in 'letitthrow' from inner try catch (%v)\n", e.GetError())
		goe.Throw(fmt.Sprintf("Re-Throwning (%v)", e.GetError()), e.GetErrorCode())
	}
}

// indexRange
func indexRange(e *goe.Goexcep) {
	x := []int{1,2} 

	for i:=0;i<5;i++ {
		fmt.Println(x[i])
	}
}


// deeper 
func deeper(e *goe.Goexcep) {
    indexRange(e)
    fmt.Println("end")
}


// with subroutine
func withSubroutine(e *goe.Goexcep) {	
	go func() {
		e.Try(func(*goe.Goexcep) {
			segViolation(e)
		})
		if e.Catch() {
			fmt.Printf("Caught in 'withSubroutine' goroutine (%v)\n", e.GetError())
		}
	}()
	divByZero(e)
}

func main() {
	e := goe.NewGoexcep()

	// one way to do it
	e.Try(func(*goe.Goexcep) {
		indexRange(e)
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
