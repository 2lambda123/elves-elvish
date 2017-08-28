// Package vector implements persistent vector.
package vector

import "github.com/xiaq/persistent/types"

const (
	chunkBits  = 5
	nodeSize   = 1 << chunkBits
	tailMaxLen = nodeSize
	chunkMask  = nodeSize - 1
)

// Vector is a persistent sequential container for arbitrary values. It supports
// O(1) lookup by index, modification by index, and insertion and removal
// operations at the end. Being a persistent variant of the data structure, it
// is immutable, and provides O(1) operations to create modified versions of the
// vector that shares the underlying data structure, making it suitable for
// concurrent access. The empty value is a valid empty vector.
type Vector interface {
	types.Equaler
	// Len returns the length of the vector.
	Len() int
	// Nth returns the i-th element of the vector. It returns nil if the index
	// is smaller than 0 or greater than or equal to the length of the vector.
	Nth(i int) interface{}
	// AssocN returns an almost identical Vector, with the i-th element
	// replaced. If the index is smaller than 0 or greater than the length of
	// the vector, it returns nil. If the index is equal to the size of the
	// vector, it is equivalent to Cons.
	AssocN(i int, val interface{}) Vector
	// Cons returns an almost identical Vector, with an additional element
	// appended to the end.
	Cons(val interface{}) Vector
	// Pop returns an almost identical Vector, with the last element removed. It
	// returns nil if the vector is already empty.
	Pop() Vector
	// SubVector returns a subvector containing the elements from i up to but
	// not including j.
	SubVector(i, j int) Vector
	// Iterator returns an iterator over the vector.
	Iterator() Iterator
}

// Equal determines whether two Vector are structurally equal.
func Equal(v1, v2 Vector) bool {
	if v1.Len() != v2.Len() {
		return false
	}
	n := v1.Len()
	it1 := v1.Iterator()
	it2 := v2.Iterator()
	for i := 0; i < n; i++ {
		elem1 := it1.Elem()
		elem2 := it2.Elem()
		if elem1eq, ok := elem1.(types.Equaler); ok {
			if !elem1eq.Equal(elem2) {
				return false
			}
		} else {
			if elem1 != elem2 {
				return false
			}
		}
		it1.Next()
		it2.Next()
	}
	return true
}

func equal(v Vector, other interface{}) bool {
	v2, ok := other.(Vector)
	if !ok {
		return false
	}
	return Equal(v, v2)
}

// Iterator is an iterator over vector elements. It can be used like this:
//
//     for it := v.Iterator(); it.HasElem(); it.Next() {
//         elem := it.Elem()
//         // do something with elem...
//     }
type Iterator interface {
	// Elem returns the element at the current position.
	Elem() interface{}
	// HasElem returns whether the iterator is pointing to an element.
	HasElem() bool
	// Next moves the iterator to the next position.
	Next()
}

type vector struct {
	count int
	// height of the tree structure, defined to be 0 when root is a leaf.
	height uint
	root   node
	tail   []interface{}
}

// Empty is an empty Vector.
var Empty Vector = &vector{}

// node is a node in the vector tree. It is always of the size nodeSize.
type node []interface{}

func newNode() node {
	return node(make([]interface{}, nodeSize))
}

func (n node) clone() node {
	m := newNode()
	copy(m, n)
	return m
}

func (v *vector) Equal(other interface{}) bool {
	return equal(v, other)
}

// Count returns the number of elements in a Vector.
func (v *vector) Len() int {
	return v.count
}

// treeSize returns the number of elements stored in the tree (as opposed to the
// tail).
func (v *vector) treeSize() int {
	if v.count < tailMaxLen {
		return 0
	}
	return ((v.count - 1) >> chunkBits) << chunkBits
}

// sliceFor returns the slice where the i-th element is stored. The index must
// be in bound.
func (v *vector) sliceFor(i int) []interface{} {
	if i >= v.treeSize() {
		return v.tail
	}
	n := v.root
	for shift := v.height * chunkBits; shift > 0; shift -= chunkBits {
		n = n[(i>>shift)&chunkMask].(node)
	}
	return n
}

func (v *vector) Nth(i int) interface{} {
	if i < 0 || i >= v.count {
		return nil
	}
	return v.sliceFor(i)[i&chunkMask]
}

func (v *vector) AssocN(i int, val interface{}) Vector {
	if i < 0 || i > v.count {
		return nil
	} else if i == v.count {
		return v.Cons(val)
	}
	if i >= v.treeSize() {
		newTail := append([]interface{}(nil), v.tail...)
		copy(newTail, v.tail)
		newTail[i&chunkMask] = val
		return &vector{v.count, v.height, v.root, newTail}
	}
	return &vector{v.count, v.height, doAssoc(v.height, v.root, i, val), v.tail}
}

// doAssoc returns an almost identical tree, with the i-th element replaced by
// val.
func doAssoc(height uint, n node, i int, val interface{}) node {
	m := n.clone()
	if height == 0 {
		m[i&chunkMask] = val
	} else {
		sub := (i >> (height * chunkBits)) & chunkMask
		m[sub] = doAssoc(height-1, m[sub].(node), i, val)
	}
	return m
}

func (v *vector) Cons(val interface{}) Vector {
	// Room in tail?
	if v.count-v.treeSize() < tailMaxLen {
		newTail := make([]interface{}, len(v.tail)+1)
		copy(newTail, v.tail)
		newTail[len(v.tail)] = val
		return &vector{v.count + 1, v.height, v.root, newTail}
	}
	// Full tail; push into tree.
	tailNode := node(v.tail)
	newHeight := v.height
	var newRoot node
	// Overflow root?
	if (v.count >> chunkBits) > (1 << (v.height * chunkBits)) {
		newRoot = newNode()
		newRoot[0] = v.root
		newRoot[1] = newPath(v.height, tailNode)
		newHeight++
	} else {
		newRoot = v.pushTail(v.height, v.root, tailNode)
	}
	return &vector{v.count + 1, newHeight, newRoot, []interface{}{val}}
}

