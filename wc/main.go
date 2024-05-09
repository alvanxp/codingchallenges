package main

import (
	"ccwc/data"
	"flag"
	"os"
)

func main() {
	printerParams := data.PrinterParams{
		ShowByteCount:  *flag.Bool("c", false, "show the number of bytes in the file"),
		ShowLinesCount: *flag.Bool("l", false, "show the number of lines in the file"),
		ShowWordsCount: *flag.Bool("w", false, "show the number of words in the file"),
		ShowCharsCount: *flag.Bool("m", false, "show the number of characters in the file"),
		FileName:       getFilePath(),
	}

	flag.Parse()
	printer := data.NewCounterPrinter(printerParams)
	printer.Print()

}

func getFilePath() string {
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i][0] != '-' {
			return os.Args[i]
		}
	}
	return ""
}
