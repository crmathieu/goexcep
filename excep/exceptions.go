package goexcep

import (
	"fmt"
	"strings"
	"strconv"
)
const (
	THROW_STR = "THROW"
	EXCEP_RUNTIME = -1
	EXCEP_UNKNOWN = -2
	EXCEP_RETHROW = 1
	EXCEP_CUSTOM1 = 10
)
type goexcep struct {
	e      chan int
	code   int
	errmsg string
}

// NewGoexcep - create an exception object
func NewGoexcep() *goexcep {
	return &goexcep{e: make(chan int, 1), code: EXCEP_RUNTIME, errmsg: ""}
}

// Throw - Throws an exception
func Throw(msg string, code int) {
	panic(fmt.Sprintf("%v:%v:%v", THROW_STR, code, msg))
}

// TryAndCatch - performs a try and returns a boolean
func (g *goexcep) TryAndCatch(f func()) bool {
	g.try(f)
	return g.catch()
}

// GetErrorCode
func (g *goexcep) GetErrorCode() int {
	return g.code
}

// GetError
func (g *goexcep) GetError() string {
	return g.errmsg
}

// try - will try a function and recover from an exception if something
// happens during its execution
func (g *goexcep) try(f func()) {
	defer func() {
		if r := recover(); r != nil {
			var t string
			// we are recovering from a panic
			if err, ok := r.(error); ok {
				t = err.Error()
			} else{
				t = fmt.Sprintf("%v", r)
			}
			var err error
			tok := strings.Split(t, ":")
			if tok[0] == THROW_STR {
				g.code, err = strconv.Atoi(tok[1])
				if err != nil {
					g.code = EXCEP_UNKNOWN
				}
				g.errmsg = tok[2]
			} else {
				g.errmsg = t
			}
			fmt.Printf("Recovering from (%v)\n", g.errmsg)
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

func (g *goexcep) Catch() bool {
	return g.catch()
}

// catch - will listen to the exception channel waiting for an exception to
// occur -or- the end of the normal execution
func (g *goexcep) catch() bool {
	excep := <-g.e
	if excep != 0 {
		return true
	}
	return false
}
