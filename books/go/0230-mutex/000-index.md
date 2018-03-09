---
Title: Mutex
Id: 149
SOId: 2607
---
In Go goroutines share memory.

It's good for performance but modifying the same memory from multiple goroutines is not safe. It can lead to data races and crashes.

One way to avoid that is by using [channels](141) to transfer exclusive ownership of data.

This is Go's motto: *do not communicate by sharing memory; instead, share memory by communicating.*

Another way to ensure exclusive access to data by goroutines is to use mutexes.

Here's a basic pattern that uses a map for a cache and ensures exclusive access with mutex locking:

@file index.go output sha1:6507399efba68176abf2eeb6e4fc71b0312aa36e goplayground:QlyPp-d8av8

Zero-value of `sync.Mutex` is a valid mutex so you don't need to explicitly initialize it.

For performance we want to minimize the time spent holding a lock.

Function `getCached` would be simpler if we kept the mutex locked for the duration of the function but we don't want to keep cache locked when we're executing `expensiveOperation()`.

Unlike many other languages, Go mutexes are non-recursive.

If the same goroutines tries to `Lock()` a mutex twice, the second `Lock()` will block forever.
