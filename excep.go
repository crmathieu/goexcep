package goexcep

import (
	"fmt"
	"errors"
)

type goexcep struct {
	e 		chan int
	excep 	bool
	errmsg 	string
}

func NewGoexcep() *goexcep {
	return &goexcep{e: make(chan int), excep: false, errmsg: ""}
}

func (g *goexcep) try (f func()) {
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

func Throw(msg string) {
	panic(fmt.Sprintf("%v", msg))
}

func (g *goexcep) TryAndCatch (f func()) error {
	g.try(f)
	return g.catch()
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