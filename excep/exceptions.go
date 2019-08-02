package goexcep

import (
	"errors"
	"fmt"
)

type goexcep struct {
	e      chan int
	errmsg string
}

// NewGoexcep - create an exception object
func NewGoexcep() *goexcep {
	return &goexcep{e: make(chan int, 1), errmsg: ""}
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
	defer func() {
		if r := recover(); r != nil {
			// we are recovering from a panic
			fmt.Printf("Recovering from (%v)\n", r)
			if err, ok := r.(error); ok {
				g.errmsg = err.Error()
			} else {
				g.errmsg = fmt.Sprintf("%v", r)
			}
			// we exit with an exception - feed the exception channel
			g.e <- 1
		}
	}()
	f()
	// we exit without exception - feed the exception channel
	g.e <- 0
}

func (g *goexcep) Try(f func()) {
	g.try(f)
}

func (g *goexcep) Catch() error {
	return g.catch()
}

// catch - will listen to the exception channel waiting for an exception to
// occur -or- the end of the normal execution
func (g *goexcep) catch() error {
	excep := <-g.e
	if excep != 0 {
		return errors.New(g.errmsg)
	}
	return nil
}
