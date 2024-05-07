package data_test

import (
	"ccwc/data"
	"testing"
)

func TestGetLinesCount(t *testing.T) {
	fileContent := []byte("Hello\nWorld\nThis\nIs\nA\nTest")
	data := data.FileData(fileContent) // Import the missing package and use the FileData function
	expectedLinesCount := 6

	linesCount := data.GetLinesCount()

	if linesCount != expectedLinesCount {
		t.Errorf("Expected lines count: %d, but got: %d", expectedLinesCount, linesCount)
	}
	t.Log("TestGetLinesCount passed")
}

func TestGetWordsCount(t *testing.T) {
	fileContent := []byte("Hello World This Is A Test")
	data := data.FileData(fileContent) // Import the missing package and use the FileData function
	expectedWordsCount := 6

	wordsCount := data.GetWordsCount()

	if wordsCount != expectedWordsCount {
		t.Errorf("Expected words count: %d, but got: %d", expectedWordsCount, wordsCount)
	}
	t.Log("TestGetWordsCount passed")
}

func TestGetCharsCount(t *testing.T) {
	fileContent := []byte("Hello World This Is A Test")
	data := data.FileData(fileContent) // Import the missing package and use the FileData function
	expectedCharsCount := 26

	charsCount := data.GetCharsCount()

	if charsCount != expectedCharsCount {
		t.Errorf("Expected chars count: %d, but got: %d", expectedCharsCount, charsCount)
	}
	t.Log("TestGetCharsCount passed")
}

func TestGetBytesCount(t *testing.T) {
	fileContent := []byte("Hello World This Is A Test")
	data := data.FileData(fileContent) // Import the missing package and use the FileData function
	expectedBytesCount := 26

	bytesCount := data.GetBytesCount()

	if bytesCount != expectedBytesCount {
		t.Errorf("Expected bytes count: %d, but got: %d", expectedBytesCount, bytesCount)
	}
	t.Log("TestGetBytesCount passed")
}
