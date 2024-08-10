package counter

type Counter struct {
	Counter    map[rune]int
	TotalChars int
}

func (c Counter) Print() {
	for k, v := range c.Counter {
		println(string(k), v)
	}
}

func NewCounter() Counter {
	return Counter{
		Counter:    make(map[rune]int),
		TotalChars: 0,
	}
}
