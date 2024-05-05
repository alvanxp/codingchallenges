package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	filename := flag.String("f", "", "file to process")
	showByteCount := flag.Bool("c", false, "show the number of bytes in the file")
	showLinesCount := flag.Bool("l", false, "show the number of lines in the file")
	showWordsCount := flag.Bool("w", false, "show the number of words in the file")
	flag.Parse()

	byteCount, linesCount, wordsCount := -1, -1, -1

	var consoleOutput string
	data, err := getDataFromFile(*filename)
	if *showByteCount {
		byteCount = getBytesCount(data)
		if err != nil {
			fmt.Println("Error processing file:", err)
		}
		consoleOutput = getstring(byteCount, consoleOutput)
	}

	if *showLinesCount {
		linesCount = getLinesCount(data)
		if err != nil {
			fmt.Println("Error processing file:", err)
		}
		consoleOutput = getstring(linesCount, consoleOutput)
	}

	if *showWordsCount {
		wordsCount = getWordsCount(data)
		if err != nil {
			fmt.Println("Error processing file:", err)
		}
		consoleOutput = getstring(wordsCount, consoleOutput)
	}

	consoleOutput = fmt.Sprintf("%s %s", consoleOutput, *filename)

	fmt.Println(consoleOutput)
}

func getstring(count int, a string) string {
	return fmt.Sprintf("%s %d", a, count)
}

func getDataFromFile(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return data, err
	}

	return data, nil
}

func getBytesCount(data []byte) int {
	return len(data)
}

func getLinesCount(data []byte) int {
	content := string(data)
	lines := strings.Split(content, "\n")
	return len(lines)
}

func getWordsCount(data []byte) int {
	content := string(data)
	words := strings.Split(content, " ")
	return len(words)
}
