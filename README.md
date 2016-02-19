# timingtree
Timing tree library for Go program performance diagnostics

# example usage
```go
n := timingtree.Start("root name")

ch := n.StartChild("some op")
SomeOp(ch)
ch.End()

ch = n.StartChild("other op")
OtherOp(ch)
ch.End()

// stop our timing tree and check if we should output debug (>1s tree)
n.End()
if n.Duration() > time.Second {
	fmt.Print(n.String())	
}
```
Will output something like this if the root node time took longer than 1 second:

	Some root: 2.605474s
		some op: 1.337457s
			more code: 0.07s
			expensive: 1.2s
		other op: 1.267011s

There's also a shortform to just time a while function using defer:

```go
n := timingtree.Start("root name")

func() {
	defer n.StartChild("some func").End()
	//... do stuff
}()

n.End()
```

# notes
It is not concurrent so you'll need to wrap with your own locking to use across goroutines.