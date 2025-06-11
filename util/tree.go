package util

import (
	"slices"
)

// id_counter is used to provide treenodes with unique IDs. TreeNode[T] need IDs so they can be compared
// TODO: investigate using the self pointer for comparison.
var id_counter int

func gen_ID() int {
	id_counter += 1
	return id_counter
}

// tree describes things with parent/child relationships. current this interface in unexported because
// concrete tree objects are far more useful if they satisfy the generic interface below. maybe someday
// we'll find a reason to export this though?? who knows.
type tree interface {
	getParentNode() tree
	addChildNode(child tree)
	removeChildNode(child tree)
	removeChildByIndex(index int)
	setParentNode(parent tree)

	ChildCount() int
	Deparent()
	getID() int
}

// TreeType[T] is the interface for trees of objects T (duh). Ts can be added as children to other Ts,
// and as such a tree structure is built. Trees ensure all contained items are unique (no object in the
// can have 2 parents).
type TreeType[T tree] interface {
	tree
	GetParent() T
	GetChildren() []T
	GetSelf() T
	AddChild(child T)
	AddChildren(children ...T)
	RemoveChild(child T)
}

// TreeNode[T] implements the TreeType[T] generic interface. This object can stand freeform, or can
// be embedded in other type to give them tree functionality. Remember that the node is not usable
// unless it has been initialized with Init. Failing to do this will crash your whole gig and make the
// computer have a headache. So don't do that.
type TreeNode[T tree] struct {
	parent   T
	self     T
	children []T
	node_id  int
}

// Initializes the node with a reference to the object in the tree that it represents. If this is not
// done the tree cannot work.
func (t *TreeNode[T]) Init(self T) {
	t.self = self
}

func (t TreeNode[T]) GetSelf() T {
	return t.self
}

func (t TreeNode[T]) GetParent() T {
	return t.parent
}

func (t TreeNode[T]) getParentNode() tree {
	return t.parent
}

func (t TreeNode[T]) GetChildren() []T {
	return t.children
}

// Retrieves the node ID. This is a unique ID to each node and is used to make nodes comparable.
// NOTE: this ID is generated at runtime when first needed, so may change between one run of the
// program to the next.
func (t *TreeNode[T]) getID() int {
	if t.node_id == 0 {
		t.node_id = gen_ID()
	}

	return t.node_id
}

func (t *TreeNode[T]) AddChild(node T) {
	t.addChildNode(node)
	node.setParentNode(t.self)
}

func (t *TreeNode[T]) addChildNode(node tree) {
	//check duplicate add. deparent if child is already parented
	if parent := node.getParentNode(); parent != nil {
		if parent.getID() == t.getID() {
			return
		} else {
			node.Deparent()
		}
	}

	if t.children == nil {
		t.children = make([]T, 0)
	}

	t.children = append(t.children, node.(T))

}

func (t *TreeNode[T]) AddChildren(nodes ...T) {
	for i := range nodes {
		t.AddChild(nodes[i])
	}
}

func (t *TreeNode[T]) RemoveChild(node T) {
	t.removeChildNode(node)
}

func (t *TreeNode[T]) removeChildNode(node tree) {
	if t.children == nil {
		return
	}

	for i, child := range t.children {
		if child.getID() == node.getID() {
			t.removeChildByIndex(i)
			break
		}
	}
}

func (t *TreeNode[T]) removeChildByIndex(index int) {
	if index >= len(t.children) {
		return
	}

	t.children[index].Deparent()
	t.children = slices.Delete(t.children, index, index+1)
}

func (t *TreeNode[T]) setParentNode(node tree) {
	t.parent = node.(T)
}

func (t *TreeNode[T]) Deparent() {
	var nil_parent T
	t.parent = nil_parent
}

func (t TreeNode[T]) ChildCount() int {
	return len(t.children)
}

// Walks the tree, starting from the provided root node. The function f is called in a depth-first manner, on leaf
// nodes first and then up. The type T must be specified, inferring it doesn't really seem to work. Optionally takes
// any number of predicates. If any are provided, they are called for each node before continuing to traverse the
// tree and if any return false, stops the walk for that node and all of its children.
// THINK: should there be non-depth-first versions?
func WalkTree[T TreeType[T]](node T, fn func(T), predicates ...func(T) bool) {
	for i := range predicates {
		if !predicates[i](node) {
			return
		}
	}

	for _, child := range node.GetChildren() {
		if predicates != nil {
			WalkTree(child, fn, predicates...)
		} else {
			WalkTree(child, fn)
		}
	}

	fn(node)
}

// Walks the tree, calling function f on all subnodes in the tree below the provided node.
// This is identical to Walktree, except it doesn't call f on the root node.
func WalkSubTrees[T TreeType[T]](node T, fn func(T), predicates ...func(T) bool) {
	for i := range predicates {
		if !predicates[i](node) {
			return
		}
	}

	for _, child := range node.GetChildren() {
		if predicates != nil {
			WalkTree(child, fn, predicates...)
		} else {
			WalkTree(child, fn)
		}
	}
}
