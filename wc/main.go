package main

import (
	"bufio"
	"ccwc/data"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	filename, _ := getFilePath()
	showByteCount := flag.Bool("c", false, "show the number of bytes in the file")
	showLinesCount := flag.Bool("l", false, "show the number of lines in the file")
	showWordsCount := flag.Bool("w", false, "show the number of words in the file")
	showCharsCount := flag.Bool("m", false, "show the number of characters in the file")
	flag.Parse()

	var byteCount, linesCount, wordsCount, charsCount int = 0, 0, 0, 0

	var consoleOutput string
	data, err := getDataFromFile(filename)
	if err != nil {
		fmt.Println("Error processing file:", err)
		return
	}

	if !*showByteCount && !*showLinesCount && !*showWordsCount && !*showCharsCount {
		byteCount = data.GetBytesCount()
		linesCount = data.GetLinesCount()
		wordsCount = data.GetWordsCount()
		consoleOutput = fmt.Sprintf("%d %d %d %s", linesCount, wordsCount, byteCount, filename)
		fmt.Println(consoleOutput)
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

	if *showCharsCount {
		charsCount = data.GetCharsCount()
		consoleOutput = AppendCounterFormat(charsCount, consoleOutput)
	}

	consoleOutput = fmt.Sprintf("%s %s", consoleOutput, filename)
	fmt.Println(consoleOutput)
}

func AppendCounterFormat(count int, input string) string {
	return fmt.Sprintf("%s %d", input, count)
}

func getDataFromFile(filename string) (data.FileData, error) {
	temp := flag.CommandLine.Args()
	if len(temp) == 0 {
		reader := bufio.NewReader(os.Stdin)
		data, err := io.ReadAll(io.Reader(reader))
		if err != nil {
			return nil, err
		}
		return data, nil
	}
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
