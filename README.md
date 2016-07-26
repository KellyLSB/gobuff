# GoBuff

Is a utility for dealing with memory mapped buffers of data.
Long term goal is to make this compatible as a generalized storage API,
for translating and routing information in models that can be utilized
by object oriented language structures with convenient data back-end switches.

These Buffers are io.ReadWriter compatible.

```
package main

import "github.com/KellyLSB/gobuff"

func main() {
  // Create a Buff with a length of 6
  buff := buff.NewBuff(6)
}
```
