package counter

type Counter struct {
	Counter   map[rune]int
	CharCount int
}

func (c Counter) Print() {
	for k, v := range c.Counter {
		println(string(k), v)
	}
}

func NewCounter() Counter {
	return Counter{
		Counter:   make(map[rune]int),
		CharCount: 0,
	}
}
