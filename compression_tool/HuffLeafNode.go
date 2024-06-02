package main

type HuffLeafNode struct {
	Left    *HuffmanInternalNode
	Right   *HuffmanInternalNode
	weight  int
	element rune
}

func (n HuffLeafNode) Value() rune {
	return n.element
}

func (n HuffLeafNode) IsLeaf() bool {
	return true
}

func (n HuffLeafNode) Weight() int {
	return n.weight
}
