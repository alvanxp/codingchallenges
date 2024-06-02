package main

type HuffBaseNode interface {
	IsLeaf() bool
	Weight() int
}
type HuffmanTree struct {
	Root HuffBaseNode
}

func (t HuffmanTree) Weight() int {
	return t.Root.Weight()
}

func NewHuffTree(el rune, wt int) *HuffmanTree {
	return &HuffmanTree{Root: &HuffLeafNode{weight: wt, element: el}}
}

func NewHuffInternalNode(left, right *HuffmanInternalNode, wt int) *HuffmanTree {
	return &HuffmanTree{
		Root: &HuffmanInternalNode{
			Left:   left,
			Right:  right,
			weight: wt,
		},
	}
}

func BuildTree() HuffmanTree {
	var tmp1, tmp2, tmp3 HuffmanTree

	for {
		tmp1 = tmp3
		tmp2 = tmp3
		if tmp1.Weight() < tmp2.Weight() {
			tmp3 = NewHuffInternalNode(tmp1.Root.(*HuffmanInternalNode), tmp2.Root.(*HuffmanInternalNode), tmp1.Weight()+tmp2.Weight())
		} else {
			tmp3 = NewHuffInternalNode(tmp2.Root.(*HuffmanInternalNode), tmp1.Root.(*HuffmanInternalNode), tmp1.Weight()+tmp2.Weight())
		}
		if tmp3.Weight() == 0 {
			break
		}
	}
	return tmp3

}
