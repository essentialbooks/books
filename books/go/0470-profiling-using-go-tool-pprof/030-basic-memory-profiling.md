---
Title: Basic memory Profiling
Id: 261
Score: 0
SOId: 25407
---

```go
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
    flag.Parse()
    if *memprofile != "" {
        f, err := os.Create(*memprofile)
        if err != nil {
            log.Fatal("could not create memory profile: ", err)
        }
        runtime.GC() // get up-to-date statistics
        if err := pprof.WriteHeapProfile(f); err != nil {
            log.Fatal("could not write memory profile: ", err)
        }
        f.Close()
    }
}
```

```go
go build main.go
main.exe -memprofile mem.prof
go tool pprof main.exe mem.prof
```
