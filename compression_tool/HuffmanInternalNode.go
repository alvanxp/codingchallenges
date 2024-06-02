package main

type HuffmanInternalNode struct {
	Left   *HuffmanInternalNode
	Right  *HuffmanInternalNode
	weight int
}

func (n HuffmanInternalNode) IsLeaf() bool {
	return false
}

func (n HuffmanInternalNode) Weight() int {
	return n.weight
}