// pushTail returns a tree with tail appended.
func (v *vector) pushTail(height uint, n node, tail node) node {
	if height == 0 {
		return tail
	}
	idx := ((v.count - 1) >> (height * chunkBits)) & chunkMask
	m := n.clone()
	child := n[idx]
	if child == nil {
		m[idx] = newPath(height-1, tail)
	} else {
		m[idx] = v.pushTail(height-1, child.(node), tail)
	}
	return m
}

// newPath creates a left-branching tree of specified height and leaf.
func newPath(height uint, leaf node) node {
	if height == 0 {
		return leaf
	}
	ret := newNode()
	ret[0] = newPath(height-1, leaf)
	return ret
}

func (v *vector) Pop() Vector {
	switch v.count {
	case 0:
		return nil
	case 1:
		return Empty
	}
	if v.count-v.treeSize() > 1 {
		newTail := make([]interface{}, len(v.tail)-1)
		copy(newTail, v.tail)
		return &vector{v.count - 1, v.height, v.root, newTail}
	}
	newTail := v.sliceFor(v.count - 2)
	newRoot := v.popTail(v.height, v.root)
	newHeight := v.height
	if v.height > 0 && newRoot[1] == nil {
		newRoot = newRoot[0].(node)
		newHeight--
	}
	return &vector{v.count - 1, newHeight, newRoot, newTail}
}

// popTail returns a new tree with the last leaf removed.
func (v *vector) popTail(level uint, n node) node {
	idx := ((v.count - 2) >> (level * chunkBits)) & chunkMask
	if level > 1 {
		newChild := v.popTail(level-1, n[idx].(node))
		if newChild == nil && idx == 0 {
			return nil
		}
		m := n.clone()
		if newChild == nil {
			// This is needed since `m[idx] = newChild` would store an
			// interface{} with a non-nil type part, which is non-nil
			m[idx] = nil
		} else {
			m[idx] = newChild
		}
		return m
	} else if idx == 0 {
		return nil
	} else {
		m := n.clone()
		m[idx] = nil
		return m
	}
}

func (v *vector) SubVector(begin, end int) Vector {
	if begin < 0 || begin > end || end > v.count {
		return nil
	}
	return &subVector{v, begin, end}
}

type subVector struct {
	v     *vector
	begin int
	end   int
}

func (s *subVector) Equal(other interface{}) bool {
	return equal(s, other)
}

func (s *subVector) Len() int {
	return s.end - s.begin
}

func (s *subVector) Nth(i int) interface{} {
	if i < 0 || s.begin+i >= s.end {
		return nil
	}
	return s.v.Nth(s.begin + i)
}

func (s *subVector) AssocN(i int, val interface{}) Vector {
	if i < 0 || s.begin+i > s.end {
		return nil
	} else if s.begin+i == s.end {
		return s.Cons(val)
	}
	return s.v.AssocN(s.begin+i, val).SubVector(s.begin, s.end)
}

func (s *subVector) Cons(val interface{}) Vector {
	return s.v.AssocN(s.end, val).SubVector(s.begin, s.end+1)
}

func (s *subVector) Pop() Vector {
	switch s.Len() {
	case 0:
		return nil
	case 1:
		return Empty
	default:
		return s.v.SubVector(s.begin, s.end-1)
	}
}

func (s *subVector) SubVector(i, j int) Vector {
	return s.v.SubVector(s.begin+i, s.begin+j)
}

func (s *subVector) Iterator() Iterator {
	return newIteratorWithRange(s.v, s.begin, s.end)
}

func (v *vector) Iterator() Iterator {
	return newIterator(v)
}

type iterator struct {
	v        *vector
	treeSize int
	index    int
	end      int
	path     []pathEntry
}

type pathEntry struct {
	node  node
	index int
}

func (e pathEntry) current() interface{} {
	return e.node[e.index]
}

func newIterator(v *vector) *iterator {
	return newIteratorWithRange(v, 0, v.Len())
}

func newIteratorWithRange(v *vector, begin, end int) *iterator {
	it := &iterator{v, v.treeSize(), begin, end, nil}
	// Find the node for begin, remembering all nodes along the path.
	n := v.root
	for shift := v.height * chunkBits; shift > 0; shift -= chunkBits {
		idx := (begin >> shift) & chunkMask
		it.path = append(it.path, pathEntry{n, idx})
		n = n[idx].(node)
	}
	it.path = append(it.path, pathEntry{n, begin & chunkMask})
	return it
}

func (it *iterator) Elem() interface{} {
	if it.index >= it.treeSize {
		return it.v.tail[it.index-it.treeSize]
	}
	return it.path[len(it.path)-1].current()
}

func (it *iterator) HasElem() bool {
	return it.index < it.end
}

func (it *iterator) Next() {
	if it.index+1 >= it.treeSize {
		// Next element is in tail. Just increment the index.
		it.index++
		return
	}
	// Find the deepest level that can be advanced.
	var i int
	for i = len(it.path) - 1; i >= 0; i-- {
		e := it.path[i]
		if e.index+1 < len(e.node) && e.node[e.index+1] != nil {
			break
		}
	}
	if i == -1 {
		panic("cannot advance; vector iterator bug")
	}
	// Advance on this node, and re-populate all deeper levels.
	it.path[i].index++
	for i++; i < len(it.path); i++ {
		it.path[i] = pathEntry{it.path[i-1].current().(node), 0}
	}
	it.index++
}
