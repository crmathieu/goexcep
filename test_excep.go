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
		fmt.Printf("Caught in 'letitthrow' from inner try catch (%v)\n",err.Error())
		goe.Throw(fmt.Sprintf("Re-Throwning (%v)", err.Error()))
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
		if err := e2.TryAndCatch(segViolation); err != nil {
			fmt.Printf("Caught in goroutine 'segViolation' (%v)\n",err.Error())
		}
	}()
	divByZero()
}

func main() {
	e := goe.NewGoexcep()

	// one way to do it
	e.Try(deeper)
	if err := e.Catch(); err != nil {
		// catch code
		fmt.Printf("Caught in 'deeper' (%v)\n",err.Error())
	}

	// and the other way
	if err := e.TryAndCatch(withSubroutine); err != nil {
       // catch code
        fmt.Printf("Caught in 'withSubroutine' (%v)\n",err.Error())
 	}
	if err := e.TryAndCatch(divByZero); err != nil {
		// catch code
		fmt.Printf("Caught in 'divByZero' (%v)\n",err.Error())
	}
	if err := e.TryAndCatch(goodboy); err != nil {
		// catch code
		fmt.Printf("Caught in 'goodboy' (%v)\n",err.Error())
	}
	if err := e.TryAndCatch(nestedProblems); err != nil {
		// catch code
		fmt.Printf("Caught in 'nestedProblems' (%v)\n",err.Error())
	}
}
