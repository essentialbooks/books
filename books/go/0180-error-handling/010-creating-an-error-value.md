---
Title: Creating errors
Id: 2705
---

There are several ways to create an error value.

## Use errors.New

The simplest way, `errors.New(msg string)` creates error value from string.

```go
return errors.New("error created with errors.New")
```

## Use fmt.Errorf

For more complex error messages, `fmt.Errorf(format string, args ...interface{})` takes a formatting string and arguments:

```go
return fmt.Errorf("error created with %s", "fmt.Errorf")
```

It's a shortcut for `errors.New(fmt.Sprintf(...))`.

## Use global error variable

```go
// ErrGlobal exported so that others can compare returned error value with this variable
var ErrGlobal = errors.New("global error variable")

return ErrGlobal
```

Sometimes you want the error to have an identity so that callers can test if returned error is this specific error. You can do it by declaring global variable as `ErrGlobal` in example above.

One example of such error from standard library is `io.EOF` although usually the naming convention for such errors is `Err*`.

## Nil indicates no error

To indicate there was no error, return nil.


Finally, you can [create custom error type](2706).

