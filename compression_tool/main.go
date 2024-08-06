package main

import (
	"bufio"
	"compression/counter"
	"compression/huffman"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"strconv"
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
			return err
		}
		c := count(reader)
		root := huffman.BuildTree(c.Counter)
		codes := make(map[rune]string)
		huffman.Traverse(root, "", codes)
		//print codes
		for ch, code := range codes {
			fmt.Printf("%c: %s\n", ch, code)
		}
		reader, err = getReaderToCompress(compressParams.FilePath)
		if err != nil {
			return err
		}

		writeToFile(compressParams.OutputPath, codes, c, reader)
	case Unzip:
		decompress(compressParams.FilePath)
	}
	return nil
}

func decompress(filePath string) {
	bits := ""
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	r := bufio.NewReader(f)

	var outputTxtBuilder strings.Builder = strings.Builder{}
	headerLength := 0
	headerLengthBuffer := make([]byte, 4)
	r.Read(headerLengthBuffer)
	headerLength = int(binary.LittleEndian.Uint32(headerLengthBuffer))
	fmt.Println("Header length: ", headerLength)
	headerBuffer := make([]byte, headerLength)
	r.Read(headerBuffer)
	charCounterBuffer:=make([]byte, 4)
	r.Read(charCounterBuffer)
	charCounter:=int(binary.LittleEndian.Uint32(charCounterBuffer))
	fmt.Println("Header Buffer length: ", len(headerBuffer))
	fmt.Println("Header Buffer: ", headerBuffer)
	header := loadHeader(headerBuffer)
	if header == nil {
		fmt.Println("Error loading header")
		return
	}

	c := 0
	for {
		b, err := r.ReadByte()
		if err != nil {
			break
		}
		if c == charCounter {
			break
		}

		for i := 0; i < 8; i++ {
			bit := (b >> uint(7-i)) & 1
			bits = bits + fmt.Sprintf("%d", bit)
			if ch, ok := header[bits]; ok {
				outputTxtBuilder.WriteString(string(ch))
				bits = ""
				c++
			}
		}
	}
	outputFileName := "output2.txt"
	//write to file
	outputTxt := outputTxtBuilder.String()
	e := os.WriteFile(outputFileName, []byte(outputTxt), fs.FileMode(0644))

	if e != nil {
		fmt.Println("Error writing to file: ", e)
		panic(e)
	}
	fmt.Println("Output written to file: ", outputFileName)
}

func loadHeader(headerBuffer []byte) map[string]rune {
	cont := string(headerBuffer)
	fmt.Println("Header: \n", cont)
	sc := strings.Split(cont, "\n")
	codes := make(map[string]rune)
	for _, l := range sc {
		if l == "" {
			continue
		}
		v := strings.Split(l, ":")
		if len(v) == 1 {
			fmt.Println(v[0])
			ch := '\n'
			code := string(v[0])
			codes[code] = ch
			continue
		}
		if len(v[1]) == 0 {
			codes[v[0]] = '\n'
			continue
		}
		code := string(v[0])
		ch := v[1]
		i, err := strconv.Atoi(ch)
		if err != nil {
			panic(err)
		}
		codes[code] = rune(i)
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

	// defer f.Close()
	return bufio.NewReader(f), nil
}

func writeToFile(fileName string, codes map[rune]string, counter counter.Counter, r *bufio.Reader) {
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	if err != nil {
		panic(err)
	}
	writeCodes(w, codes)
	sb:=make([]byte, 4)
	binary.LittleEndian.PutUint32(sb, uint32(counter.CharCount))
	w.Write(sb)
	w.WriteByte(byte(counter.CharCount))
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
	sb := strings.Builder{}
	for ch, code := range codes {
		if ch != '\n' {
			sb.WriteString(code)
			sb.WriteRune(':')
			sb.WriteString(fmt.Sprintf("%v", ch))
			sb.WriteRune('\n')
		}
	}
	code := codes['\n']
	sb.WriteString(code)
	sb.WriteRune('\n')
	x := sb.Len()
	//convert x to bytes and write to file
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(x))

	w.Write(bs)
	fmt.Println("Header length: ", x)
	w.WriteString(sb.String())
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
	charCounter := 0
	for {
		rc, _, err := r.ReadRune()
		if err != nil {
			break
		}
		// if rc != ' ' {
		c.Counter[rc]++
		charCounter++
		// }
	}

	c.CharCount = charCounter
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
