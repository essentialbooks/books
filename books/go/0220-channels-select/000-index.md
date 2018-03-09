---
Title: Channels and select
Id: 141
SOId: 1263
---
A channel is a thread-safe queue of values of a given type.

A primary use for channels is to communicate between goroutines.

For that reason we talk about sending values to a channel (`ch <- value`) and receving values from a channel (`value <- ch`).

Basic of channels:

@file index.go output sha1:5c34ffb9f92451fe53fcfb05de6b620b2afb31ee goplayground:4XuAXmrbElQ

A [zero value](29) of a channel is `nil` so the first thing to do is to create a channel with `make(chan ${type})`.

Send operator `chan <- value` enqueues value at the end.

If channel is full, `<-` will block.

Send on a `nil` channel blocks forever.

Retrieve statement `value = <- chan` picks up the value from the front of the queue.

If channel is empty, retrieve will block.

Another way to retrieve a value form channel is to use `select` statement.

Using `select` allows to wait on multiple channels, do a non-blocking wait and implement [timeouts](143).

Yet another is to use [range](142).

Channels have a fixed capacity.

Channel created with `make(chan bool)` is called unbuffered channel. Send on unbuffered channel blocks until a corresponding receive.

Channel created with `make(chan int, 3)` is a channel of integers with capacity of 3. It's called a buffered channel.

The first 3 sends will finish immediately, the 4th will block until a value is recieved from a channel.

You can [close](144) a channel with `close(chan)`.

Closing channel twice [panics](131).

Sending to closed channel [panics](131).

A receive from closed channels returns zero value immediately.

A closed channel finishes `range` over a channel.
