# goexcep
exception in Go for the nostalgic.

## Installation
Run ```go get https://github.com/crmathieu/goexcep``` and import the package at the beginning of your code:
```go
import (
    goe "github.com/crmathieu/goexcep/excep"
)
```

## Introduction
For those developers used to work in C++, Java, PHP, Python etc... it might be a little hard to realize that Go does not have exceptions. But there are ways where we can sort of simulate the type of exception handling known in other languages.


Go allows functions to return an error type in addition to a result via its support for multiple return values. By declaring an error return value you indicate to the caller that this method could go wrong. If a function returns a value and an error, then you can’t assume anything about the value until you’ve inspected the error. 

Now, let's examine a few go functionalities that we will be using in our exception implementation:

From the Go documentation, we learned that:

A **defer** statement pushes a function call onto a list. The list of saved calls is executed after the surrounding function returns. Defer is commonly used to simplify functions that perform various clean-up actions. A few things to keep in mind:
- Deferred function's arguments are evaluated when the defer statement is evaluated. In this example, the expression "i" is evaluated when the Println call is deferred. The deferred call will print "0" after the function returns.
```go
func a() {
    i := 0
    defer fmt.Println(i)
    i++
    return
}
```
- Deferred function calls are executed in Last In First Out order after the surrounding function returns. This function prints "3210":
```go
func b() {
    for i := 0; i < 4; i++ {
        defer fmt.Print(i)
    }
}
```
- Deferred functions may read and assign to the returning function's named return values.In this example, a deferred function increments the return value i after the surrounding function returns. Thus, this function returns 2:
```go
func c() (i int) {
    defer func() { i++ }()
    return 1
}
```

**Panic** is a built-in function that stops the ordinary flow of control and begins panicking. When the function F calls panic, execution of F stops, any deferred functions in F are executed normally, and then F returns to its caller. To the caller, F then behaves like a call to panic. The process continues up the stack until all functions in the current goroutine have retu    rned, at which point the program crashes. Panics can be initiated by invoking panic directly. They can also be caused by runtime errors, such as out-of-bounds array accesses.

**Recover** is a built-in function that regains control of a panicking goroutine. Recover is only useful inside deferred functions. During normal execution, a call to recover will return nil and have no other effect. If the current goroutine is panicking, a call to recover will capture the value given to panic and resume normal execution.

Having said that, we are going to define an _exception_ object as:
```go
type goexcep struct {
    e      chan int
    errmsg string
}
```
Where
- The channel **e** is used to synchronize a _try_ and _catch_ private methods with the exception type (or lack thereof). 
- The string **errmsg** is used to hold the error message corresponding to the exception.

A private _try_ method takes a function as a parameter. This function corresponds to the code we want to wrap and protect from errors with our exception handling system. The _try_ method has a differed function that recovers from a _panic_ call, whether it is generated by the go runtime, or by the code itself using the _Throw_ method.

```go
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
            if tok[0] == "THROW" {
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
```
Here we first check if we are recovering from a call to **panic**. When this is the case (test is true), then we extract the error message and notify the exception channel. We also filter whether the exception originates from a **throw** and extract the error code if this is the case.

A private _catch_ method waits for the exception channel to be unblocked. This is accomplished when an exception occurs from the recovery of a call to panic, or when the code finishes normally.

```go
func (g *goexcep) catch() bool {
    excep := <-g.e
    if excep != 0 {
        return true
    }
    return false
}
```
_catch_ returns a boolean that is true when an exception occured.

The API function **TryAndCatch** calls the _try_ and _catch_ methods and returns a boolean.
```go
func (g *goexcep) TryAndCatch(f func()) bool {
    g.try(f)
    return g.catch()
}
```

```go
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
            if tok[0] == "THROW" {
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
```

## API

#### Create an exception object
```go
func NewGoexcep() *goexcep
```

#### Throw an exception
```go
func Throw(msg string, errcode int) 
```

#### Try
```go
func (g *goexcep) Try(f func())
```

