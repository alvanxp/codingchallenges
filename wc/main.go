package main

import (
	"ccwc/data"
	"flag"
	"os"
)

func main() {

	var showByteCount, showLinesCount, showWordsCount, showCharsCount bool

	flag.BoolVar(&showByteCount, "c", false, "show the number of bytes in the file")
	flag.BoolVar(&showLinesCount, "l", false, "show the number of lines in the file")
	flag.BoolVar(&showWordsCount, "w", false, "show the number of words in the file")
	flag.BoolVar(&showCharsCount, "m", false, "show the number of characters in the file")

	flag.Parse()
	printerParams := data.PrintParams{
		ShowByteCount:  showByteCount,
		ShowLinesCount: showLinesCount,
		ShowWordsCount: showWordsCount,
		ShowCharsCount: showCharsCount,
		FileName:       getFilePath(),
	}
	// fmt.Println(printerParams)
	countPrinter := data.NewCounterPrinter(printerParams)
	countPrinter.Print()
}

func getFilePath() string {
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i][0] != '-' {
			return os.Args[i]
		}
	}
	return ""
}
