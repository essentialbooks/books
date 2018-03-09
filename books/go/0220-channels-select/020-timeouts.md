---
Title: Timout reading from channel with select
Id: 143
SOId: 6050
---
Receiving from a channel with `<- chan` or `for range` loop blocks.

Sometimes you want to limit time waiting for a value on a channel.

It's possible with `select`:

@file timeout.go output sha1:1a3592bd922445233dfa604d927c9b5ca39fcc46 goplayground:kWFgA_-fMgN
