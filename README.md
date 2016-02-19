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

n.End()
if n.Duration() > time.Second {
	log.Print(n.String())	
}
```

```go
n := timingtree.Start("root name")

func() {
	defer n.StartChild("some func").End()
	//... do stuff
}()

// stop our timing tree and check if we should output debug (>1s tree)
n.End()
if n.Duration() > time.Second {
	log.Print(n.String())	
}
```

# notes
It is not concurrent so you'll need to wrap with your own locking to use across goroutines.