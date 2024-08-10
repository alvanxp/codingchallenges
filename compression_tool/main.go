package main

import (
	"bufio"
	"compression/counter"
	"compression/huffman"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	var inputFileName string
	flag.StringVar(&inputFileName, "zip", "", "file to compress")
	var outputFileName string
	flag.StringVar(&outputFileName, "o", "", "output file of the operation")
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

	err := Process(CompressParams{FilePath: filePath, OutputPath: outputFileName, Operation: operation})
	if err != nil {
		fmt.Println(err)
	}
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
	var err error
	switch compressParams.Operation {
	case Zip:
		err = compress(compressParams)
	case Unzip:
		err = decompress(compressParams.FilePath, compressParams.OutputPath)
	}
	return err
}

func compress(compressParams CompressParams) error {
	c, err := count(compressParams)
	if err != nil {
		return err
	}
	root := huffman.BuildTree(c.Counter)
	codes := make(map[rune]string)
	huffman.Traverse(root, "", codes)
	err = writeToFile(compressParams, codes, c)
	if err != nil {
		return err
	}
	return nil
}

func decompress(filePath string, outputFileName string) error {
	bits := ""
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)

	// Open the output file for writing
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	headerLength := readNextInt(r)

	headerBuffer := make([]byte, headerLength)
	r.Read(headerBuffer)

	totalChars := readNextInt(r)

	header, err := loadHeader(headerBuffer)
	if err != nil {
		return err
	}

	// Buffer to accumulate the decompressed data before writing it to the file
	chunkSize := 1024
	var chunkBuilder strings.Builder

	charCounting := 0
	for {
		b, err := r.ReadByte()
		if err != nil {
			break
		}
		if charCounting == totalChars {
			break
		}

		for i := 0; i < 8; i++ {
			bit := (b >> uint(7-i)) & 1
			bits += fmt.Sprintf("%d", bit)
			if ch, ok := header[bits]; ok {
				chunkBuilder.WriteRune(ch)
				bits = ""
				charCounting++

				// Write to the file when the chunk reaches the specified size
				if chunkBuilder.Len() >= chunkSize {
					_, err := outputFile.WriteString(chunkBuilder.String())
					if err != nil {
						return err
					}
					chunkBuilder.Reset()
				}
			}
		}
	}

	// Write any remaining data in the chunk buffer to the file
	if chunkBuilder.Len() > 0 {
		_, err := outputFile.WriteString(chunkBuilder.String())
		if err != nil {
			return err
		}
	}

	fmt.Println("Decompression complete. Output written to", outputFileName)
	return nil
}

func readNextInt(r *bufio.Reader) int {
	result := 0
	headerLengthBuffer := make([]byte, 4)
	r.Read(headerLengthBuffer)
	result = int(binary.LittleEndian.Uint32(headerLengthBuffer))
	return result
}

func loadHeader(headerBuffer []byte) (map[string]rune, error) {
	headerContent := string(headerBuffer)
	headerLines := strings.Split(headerContent, "|")
	codes := make(map[string]rune)
	for _, line := range headerLines {
		if len(line) == 0 {
			continue
		}
		lineValues := strings.Split(line, ":")
		if len(lineValues) != 2 {
			return nil, errors.New("invalid header line")
		}
		code := string(lineValues[0])
		ch := lineValues[1]
		i, err := strconv.Atoi(ch)
		if err != nil {
			return nil, err
		}
		codes[code] = rune(i)
	}
	return codes, nil
}

func writeToFile(compressParams CompressParams, codes map[rune]string, counter counter.Counter) error {
	f, err := os.Create(compressParams.OutputPath)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	err = writeCodes(w, codes)
	if err != nil {
		return err
	}
	chartCountBuffer := convertIntToBytes(counter.TotalChars)
	_, err = w.Write(chartCountBuffer)
	if err != nil {
		return err
	}
	err = writeContent(w, compressParams, codes)
	if err != nil {
		return err
	}
	err = w.Flush()
	if err != nil {
		return err
	}
	return nil
}

func convertIntToBytes(input int) []byte {
	sb := make([]byte, 4)
	binary.LittleEndian.PutUint32(sb, uint32(input))
	return sb
}

func writeContent(w *bufio.Writer, compressParams CompressParams, codes map[rune]string) error {

	f, err := os.Open(compressParams.FilePath)
	if err != nil {
		return err
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
	return nil
}

func writeCodes(w *bufio.Writer, codes map[rune]string) error {
	sb := strings.Builder{}
	for ch, code := range codes {
		sb.WriteString(code)
		sb.WriteRune(':')
		sb.WriteString(fmt.Sprintf("%v", ch))
		sb.WriteRune('|')
	}
	count := sb.Len()
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(count))
	_, err := w.Write(bs)
	if err != nil {
		return err
	}
	_, err = w.WriteString(sb.String())
	if err != nil {
		return err
	}
	return nil
}

func count(compressParams CompressParams) (counter.Counter, error) {
	f, err := os.Open(compressParams.FilePath)
	if err != nil {
		return counter.Counter{}, err
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

	c.TotalChars = charCounter
	return c, nil
}
