package data

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type DataCounter struct {
	filename     string
	data         []byte
	loadFromFile bool
}

func ProcessFileData(filename string, loadFromFile bool) *DataCounter {
	dataCounter := &DataCounter{filename: filename, loadFromFile: loadFromFile}
	dataCounter.loadData()
	return dataCounter
}

func (fp *DataCounter) loadData() {
	if fp.loadFromFile {
		d, err := os.ReadFile(fp.filename)
		if err != nil {
			panic("error loading file")
		}

		fp.data = d
		return
	}
	reader := bufio.NewReader(os.Stdin)
	d, err := io.ReadAll(io.Reader(reader))
	if err != nil {
		panic("error loading data")
	}
	fp.data = d
}

func (fp DataCounter) GetBytesCount() int {
	return len(fp.data)
}

func (fp DataCounter) GetCharsCount() int {
	content := string(fp.data)
	chars := strings.Split(content, "")
	return len(chars)
}

func (fp DataCounter) GetLinesCount() int {
	content := string(fp.data)
	lines := strings.Split(content, "\n")
	return len(lines)
}

func (fp DataCounter) GetWordsCount() int {
	content := string(fp.data)
	words := strings.Split(content, " ")
	return len(words)
}
