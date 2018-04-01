---
Title: Interface & Polymorphism
Id: 284
Score: 1
SOId: 27429
---
Interfaces provide a way to specify the behaviour of an object, if something can do this then it can be used here. an interface defines a set of methods, but these methods do not contain code as they are abstract or the implementation is left to the user of the interface. unlike most Object Oriented languages interfaces can contain variables in Go.

Polymorphism is the essence of object-oriented programming: the ability to treat objects of different types uniformly as long as they adhere to the same interface. Go interfaces provide this capability in a very direct and intuitive way

```go
package main

type Runner interface {
    Run()
}

type Admin struct {
    Username, Password string
}

func (admin Admin) Run() {
    fmt.Println("Admin ==> Run()");
}

type User struct {
    ID uint64
    FullName, Email string
}

func (user User) Run() {
    fmt.Println("User ==> Run()")
}

// RunnerExample takes any type that fullfils the Runner interface
func RunnerExample(r Runner) {
    r.Run()
}

func main() {
    admin := Admin{
        "zola",
        "supersecretpassword",
    }

    user := User{
        1,
        "Zelalem Mekonen",
        "zola.mk.27@gmail.com",
    }

    RunnerExample(admin)

    RunnerExample(user)
}
```
