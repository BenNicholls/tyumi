package util

import "testing"

// testTree is a small interface wrapping TreeType, mirroring how real consumers of TreeNode (e.g. gfx/ui's
// ElementInterface) use an interface type as the tree's type parameter rather than a concrete pointer type. This
// matters because TreeNode's parent-nil checks rely on comparing an interface value to nil, which only behaves as
// expected when the type parameter itself is an interface.
type testTree interface {
	TreeType[testTree]
	Name() string
}

type testNode struct {
	TreeNode[testTree]

	name string
}

func (n *testNode) Name() string { return n.name }

func newTestNode(name string) testTree {
	n := &testNode{name: name}
	n.Init(n)
	return n
}

func TestTreeAddChild(t *testing.T) {
	root := newTestNode("root")
	child := newTestNode("child")

	root.AddChild(child)

	if root.ChildCount() != 1 {
		t.Errorf("ChildCount() = %d, want 1", root.ChildCount())
	}

	if child.GetParent() != root {
		t.Error("child's parent should be root")
	}

	if root.GetChildren()[0] != child {
		t.Error("root's child should be the added child")
	}
}

func TestTreeAddChildren(t *testing.T) {
	root := newTestNode("root")
	c1, c2, c3 := newTestNode("c1"), newTestNode("c2"), newTestNode("c3")

	root.AddChildren(c1, c2, c3)

	if root.ChildCount() != 3 {
		t.Errorf("ChildCount() = %d, want 3", root.ChildCount())
	}
}

func TestTreeAddChildTwiceIsNoOp(t *testing.T) {
	root := newTestNode("root")
	child := newTestNode("child")

	root.AddChild(child)
	root.AddChild(child)

	if root.ChildCount() != 1 {
		t.Errorf("ChildCount() = %d, want 1 (duplicate add should be no-op)", root.ChildCount())
	}
}

// NOTE: reparenting via AddChild updates the child's parent pointer and adds it to the new parent's children, but
// (via addChildNode's call to node.Deparent()) does not remove the child from the *old* parent's children slice -
// that only happens through an explicit RemoveChild call. This test documents that current behaviour rather than an
// idealized one, since fixing it would change the public API's observable behaviour.
func TestTreeReparenting(t *testing.T) {
	root1 := newTestNode("root1")
	root2 := newTestNode("root2")
	child := newTestNode("child")

	root1.AddChild(child)
	root2.AddChild(child)

	if root1.ChildCount() != 1 {
		t.Errorf("root1.ChildCount() = %d, want 1 (stale reference remains after reparenting)", root1.ChildCount())
	}
	if root2.ChildCount() != 1 {
		t.Errorf("root2.ChildCount() = %d, want 1", root2.ChildCount())
	}
	if child.GetParent() != root2 {
		t.Error("child's parent should now be root2")
	}
}

func TestTreeRemoveChild(t *testing.T) {
	root := newTestNode("root")
	c1, c2 := newTestNode("c1"), newTestNode("c2")
	root.AddChildren(c1, c2)

	root.RemoveChild(c1)

	if root.ChildCount() != 1 {
		t.Errorf("ChildCount() = %d, want 1", root.ChildCount())
	}
	if root.GetChildren()[0] != c2 {
		t.Error("remaining child should be c2")
	}
	var nilNode testTree
	if c1.GetParent() != nilNode {
		t.Error("removed child should have nil parent")
	}
}

func TestTreeDeparent(t *testing.T) {
	root := newTestNode("root")
	child := newTestNode("child")
	root.AddChild(child)

	child.Deparent()

	var nilNode testTree
	if child.GetParent() != nilNode {
		t.Error("expected child's parent to be nil after Deparent")
	}
	// NOTE: Deparent only clears the child's parent pointer, it does not remove the child from the parent's
	// children slice - that bookkeeping is done by TreeNode.RemoveChild.
	if root.ChildCount() != 1 {
		t.Errorf("root.ChildCount() = %d, want 1 (Deparent alone doesn't update the parent's children)", root.ChildCount())
	}
}

func TestTreeGetSelf(t *testing.T) {
	root := newTestNode("root")
	if root.GetSelf() != root {
		t.Error("GetSelf() should return the node itself")
	}
}

func TestWalkTree(t *testing.T) {
	root := newTestNode("root")
	c1 := newTestNode("c1")
	c2 := newTestNode("c2")
	gc1 := newTestNode("gc1")

	root.AddChildren(c1, c2)
	c1.AddChild(gc1)

	var visited []string
	WalkTree(root, func(n testTree) {
		visited = append(visited, n.Name())
	})

	if len(visited) != 4 {
		t.Fatalf("WalkTree visited %d nodes, want 4: %v", len(visited), visited)
	}

	// depth-first, leaves before parents: gc1 must come before c1, and both children before root (last).
	positions := make(map[string]int)
	for i, name := range visited {
		positions[name] = i
	}

	if positions["gc1"] >= positions["c1"] {
		t.Errorf("expected gc1 to be visited before c1, order was %v", visited)
	}
	if positions["root"] != len(visited)-1 {
		t.Errorf("expected root to be visited last, order was %v", visited)
	}
}

func TestWalkTreePredicateStopsTraversal(t *testing.T) {
	root := newTestNode("root")
	c1 := newTestNode("c1")
	gc1 := newTestNode("gc1")

	root.AddChild(c1)
	c1.AddChild(gc1)

	var visited []string
	WalkTree(root, func(n testTree) {
		visited = append(visited, n.Name())
	}, func(n testTree) bool {
		return n.Name() != "c1" // stop traversal at c1, skipping its subtree
	})

	for _, name := range visited {
		if name == "c1" || name == "gc1" {
			t.Errorf("expected c1 and its children to be skipped, but visited %v", visited)
		}
	}
}

func TestWalkSubTreesExcludesRoot(t *testing.T) {
	root := newTestNode("root")
	c1 := newTestNode("c1")
	root.AddChild(c1)

	var visited []string
	WalkSubTrees(root, func(n testTree) {
		visited = append(visited, n.Name())
	})

	for _, name := range visited {
		if name == "root" {
			t.Error("WalkSubTrees should not visit the root node")
		}
	}
	if len(visited) != 1 || visited[0] != "c1" {
		t.Errorf("WalkSubTrees visited = %v, want [c1]", visited)
	}
}
