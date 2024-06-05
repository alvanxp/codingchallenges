package huffman

import "container/heap"

// Node represents a node in the Huffman tree.
type Node struct {
	Freq  int  // Frequency of the symbol
	Ch    rune // The symbol itself (can be extended for wider characters)
	Left  *Node
	Right *Node
}

// ByFreq implements the heap.Interface for building the priority queue.
type ByFreq []*Node

func (bf ByFreq) Len() int            { return len(bf) }
func (bf ByFreq) Less(i, j int) bool  { return bf[i].Freq < bf[j].Freq }
func (bf ByFreq) Swap(i, j int)       { bf[i], bf[j] = bf[j], bf[i] }
func (bf *ByFreq) Push(x interface{}) { *bf = append(*bf, x.(*Node)) }
func (bf *ByFreq) Pop() interface{} {
	old := *bf
	n := len(old)
	x := old[n-1]
	*bf = old[0 : n-1]
	return x
}

// BuildTree constructs a Huffman tree from symbol frequencies.
func BuildTree(freqs map[rune]int) *Node {
	h := &ByFreq{}
	heap.Init(h)

	// Create leaf nodes and push them into the priority queue.
	for ch, freq := range freqs {
		heap.Push(h, &Node{Freq: freq, Ch: ch})
	}

	// Build the tree by merging nodes.
	for h.Len() > 1 {
		left := heap.Pop(h).(*Node)
		right := heap.Pop(h).(*Node)
		heap.Push(h, &Node{Freq: left.Freq + right.Freq, Left: left, Right: right})
	}

	return heap.Pop(h).(*Node)
}

// Traverse generates Huffman codes from the tree.
func Traverse(node *Node, prefix string, codes map[rune]string) {
	if node == nil {
		return
	}
	if node.Ch != 0 { // Leaf node
		codes[node.Ch] = prefix
	} else { // Internal node
		Traverse(node.Left, prefix+"0", codes)
		Traverse(node.Right, prefix+"1", codes)
	}
}
