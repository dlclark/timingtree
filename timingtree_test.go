package timingtree

import (
	"regexp"
	"testing"
	"time"
)

//Some root: 1.295717ms
//	a child: 1.29517ms
var basic = regexp.MustCompile("Some root: 1\\.[0-9]+ms\n\ta child: 1\\.[0-9]+ms")

func TestBasic(t *testing.T) {
	n := Start("Some root")
	c := n.StartChild("a child")
	time.Sleep(time.Millisecond)
	c.End()
	n.End()
	if !basic.MatchString(n.String()) {
		t.Fatalf("Expected match but got %v", n.String())
	}
}

func TestGroupEnd(t *testing.T) {
	n := Start("Some root")
	n.StartChild("a child")
	time.Sleep(time.Millisecond)
	// make sure our parent can end our children
	n.End()
	if !basic.MatchString(n.String()) {
		t.Fatalf("Expected match but got %v", n.String())
	}
}

func TestCheckEndedDuration(t *testing.T) {
	n := Start("Some root")
	n.StartChild("a child")
	time.Sleep(time.Millisecond)
	// make sure our parent can end our children
	n.End()
	if got := n.Duration(); got > 2*time.Millisecond || got < time.Millisecond {
		t.Fatalf("Expected total time between 1-2 milli, got %v", got.String())
	}
}

func TestSlidingDuration(t *testing.T) {
	n := Start("Some root")
	n.StartChild("a child")
	// sleep then check again, time should be moving on
	time.Sleep(time.Millisecond)
	if got := n.Duration(); got > 2*time.Millisecond || got < time.Millisecond {
		t.Fatalf("Expected first time between 1-2 milli, got %v", got.String())
	}
	// sleep then check again, time should be moving on
	time.Sleep(time.Millisecond)
	if got := n.Duration(); got > 3*time.Millisecond || got < 2*time.Millisecond {
		t.Fatalf("Expected second time between 2-3 milli, got %v", got.String())
	}
	// now end and sleep, time should be still
	n.End()
	time.Sleep(time.Millisecond)
	if got := n.Duration(); got > 3*time.Millisecond || got < 2*time.Millisecond {
		t.Fatalf("Expected final time between 2-3 milli, got %v", got.String())
	}
}

func TestSlidingDurationEndedChild(t *testing.T) {
	// confirm it works if we end a child but our parent isn't ended

	n := Start("Some root")
	c := n.StartChild("a child")
	// sleep then check again, time should be moving on
	time.Sleep(time.Millisecond)
	if got := n.Duration(); got > 2*time.Millisecond || got < time.Millisecond {
		t.Fatalf("Expected first root time between 1-2 milli, got %v", got.String())
	}
	if got := c.Duration(); got > 2*time.Millisecond || got < time.Millisecond {
		t.Fatalf("Expected first child time between 1-2 milli, got %v", got.String())
	}
	c.End()
	// sleep then check again, time should be moving on still with parent but not child
	time.Sleep(time.Millisecond)
	if got := n.Duration(); got > 3*time.Millisecond || got < 2*time.Millisecond {
		t.Fatalf("Expected second root time between 2-3 milli, got %v", got.String())
	}
	if got := c.Duration(); got > 2*time.Millisecond || got < time.Millisecond {
		t.Fatalf("Expected second child time between 1-2 milli, got %v", got.String())
	}
	// now end parent and sleep, time should be still for both
	n.End()
	time.Sleep(time.Millisecond)
	if got := n.Duration(); got > 3*time.Millisecond || got < 2*time.Millisecond {
		t.Fatalf("Expected final root time between 2-3 milli, got %v", got.String())
	}
	if got := c.Duration(); got > 2*time.Millisecond || got < time.Millisecond {
		t.Fatalf("Expected final child time between 1-2 milli, got %v", got.String())
	}
}

func TestDeferedNodeEnd(t *testing.T) {
	n := Start("Some root")
	func() {
		defer n.StartChild("a child").End()
		time.Sleep(time.Millisecond)
	}()
	// make sure our root goes beyond 2 millis
	time.Sleep(time.Millisecond)
	if got := n.children[0].duration; got > 2*time.Millisecond || got < time.Millisecond {
		t.Fatalf("Expected first root time between 1-2 milli, got %v", got.String())
	}
}

func TestNilNode(t *testing.T) {
	// none of this should throw -- this allows a simple niling of the root
	// and subsequent calls become cheaper, just don't output the timing tree!
	var n *Node
	ch := n.StartChild("childName")

	if ch != nil {
		t.Fatalf("Expected nil child, got %v", ch)
	}
	ch.End()
	if want, got := time.Duration(0), n.Duration(); want != got {
		t.Fatalf("timing wanted %v, got %v", want, got)
	}
}

var multi = regexp.MustCompile("Some root: 2\\.[0-9]+ms\n\ta child: 1\\.[0-9]+ms\n\tanother child: 1\\.[0-9]+ms")

func TestMutliChildrenString(t *testing.T) {
	n := Start("Some root")
	ch := n.StartChild("a child")
	time.Sleep(time.Millisecond)
	ch.End()
	ch = n.StartChild("another child")
	time.Sleep(time.Millisecond)
	ch.End()
	n.End()
	if !multi.MatchString(n.String()) {
		t.Fatalf("Expected match but got %v", n.String())
	}
}
