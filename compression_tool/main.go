package main

import (
	"bufio"
	"compression/counter"
	"compression/huffman"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	var inputFileName string
	flag.StringVar(&inputFileName, "zip", "", "file to compress")
	var outputFileName string
	flag.StringVar(&outputFileName, "o", "", "output of the compressed file")
	var fileToDecompress string
	flag.StringVar(&fileToDecompress, "unzip", "example.txt", "input file name")
	flag.Parse()
	var filePath string
	var operation OperationType
	if inputFileName == "" {
		fmt.Println("Unzipping file: ", fileToDecompress)
		filePath = fileToDecompress
		operation = Unzip
	} else {
		fmt.Println("Zipping file: ", inputFileName)
		filePath = inputFileName
		operation = Zip
	}

	fmt.Println("Operation: ", operation)
	Process(CompressParams{FilePath: filePath, OutputPath: outputFileName, Operation: operation})
}

// CompressParams represents the parameters for compression.
type CompressParams struct {
	FilePath   string
	OutputPath string
	Operation  OperationType
}

// OperationType represents the type for compression operations.
type OperationType int

const (
	// Zip represents the operation type for compression.
	Zip OperationType = iota
	// Unzip represents the operation type for decompression.
	Unzip
)

// Process compresses or decompresses a file based on the provided parameters.
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
		reader, err = getReaderToCompress(compressParams.FilePath)
		if err != nil {
			fmt.Println("Error: ", err)
			return err
		}
		writeToFile(compressParams.OutputPath, c, codes, reader)
	case Unzip:
		decompress(compressParams.FilePath)
	}
	return nil
}

func decompress(filePath string) {
	header := loadHeader(filePath)
	bits := ""
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(f)
	for {
		b, err := r.ReadByte()
		if err != nil {
			break
		}
		for i := 0; i < 8; i++ {
			bit := (b >> uint(7-i)) & 1
			bits = bits + fmt.Sprintf("%d", bit)
			if ch, ok := header[bits]; ok {
				if ch != 239 {
					fmt.Print(string(ch))
				}
				bits = ""
			}
		}
	}
}

func loadHeader(filePath string) map[string]rune {
	hf := filePath + ".header"
	f, err := os.Open(hf)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(f)
	codes := make(map[string]rune)

	for {
		line, _, err := r.ReadLine()
		if err != nil {
			break
		}
		l := strings.Trim(string(line), " ")
		fmt.Println(l)
		v := strings.Split(l, ":")
		if len(v) > 1 {
			code := v[0]
			ch := '\n'
			codes[code] = ch
			continue
		}
		if len(v) > 2 {
			ch := ' '
			code := string(v[2])
			codes[code] = ch
			continue
		}

		ch := rune(v[0][0])
		code := string(v[1])
		codes[code] = ch
	}
	return codes
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

func writeToFile(fileName string, c counter.Counter, codes map[rune]string, r *bufio.Reader) {
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	hf, err := os.Create(fileName + ".header")
	if err != nil {
		panic(err)
	}
	defer hf.Close()
	hw := bufio.NewWriter(hf)
	writeCodes(hw, codes)
	hw.Flush()
	writeContent(w, r, codes)
	w.Flush()
}

func writeContent(w *bufio.Writer, r *bufio.Reader, codes map[rune]string) {

	temp := byte(0)
	i := 0
	for {
		rc, _, err := r.ReadRune()
		if err != nil {
			break
		}
		// fmt.Println(rc)
		for _, c := range codes[rc] {

			if c == '1' {
				//add a bit to the left of the byte
				temp = temp | 1<<uint(7-i)
			} else {
				//add a bit to the right of the byte
				temp = temp | 0<<uint(7-i)
			}
			i = i + 1
			if i == 8 {
				w.WriteByte(temp)
				temp = 0
				i = 0
			}
		}
	}
	w.WriteByte(temp)
	temp = 0
	i = 0
}

func writeCodes(w *bufio.Writer, codes map[rune]string) {
	for ch, code := range codes {
		if ch != '\n' {
			w.WriteString(code)
			w.WriteRune(':')
			w.WriteRune(ch)
			w.WriteRune('\n')
		}
	}
	code := codes['\n']
	w.WriteString(code)
	w.WriteRune('\n')
}

// func writeHeader(w *bufio.Writer, c counter.Counter) {
// 	for ch, freq := range c.Counter {
// 		w.WriteRune(ch)
// 		w.WriteRune(' ')
// 		w.WriteRune(rune(freq))
// 		w.WriteRune('\n')
// 	}
// 	w.WriteRune('\n')
// }

func count(r *bufio.Reader) counter.Counter {
	c := counter.NewCounter()
	for {
		rc, _, err := r.ReadRune()
		if err != nil {
			break
		}
		// if rc != ' ' {
		c.Counter[rc]++
		// }
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
