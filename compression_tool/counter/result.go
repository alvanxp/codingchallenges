package counter

type Counter struct {
	Counter map[rune]uint
}

func (c Counter) Print() {
	for k, v := range c.Counter {
		println(string(k), v)
	}
}

func NewCounter() Counter {
	return Counter{make(map[rune]uint)}
}
