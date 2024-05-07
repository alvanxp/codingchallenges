package data

import "strings"

type FileData []byte

func (f FileData) GetBytesCount() int {
	return len(f)
}

func (f FileData) GetLinesCount() int {
	content := string(f)
	lines := strings.Split(content, "\n")
	return len(lines)
}

func (f FileData) GetWordsCount() int {
	content := string(f)
	words := strings.Split(content, " ")
	return len(words)
}

func (f FileData) GetCharsCount() int {
	content := string(f)
	chars := strings.Split(content, "")
	return len(chars)
}
