package main

import (
	"bufio"
	"compression/counter"
	"compression/huffman"
	"encoding/binary"
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
		c := count(compressParams)
		root := huffman.BuildTree(c.Counter)
		codes := make(map[rune]string)
		huffman.Traverse(root, "", codes)
		//print codes
		for ch, code := range codes {
			fmt.Printf("%c: %s\n", ch, code)
		}

		writeToFile(compressParams, codes, c)
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
	headerLength := ReadNextInt(r)

	headerBuffer := make([]byte, headerLength)
	r.Read(headerBuffer)

	totalChars := ReadNextInt(r)

	header := loadHeader(headerBuffer)
	if header == nil {
		fmt.Println("Error loading header")
		return
	}

	charCouting := 0
	for {
		b, err := r.ReadByte()
		if err != nil {
			break
		}
		if charCouting == totalChars {
			break
		}

		for i := 0; i < 8; i++ {
			bit := (b >> uint(7-i)) & 1
			bits = bits + fmt.Sprintf("%d", bit)
			if ch, ok := header[bits]; ok {
				outputTxtBuilder.WriteRune(ch)
				bits = ""
				charCouting++
			}
		}
	}
	outputFileName := "output2.txt"
	outputTxt := outputTxtBuilder.String()
	e := os.WriteFile(outputFileName, []byte(outputTxt), fs.FileMode(0644))

	if e != nil {
		fmt.Println("Error writing to file: ", e)
		panic(e)
	}
	fmt.Println("Output written to file: ", outputFileName)
}

func ReadNextInt(r *bufio.Reader) int {
	result := 0
	headerLengthBuffer := make([]byte, 4)
	r.Read(headerLengthBuffer)
	result = int(binary.LittleEndian.Uint32(headerLengthBuffer))
	return result
}

func loadHeader(headerBuffer []byte) map[string]rune {
	header := ""
	if len(headerBuffer) > 0 {
		header = string(headerBuffer)
	} else {
		panic("Header is empty")
	}

	headerLines := strings.Split(header, "\n")
	codes := make(map[string]rune)
	for _, l := range headerLines {
		if l == "" {
			continue
		}
		v := strings.Split(l, ":")
		if len(v) != 2 {
			fmt.Println("Invalid header line:", l)
			continue
		}
		code := strings.TrimSpace(v[0])
		ch := strings.TrimSpace(v[1])
		inputNumber, err := strconv.Atoi(ch)
		if err != nil {
			fmt.Println("Invalid character code:", ch)
			continue
		}
		codes[code] = rune(inputNumber)
	}
	return codes
}

func writeToFile(compressParams CompressParams, codes map[rune]string, counter counter.Counter) {
	f, err := os.Create(compressParams.OutputPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	if err != nil {
		panic(err)
	}
	writeCodes(w, codes)
	chartCountBuffer := ConvertIntToBytes(counter.CharCount)
	w.Write(chartCountBuffer)
	writeContent(w, compressParams, codes)
	w.Flush()
}

func ConvertIntToBytes(input int) []byte {
	sb := make([]byte, 4)
	binary.LittleEndian.PutUint32(sb, uint32(input))
	return sb
}

func writeContent(w *bufio.Writer, compressParams CompressParams, codes map[rune]string) {

	f, err := os.Open(compressParams.FilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)

	byteBuffer := byte(0)
	bitIndex := 0
	for {
		rc, _, err := r.ReadRune()
		if err != nil {
			break
		}
		// fmt.Println(rc)
		for _, c := range codes[rc] {

			if c == '1' {
				//add a bit to the left of the byte
				byteBuffer = byteBuffer | 1<<uint(7-bitIndex)
			} else {
				//add a bit to the right of the byte
				byteBuffer = byteBuffer | 0<<uint(7-bitIndex)
			}
			bitIndex = bitIndex + 1
			if bitIndex == 8 {
				w.WriteByte(byteBuffer)
				byteBuffer = 0
				bitIndex = 0
			}
		}
	}
	w.WriteByte(byteBuffer)
	byteBuffer = 0
	bitIndex = 0
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
	count := sb.Len()
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(count))
	w.Write(bs)
	w.WriteString(sb.String())
}

func count(compressParams CompressParams) counter.Counter {
	f, err := os.Open(compressParams.FilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	c := counter.NewCounter()
	charCounter := 0
	for {
		rc, _, err := r.ReadRune()
		if err != nil {
			break
		}
		c.Counter[rc]++
		charCounter++
	}

	c.CharCount = charCounter
	return c
}
