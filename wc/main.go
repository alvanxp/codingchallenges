package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {

	filename := getFilePath()
	showByteCount := flag.Bool("c", false, "show the number of bytes in the file")
	showLinesCount := flag.Bool("l", false, "show the number of lines in the file")
	showWordsCount := flag.Bool("w", false, "show the number of words in the file")
	flag.Parse()

	byteCount, linesCount, wordsCount := -1, -1, -1

	var consoleOutput string
	data, err := getDataFromFile(filename)
	if *showByteCount {
		byteCount = data.GetBytesCount()
		if err != nil {
			fmt.Println("Error processing file:", err)
		}
		consoleOutput = AddCounterFormat(byteCount, consoleOutput)
	}

	if *showLinesCount {
		linesCount = data.GetLinesCount()
		if err != nil {
			fmt.Println("Error processing file:", err)
		}
		consoleOutput = AddCounterFormat(linesCount, consoleOutput)
	}

	if *showWordsCount {
		wordsCount = data.GetWordsCount()
		if err != nil {
			fmt.Println("Error processing file:", err)
		}
		consoleOutput = AddCounterFormat(wordsCount, consoleOutput)
	}

	consoleOutput = fmt.Sprintf("%s %s", consoleOutput, filename)

	fmt.Println(consoleOutput)
}

func AddCounterFormat(count int, input string) string {
	return fmt.Sprintf("%s %d", input, count)
}

func getDataFromFile(filename string) (fileData, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return data, err
	}

	return data, nil
}

func getFilePath() string {
	for i := range os.Args {
		if i == 0 {
			continue
		}
		if os.Args[i][0] != '-' {
			return os.Args[i]
		}
	}
	return ""
}

type fileData []byte

func (f fileData) GetBytesCount() int {
	return len(f)
}

func (f fileData) GetLinesCount() int {
	content := string(f)
	lines := strings.Split(content, "\n")
	return len(lines)
}

func (f fileData) GetWordsCount() int {
	content := string(f)
	words := strings.Split(content, " ")
	return len(words)
}
