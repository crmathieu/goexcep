# goexcep
exception in Go for the nostalgic

Go does not have exceptions. Go allows functions to return an error type in addition to a result via its support for multiple return values. By declaring an error return value you indicate to the caller that this method could go wrong. If a function returns a value and an error, then you can’t assume anything about the value until you’ve inspected the error. 

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

Having said that, we are going to define our _exception_ object as:
```go
type goexcep struct {
	e      chan int
	excep  bool
	errmsg string
}
```

- The channel **e** is used to synchronize the _try_ and _catch_ private methods. 
- The boolean **excep** is used to specify whether an exception occured or not.
- The string **errmsg** is used to hold the error message corresponding to the exception.

The code to wrap with our exception handling is provided as a function to the **TryAndCatch** method. If the method returns an error, it means that either a runtime exception or an exception purposely thrown has been detected, and it is your job to do something about it.






Running the test program will generate the following messages

```text
Recovering from runtime error: integer divide by zero
Caught: runtime error: integer divide by zero
It's all good...
Recovering from let's throw an exception
Caught: let's throw an exception
Recovering from let's throw an exception
Caught: let's throw an exception
Recovering from Thrown from ComplexStuff
Caught: Thrown from ComplexStuff
```
