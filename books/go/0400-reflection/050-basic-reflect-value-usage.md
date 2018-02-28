---
Title: Basic reflect.Value Usage
Id: 6070
Score: -1
---
```go
import "reflect"

value := reflect.ValueOf(4)

// Interface returns an interface{}-typed value, which can be type-asserted
value.Interface().(int) // 4

// Type gets the reflect.Type, which contains runtime type information about
// this value
value.Type().Name() // int

value.SetInt(5) // panics -- non-pointer/non-slice/non-array types are not addressable

x := 4
reflect.ValueOf(&x).Elem().SetInt(5) // works
```
