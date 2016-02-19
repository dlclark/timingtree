package timingtree

import (
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

//Some root: 1.295717ms
//	a child: 1.29517ms
var basic = regexp.MustCompile("^Some root: 1\\.[0-9]+ms\n\ta child: 1\\.[0-9]+ms")

func TestBasic(t *testing.T) {
	n := Start("Some root", true)
	c := n.StartChild("a child")
	time.Sleep(time.Millisecond)
	c.End()
	n.End()
	if !basic.MatchString(n.String()) {
		t.Fatalf("Expected match but got %v", n.String())
	}
}

func TestGroupEnd(t *testing.T) {
	n := Start("Some root", true)
	n.StartChild("a child")
	time.Sleep(time.Millisecond)
	// make sure our parent can end our children
	n.End()
	if !basic.MatchString(n.String()) {
		t.Fatalf("Expected match but got %v", n.String())
	}
}

func TestCheckEndedDuration(t *testing.T) {
	n := Start("Some root", true)
	n.StartChild("a child")
	time.Sleep(time.Millisecond)
	// make sure our parent can end our children
	n.End()
	if got := n.Duration(); got > 2*time.Millisecond || got < time.Millisecond {
		t.Fatalf("Expected total time between 1-2 milli, got %v", got.String())
	}
}

func TestSlidingDuration(t *testing.T) {
	n := Start("Some root", true)
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

	n := Start("Some root", true)
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
	n := Start("Some root", true)
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
	n := Start("root", false)
	ch := n.StartChild("childName")

	if ch != nil {
		t.Fatalf("Expected nil child, got %v", ch)
	}
	ch.End()
	if want, got := time.Duration(0), n.Duration(); want != got {
		t.Fatalf("timing wanted %v, got %v", want, got)
	}
}

func TestMutliChildrenString(t *testing.T) {
	multi := regexp.MustCompile("^Some root: 2\\.[0-9]+ms\n\ta child: 1\\.[0-9]+ms\n\tanother child: 1\\.[0-9]+ms")

	n := Start("Some root", true)
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

func TestLargeChildCount_MoveLargest(t *testing.T) {
	multi := regexp.MustCompile("^Root: 1\\.[0-9]+ms\n\t\\*\\* 50 children truncated \\*\\*\n\tchild 70: 1\\.[0-9]+ms\n\tchild ")
	//Root: 1.411268ms
	//	** 50 children truncated **
	//	child 70: 1.309909ms

	n := Start("Root", true)

	for i := 0; i < 100; i++ {
		ch := n.StartChild("child " + strconv.Itoa(i))
		if i == 70 {
			time.Sleep(time.Millisecond)
		}
		ch.End()
	}
	n.End()
	if !multi.MatchString(n.String()) {
		t.Fatalf("Expected truncated children with #70 at the top, got %v", n.String())
	}
}

func TestLargeChildCount_NotMoveLargest(t *testing.T) {
	multi := regexp.MustCompile("^Root: 1\\.[0-9]+ms\n\t\\*\\* 50 children truncated \\*\\*\n\tchild 0: ")

	n := Start("Root", true)

	for i := 0; i < 100; i++ {
		ch := n.StartChild("child " + strconv.Itoa(i))
		if i == 29 {
			time.Sleep(time.Millisecond)
		}
		ch.End()
	}
	n.End()
	s := n.String()
	if !multi.MatchString(s) {
		t.Fatalf("Expected truncated children with #0 at the top, got %v", s)
	}
	// should be 52 lines (root, 1 truncation notice, 50 children)
	if want, got := 52, len(strings.Split(s, "\n")); want != got {
		t.Fatalf("Want line count %v got %v", want, got)
	}
}

func TestLargeChildCount_TestEdge(t *testing.T) {
	multi := regexp.MustCompile("^Root: 1\\.[0-9]+ms\n\t\\*\\* 70 children truncated \\*\\*\n\tchild 0: ")

	n := Start("Root", true)

	for i := 0; i < 100; i++ {
		ch := n.StartChild("child " + strconv.Itoa(i))
		if i == 29 {
			time.Sleep(time.Millisecond)
		}
		ch.End()
	}
	n.End()
	s := n.LimitString(30)
	if !multi.MatchString(s) {
		t.Fatalf("Expected truncated children with #0 at the top, got %v", s)
	}
	// should be 32 lines (root, 1 truncation notice, 30 children)
	if want, got := 32, len(strings.Split(s, "\n")); want != got {
		t.Fatalf("Want line count %v got %v", want, got)
	}

}

func TestChildCountPrintWithZero(t *testing.T) {
	re := regexp.MustCompile("^Root: [0-9\\.]+[µnm]s$")

	n := Start("Root", true)

	for i := 0; i < 100; i++ {
		n.StartChild("child " + strconv.Itoa(i))
	}
	n.End()

	if got := n.LimitString(0); !re.MatchString(got) {
		t.Fatalf("Expected just root with no children, got %v", got)
	}
}

func TestLargeChildCount_NoLimit(t *testing.T) {
	multi := regexp.MustCompile("^Root: 1\\.[0-9]+ms\n\tchild 0: ")

	n := Start("Root", true)

	time.Sleep(time.Millisecond)
	for i := 0; i < 100; i++ {
		n.StartChild("child " + strconv.Itoa(i))
	}
	n.End()

	if got := n.LimitString(-1); !multi.MatchString(got) {
		t.Fatalf("Expected truncated children with #0 at the top, got %v", got)
	}
	// should be 101 lines
	if want, got := 101, len(strings.Split(n.LimitString(-1), "\n")); want != got {
		t.Fatalf("Want line count %v got %v", want, got)
	}
}

func TestNoChildrenString(t *testing.T) {
	multi := regexp.MustCompile("^Root: 1\\.[0-9]+ms$")

	n := Start("Root", true)

	time.Sleep(time.Millisecond)

	n.End()

	if got := n.String(); !multi.MatchString(got) {
		t.Fatalf("Expected root with no children, got %v", got)
	}
}
