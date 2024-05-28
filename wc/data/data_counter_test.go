package data_test

import (
	"ccwc/data"
	"testing"
)

func TestGetLinesCount(t *testing.T) {
	dataCounter := data.ProcessCouting("../test.txt", true)
	expectedLinesCount := 7145
	linesCount := dataCounter.GetLinesCount()

	if linesCount != expectedLinesCount {
		t.Errorf("Expected lines count: %d, but got: %d", expectedLinesCount, linesCount)
	}
	t.Log("TestGetLinesCount passed")
}

func TestGetWordsCount(t *testing.T) {
	dataCounter := data.ProcessCouting("../test.txt", true)
	expectedWordsCount := 58164

	wordsCount := dataCounter.GetWordsCount()

	if wordsCount != expectedWordsCount {
		t.Errorf("Expected words count: %d, but got: %d", expectedWordsCount, wordsCount)
	}
	t.Log("TestGetWordsCount passed")
}

func TestGetCharsCount(t *testing.T) {
	dataCounter := data.ProcessCouting("../test.txt", true)
	expectedCharsCount := 339292

	charsCount := dataCounter.GetCharsCount()

	if charsCount != expectedCharsCount {
		t.Errorf("Expected chars count: %d, but got: %d", expectedCharsCount, charsCount)
	}
	t.Log("TestGetCharsCount passed")
}

func TestGetBytesCount(t *testing.T) {
	dataCount := data.ProcessCouting("../test.txt", true)

	expectedBytesCount := 342190

	bytesCount := dataCount.GetBytesCount()

	if bytesCount != expectedBytesCount {
		t.Errorf("Expected bytes count: %d, but got: %d", expectedBytesCount, bytesCount)
	}
	t.Log("TestGetBytesCount passed")
}
