---
Title: Checking if channel has data available
Id: 148
SOId: 801000tb
---

Receiving on a channel blocks if there is no data in the channel.

What if you don't want to block?

You might be tempted to check if channel has data before doing a receive.

You can't do that.

For one, it couldn't possibly work correctly.

Between the time you check for availability and the time you do a receive, some other goroutine could have picked up the value.

If you want to avoid infinite waiting, you can add a [timeout](143) or do a [non-blocking wait](146) by using `select`.
