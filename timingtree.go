package timingtree

import (
	"bytes"
	"fmt"
	"time"
)

// DefaultNodePrintLimit prevents the String output from being too large
// String will makes sure it includes the slowest child, but only up to this number
// of children.
var DefaultPrintChildLimit = 50

type Node struct {
	Name string

	children  []*Node
	ended     bool
	startTime time.Time
	duration  time.Duration
}

// Start creates a new timer tree if enabled is true, otherwise nil
func Start(name string, enabled bool) *Node {
	if enabled {
		return &Node{Name: name, startTime: time.Now()}
	}
	return nil
}

// StartChild adds a child node to the current node and starts a sub-timer
func (n *Node) StartChild(childName string) *Node {
	if n == nil {
		return nil
	}
	if n.ended {
		panic("Cannot start a child on a node that is already ended")
	}
	child := Start(childName, true)
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
	return n.LimitString(DefaultPrintChildLimit)
}

// LimitString outputs our timing tree as a human-readable string, but only including up to
// childCountLimit number of child nodes, making sure to include the slowest child
func (n Node) LimitString(childCountLimit int) string {
	if !n.ended {
		panic("Cannot print a node that hasn't ended")
	}
	buf := &bytes.Buffer{}
	n.appendString(0, buf, childCountLimit)
	return buf.String()
}

func (n Node) appendString(nesting int, buf *bytes.Buffer, childCountLimit int) {
	for i := 0; i < nesting; i++ {
		buf.WriteRune('\t')
	}
	buf.WriteString(n.Name)
	buf.WriteString(": ")
	buf.WriteString(n.duration.String())
	if len(n.children) > 0 && childCountLimit != 0 {
		children := n.children
		var others int
		if childCountLimit >= 0 && len(children) > childCountLimit {
			var (
				max  time.Duration
				maxI int
			)
			children = n.children[0:childCountLimit]
			for i, child := range n.children {
				if child.duration > max {
					max = child.duration
					maxI = i
				}
			}
			if maxI >= childCountLimit {
				children[0] = n.children[maxI]
			}
			others = len(n.children) - childCountLimit
		}
		nesting++
		if others > 0 {
			buf.WriteRune('\n')
			for i := 0; i < nesting; i++ {
				buf.WriteRune('\t')
			}
			fmt.Fprintf(buf, "** %d children truncated **", others)
		}
		for _, child := range children {
			buf.WriteRune('\n')
			child.appendString(nesting, buf, childCountLimit)
		}
	}
}
