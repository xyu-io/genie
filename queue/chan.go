package queue

var _ Queue = (*Chann)(nil)

type Chann struct {
	Ch chan interface{}
}

// NewChan news a channel for simulate the queue
func NewChan() *Chann {
	return &Chann{
		Ch: make(chan interface{}),
	}
}

func (c *Chann) Pop() interface{} {
	return <-c.Ch
}

func (c *Chann) PopBatch(max int) []interface{} {
	if max > len(c.Ch) {
		max = len(c.Ch)
	}
	data := make([]interface{}, 0, max)
	for i := 0; i < max; i++ {
		data = append(data, <-c.Ch)
	}
	return data
}

func (c *Chann) Push(value interface{}) bool {
	c.Ch <- value
	return true
}

func (c *Chann) PushBatch(vs []interface{}) bool {
	for _, v := range vs {
		c.Ch <- v
	}
	return true
}

func (c *Chann) Len() int {
	return len(c.Ch)
}
