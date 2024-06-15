package main

import (
	"bufio"
	"compression/counter"
	"compression/huffman"
	"errors"
	"flag"
	"fmt"
	"os"
)

func main() {
	filePath := getFilePath()
	// if filePath == "" {
	// 	reader, err := getReaderToCompress(filePath)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	c := count(reader)
	// 	root := huffman.BuildTree(c.Counter)
	// 	codes := make(map[rune]string)
	// 	huffman.Traverse(root, "", codes)
	// }
	fmt.Println(filePath)
	var outputFileName string
	flag.StringVar(&outputFileName, "o", "example.txt", "output file name")

	var fileToDecompress string
	flag.StringVar(&fileToDecompress, "d", "example.txt", "input file name")
	flag.Parse()

	fmt.Println(outputFileName)
	Process(CompressParams{FilePath: filePath, OutputPath: outputFileName, Operation: Zip})

	// fmt.Println("Huffman Codes:")
	// for ch, code := range codes {
	// 	fmt.Printf("%c: %s\n", ch, code)
	// }
}

type CompressParams struct {
	FilePath   string
	OutputPath string
	Operation  OperationType
}

type OperationType int

const (
	Zip OperationType = iota
	Unzip
)

func Process(compressParams CompressParams) error {
	switch compressParams.Operation {
	case Zip:
		reader, err := getReaderToCompress(compressParams.FilePath)
		if err != nil {
			fmt.Println("Error: ", err)
			return err
		}
		c := count(reader)
		root := huffman.BuildTree(c.Counter)
		codes := make(map[rune]string)
		huffman.Traverse(root, "", codes)
		writeToFile(compressParams.OutputPath, c, codes)
	case Unzip:
		panic("Not implemented")
	}
	return nil
}

func getReaderToCompress(filePath string) (*bufio.Reader, error) {
	if filePath == "" {
		return nil, errors.New("file path is empty")
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return bufio.NewReader(f), nil
}

func writeToFile(fileName string, c counter.Counter, codes map[rune]string) {
	fmt.Println(fileName)
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	// writeHeader(w, c)
	w.WriteString("some text")
	// writeCodes(w, codes)
	w.Flush()
}

func writeCodes(w *bufio.Writer, codes map[rune]string) {
	for ch, code := range codes {
		w.WriteRune(ch)
		w.WriteRune(' ')
		w.WriteString(code)
		w.WriteRune('\n')
	}
}

func writeHeader(w *bufio.Writer, c counter.Counter) {
	for ch, freq := range c.Counter {
		w.WriteRune(ch)
		w.WriteRune(' ')
		w.WriteRune(rune(freq))
		w.WriteRune('\n')
	}
	w.WriteRune('\n')
}

func count(r *bufio.Reader) counter.Counter {
	c := counter.NewCounter()
	for {
		rc, _, err := r.ReadRune()
		if err != nil {
			break
		}
		if rc != ' ' {
			c.Counter[rc]++
		}
	}
	return c
}

func getFilePath() string {
	args := os.Args[1:]
	if len(args) > 0 {
		path := args[len(args)-1]
		fmt.Println(path)
		return path
	} else {
		fmt.Println("No path provided")
	}
	return ""
}
