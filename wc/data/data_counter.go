package data

import (
	"bufio"
	"io"
	"log"
	"os"
	"unicode"
)

type DataCounter struct {
	filename     string
	wordsCount   int
	bytesCount   int
	charsCount   int
	linesCount   int
	loadFromFile bool
}

func ProcessCouting(filename string, loadFromFile bool) *DataCounter {
	dataCounter := &DataCounter{filename: filename, loadFromFile: loadFromFile}
	dataCounter.loadData()
	return dataCounter
}

func (fp *DataCounter) loadData() {

	var reader *bufio.Reader
	if fp.loadFromFile {
		file, err := os.Open(fp.filename)
		defer file.Close()
		if err != nil {
			panic("error loading file")
		}
		reader = bufio.NewReader(file)
		fp.count(reader)
		return
	}
	reader = bufio.NewReader(os.Stdin)

	fp.count(reader)
}

func (fp *DataCounter) count(r *bufio.Reader) {
	var previousChar rune
	for {
		c, sz, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}
		fp.bytesCount += sz
		fp.charsCount++
		if c == '\n' {
			fp.linesCount++
		}
		if unicode.IsSpace(c) && !unicode.IsSpace(previousChar) {
			fp.wordsCount++
		}
		previousChar = c
	}

}

func (fp DataCounter) GetBytesCount() int {
	return fp.bytesCount
}

func (fp DataCounter) GetCharsCount() int {
	return fp.charsCount
}

func (fp DataCounter) GetLinesCount() int {
	return fp.linesCount
}

func (fp DataCounter) GetWordsCount() int {
	return fp.wordsCount
}
