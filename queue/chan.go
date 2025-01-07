package queue

var _ Queue = (*Chann)(nil)

type Chann struct {
	ch chan interface{}
}

// NewChan news a channel for simulate the queue
func NewChan() *Chann {
	return &Chann{
		ch: make(chan interface{}),
	}
}

func (c *Chann) Pop() interface{} {
	return <-c.ch
}

func (c *Chann) PopBatch(max int) []interface{} {
	if max > len(c.ch) {
		max = len(c.ch)
	}
	data := make([]interface{}, 0, max)
	for i := 0; i < max; i++ {
		data = append(data, <-c.ch)
	}
	return data
}

func (c *Chann) Push(value interface{}) bool {
	c.ch <- value
	return true
}

func (c *Chann) PushBatch(vs []interface{}) bool {
	for _, v := range vs {
		c.ch <- v
	}
	return true
}

func (c *Chann) Len() int {
	return len(c.ch)
}
