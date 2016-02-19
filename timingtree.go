package timingtree

import (
	"bytes"
	"time"
)

type Node struct {
	Name string

	children  []*Node
	ended     bool
	startTime time.Time
	duration  time.Duration
}

// Start creates a new timer tree
func Start(name string) *Node {
	return &Node{Name: name, startTime: time.Now()}
}

// StartChild adds a child node to the current node and starts a sub-timer
func (n *Node) StartChild(childName string) *Node {
	if n == nil {
		return nil
	}
	if n.ended {
		panic("Cannot start a child on a node that is already ended")
	}
	child := Start(childName)
	n.children = append(n.children, child)
	return child
}

// End the node and any un-ended children
func (n *Node) End() {
	if n == nil {
		return
	}
	if n.ended {
		panic("Cannot end a node that is already ended")
	}

	for _, child := range n.children {
		if !child.ended {
			child.End()
		}
	}
	n.ended = true
	n.duration = time.Now().Sub(n.startTime)
}

// Duration returns the current duration of the tree
func (n *Node) Duration() time.Duration {
	if n == nil {
		return 0
	}
	if n.ended {
		return n.duration
	}
	return time.Now().Sub(n.startTime)
}

// String outputs our timing tree as a human-readable string
func (n Node) String() string {
	if !n.ended {
		panic("Cannot print a node that hasn't ended")
	}
	buf := &bytes.Buffer{}
	n.appendString(0, buf)
	return buf.String()
}

func (n Node) appendString(nesting int, buf *bytes.Buffer) {
	for i := 0; i < nesting; i++ {
		buf.WriteRune('\t')
	}
	buf.WriteString(n.Name)
	buf.WriteString(": ")
	buf.WriteString(n.duration.String())
	if len(n.children) > 0 {
		buf.WriteRune('\n')
		nesting++
		for _, child := range n.children {
			child.appendString(nesting, buf)
		}
	}
}
