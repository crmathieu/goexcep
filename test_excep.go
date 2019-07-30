package main

import (
	"fmt"
	"github.com/crmathieu/goexcep"
)

type goexcep struct {
	e 		chan int
	excep 	bool
	errmsg 	string
}

func NewGoexcep() *goexcep {
	return &goexcep{e: make(chan int), excep: false, errmsg: ""}
}

func (g *goexcep) Try (f func()) {
   go func() {
		defer func() {
			if r := recover(); r != nil {
				// we are recovering from a panic
				fmt.Println("Recovered in f", r)
				if err, ok := r.(error); ok {
					g.errmsg = err.Error()
				} else {
					g.errmsg = fmt.Sprintf("%v", r)
				}
				g.excep = true
				g.e <- 1
			} 
		}()
		f()
		g.excep = false
		g.e <- 1
   }()		
}

func (g *goexcep) catch () error {
	fmt.Println("Waiting to catch...")
	select {
	case <- g.e: if g.excep == true {
					fmt.Println("Caught")
					return errors.New(g.errmsg)
				 }
	}
	return nil
}

func (g *goexcep) Catch (f func(string)) {
	fmt.Println("Waiting to catch...")
	select {
	case <- g.e: if g.excep == true {
					fmt.Println("Caught")
					f(g.errmsg)
				 }
	}
}

func Throw(msg string) {
	panic(fmt.Sprintf("%v", msg))
}

func (g *goexcep) TryAndCatch (f func()) error {
	g.Try(f)
	return g.catch()
}

// f1 will trigger a runtime error (division by 0)
func f1() {
	a, b := 1, 0
	c := a/b
	fmt.Println(c)
}

// trigger a throw for a particular reason
func f2() {
	Throw("something triggering a throw")
}

// nothin happened
func f3() {
	fmt.Println("It's all good...")
}

func handler(msg string) {
	fmt.Println(msg)
}

var e = NewGoexcep()

func main2() {
	e.Try(f3)
	e.Catch(handler)
}

func main() {
	if err := e.TryAndCatch(f2); err != nil {
		// catch code
		fmt.Println(err.Error())
	}
}

/*
package main

import "fmt"

func main() {
    f()
    fmt.Println("Returned normally from f.")
}

func f() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered in f", r)
        }
    }()
    fmt.Println("Calling g.")
    g(0)
    fmt.Println("Returned normally from g.")
}

func g(i int) {
    if i > 3 {
        fmt.Println("Panicking!")
        panic(fmt.Sprintf("%v", i))
    }
    defer fmt.Println("Defer in g", i)
    fmt.Println("Printing in g", i)
    g(i + 1)
}
*/