#### Catch
```go
func (g *goexcep) Catch() error 
```

#### or Try and Catch in one call
```go
func (g *goexcep) TryAndCatch(f func()) error 
```

#### GetErrorCode
```go
func (g *goexcep) GetErrorCode() int 
```

#### GetError
```go
func (g *goexcep) GetError() string 
```

## Examples
The following illustrates the capture of different type of exceptions (runtime, code generated) as well as an example of nested exceptions

#### runtime error (division by 0)
```go
func divByZero() {
    a, b := 1, 0
    c := a / b
    fmt.Println(c)
}
```

#### exception thrown
```go
func letitthrow() {
    goe.Throw("let's throw an exception of type CUSTOM1", goe.EXCEP_CUSTOM1)
}
```

#### nicely behaving function
```go
func goodboy() {
    fmt.Println("It's all good...")
}
```

#### runtime error (memory violation)
```go
func segViolation() {
    var p *int
    *p = 1
}
```

#### nested exception
```go
func nestedProblems() {
    var e2 = goe.NewGoexcep()
    if e2.TryAndCatch(letitthrow) {
        // catch code
        fmt.Printf("Caught in 'letitthrow' from inner try catch (%v)\n", e2.GetError())
        goe.Throw(fmt.Sprintf("Re-Throwning (%v)", e2.GetError()), e2.GetErrorCode())
    }
}
```

#### index range
```go
func indexRange() {
    x := []int{1,2} 
    
    for i:=0;i<5;i++ {
        fmt.Println(x[i])
    }
}
```

#### stack overflow
```go
func stackOverflow() {
    stackOverflow()
}
```

#### deeper function
```go
func deeper() {
    indexRange()
    fmt.Println("end")
}
```
#### function with goroutine
If your function creates go subroutines, each subroutine will operate on its own stack and will be out of scope as far as the current **TryAndCatch** function call is concerned. To make this work, each subroutine MUST create their own exception object and call their own **TryAndCatch** function, like in the example below:

```go
func withSubroutine() {	
    go func() {
        var e2 = goe.NewGoexcep()
        if e2.TryAndCatch(segViolation) {
            fmt.Printf("Caught in goroutine 'segViolation' (%v)\n", e2.GetError())
        }
    }()
    divByZero()
}
```
**withSubroutine** will generate 2 exceptions: one during the execution of the _divByZero_ function and one inside the goroutine originating from the _segViolation_ function.

#### putting everything together
```go
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
```

#### The above code returns the following messages
```text
1
2
Recovering from (runtime error: index out of range)
Caught in 'anonymous' (runtime error: index out of range) - Code (-1)
Recovering from (runtime error: integer divide by zero)
Caught in 'withSubroutine' (runtime error: integer divide by zero)
Recovering from (runtime error: integer divide by zero)
Caught in 'divByZero' (runtime error: integer divide by zero) - code (-1)
It's all good...
Recovering from (let's throw an exception of type CUSTOM1)
Caught in 'letitthrow' from inner try catch (let's throw an exception of type CUSTOM1)
Recovering from (Re-Throwning (let's throw an exception of type CUSTOM1))
Caught in 'nestedProblems' (Re-Throwning (let's throw an exception of type CUSTOM1))
```

Because the deferred block is defined at the _try_ goroutine block level, a panic generated within the function provided as a parameter will bubble up from its origin in the call stack until it reaches the TryAndCatch code. This, in turns, triggers a call to the deferred function which captures _panic_ using the _recover_ function.

For that reason, we can see for example that in the **deeper** function, the instruction to display the message **end** never gets a chance to be executed simply because when the runtime error happens within the **indexRange** function, it bubbles up to the **deeper** function body and since there is no differed block with recovery yet, it continues on bubbling up the call stack. Hence any code following the function call where an error originated is ignored.  

**Note**
This exception mechanism will NOT work for _Stack Overflow_ errors. Panic is bubbling up using the call stack until it reaches the **main** or a **recover** function. If the stack is corrupted, this can't work and you will get a **fatal error** and a stack trace.
