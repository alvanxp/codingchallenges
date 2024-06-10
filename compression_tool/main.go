package main

import (
	"bufio"
	"compression/counter"
	"compression/huffman"
	"fmt"
	"os"
)

func main() {
	fileName := getFilePath()
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	c := count(reader)
	root := huffman.BuildTree(c.Counter)
	codes := make(map[rune]string)
	huffman.Traverse(root, "", codes)

	fmt.Println("Huffman Codes:")
	for ch, code := range codes {
		fmt.Printf("%c: %s\n", ch, code)
	}
}

func count(r *bufio.Reader) counter.Counter {
	c := counter.NewCounter()
	for {
		rc, _, err := r.ReadRune()
		if err != nil {
			break
		}
		c.Counter[rc]++
	}
	return c
}

func getFilePath() string {
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i][0] != '-' {
			return os.Args[i]
		}
	}
	return ""
}
