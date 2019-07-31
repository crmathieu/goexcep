package goexcep

import (
	"errors"
	"fmt"
)

type goexcep struct {
	e      chan int
	excep  bool
	errmsg string
}

// NewGoexcep - create an exception object
func NewGoexcep() *goexcep {
	return &goexcep{e: make(chan int), excep: false, errmsg: ""}
}

// Throw - Throws an exception
func Throw(msg string) {
	panic(fmt.Sprintf("%v", msg))
}

// TryAndCatch - performs a try and returns error when exception is caught 
func (g *goexcep) TryAndCatch(f func()) error {
	g.try(f)
	return g.catch()
}

// try - will try a function and recover from an exception if something
// happens during its execution 
func (g *goexcep) try(f func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// we are recovering from a panic
				fmt.Println("Recovering from", r)
				if err, ok := r.(error); ok {
					g.errmsg = err.Error()
				} else {
					g.errmsg = fmt.Sprintf("%v", r)
				}
				// we exit with an exception - feed the exception channel
				g.excep = true
				g.e <- 1
			}
		}()
		f()
		// we exit without exception - feed the exception channel
		g.excep = false
		g.e <- 1
	}()
}

// catch - will listen to the exception channel waiting for an exception to
// occur -or- the end of the normal execution 
func (g *goexcep) catch() error {
	<-g.e
	if g.excep == true {
		return errors.New(g.errmsg)
	}
	return nil
}

