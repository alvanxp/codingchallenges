package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {

	filename, err := getFilePath()
	if err != nil {
		fmt.Println("Error processing file:", err)
		return
	}
	showByteCount := flag.Bool("c", false, "show the number of bytes in the file")
	showLinesCount := flag.Bool("l", false, "show the number of lines in the file")
	showWordsCount := flag.Bool("w", false, "show the number of words in the file")
	flag.Parse()

	var byteCount, linesCount, wordsCount int = 0, 0, 0

	var consoleOutput string
	data, err := getDataFromFile(filename)
	if err != nil {
		fmt.Println("Error processing file:", err)
		return
	}
	if *showByteCount {
		byteCount = data.GetBytesCount()
		consoleOutput = AppendCounterFormat(byteCount, consoleOutput)
	}

	if *showLinesCount {
		linesCount = data.GetLinesCount()
		consoleOutput = AppendCounterFormat(linesCount, consoleOutput)
	}

	if *showWordsCount {
		wordsCount = data.GetWordsCount()
		consoleOutput = AppendCounterFormat(wordsCount, consoleOutput)
	}

	consoleOutput = fmt.Sprintf("%s %s", consoleOutput, filename)

	fmt.Println(consoleOutput)
}

func AppendCounterFormat(count int, input string) string {
	return fmt.Sprintf("%s %d", input, count)
}

func getDataFromFile(filename string) (fileData, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return data, err
	}

	return data, nil
}

func getFilePath() (string, error) {
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i][0] != '-' {
			return os.Args[i], nil
		}
	}
	return "", fmt.Errorf("No file path provided")
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
