package data

import (
	"flag"
	"fmt"
)

type counterPrinter struct {
	fileCounter   DataCounter
	printerParams PrintParams
}

func NewCounterPrinter(printerParams PrintParams) *counterPrinter {
	loadFromFile := len(flag.CommandLine.Args()) > 0
	return &counterPrinter{printerParams: printerParams,
		fileCounter: *ProcessCouting(printerParams.FileName, loadFromFile)}
}

func (p *counterPrinter) Print() {
	var byteCount, linesCount, wordsCount, charsCount int = 0, 0, 0, 0
	var consoleOutput string
	if !p.printerParams.ShowByteCount && !p.printerParams.ShowLinesCount && !p.printerParams.ShowWordsCount && !p.printerParams.ShowCharsCount {
		byteCount = p.fileCounter.GetBytesCount()
		linesCount = p.fileCounter.GetLinesCount()
		wordsCount = p.fileCounter.GetWordsCount()
		consoleOutput = fmt.Sprintf("%d %d %d %s", linesCount, wordsCount, byteCount, p.printerParams.FileName)
		fmt.Println(consoleOutput)
		return
	}

	if p.printerParams.ShowByteCount {
		byteCount = p.fileCounter.GetBytesCount()
		consoleOutput = p.appendCounterFormat(byteCount, consoleOutput)
	}

	if p.printerParams.ShowLinesCount {
		linesCount = p.fileCounter.GetLinesCount()
		consoleOutput = p.appendCounterFormat(linesCount, consoleOutput)
	}

	if p.printerParams.ShowWordsCount {
		wordsCount = p.fileCounter.GetWordsCount()
		consoleOutput = p.appendCounterFormat(wordsCount, consoleOutput)
	}

	if p.printerParams.ShowCharsCount {
		charsCount = p.fileCounter.GetCharsCount()
		consoleOutput = p.appendCounterFormat(charsCount, consoleOutput)
	}

	consoleOutput = fmt.Sprintf("%s %s", consoleOutput, p.printerParams.FileName)
	fmt.Println(consoleOutput)
}

func (p counterPrinter) appendCounterFormat(count int, input string) string {
	return fmt.Sprintf("%s %d", input, count)
}

type PrintParams struct {
	ShowByteCount  bool
	ShowCharsCount bool
	ShowLinesCount bool
	ShowWordsCount bool
	FileName       string
}

func (p PrintParams) String() string {
	return fmt.Sprintf("ShowByteCount: %t, ShowCharsCount: %t, ShowLinesCount: %t, ShowWordsCount: %t, FileName: %s", p.ShowByteCount, p.ShowCharsCount, p.ShowLinesCount, p.ShowWordsCount, p.FileName)
}
