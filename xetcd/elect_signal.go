package xetcd

type (
	ElectSignal struct {
		signal chan int
	}
)

func NewElectSignal(num int) *ElectSignal {
	if num < 1 {
		num = 1
	}
	ret := &ElectSignal{
		signal: make(chan int, num),
	}
	for i := 0; i < num; i++ {
		ret.signal <- 1
	}
	return ret
}
func (c *ElectSignal) GetSignal() bool {
	select {
	case <-c.signal:
		return true
	default:
		return false
	}
}
func (c *ElectSignal) Release() bool {
	select {
	case c.signal <- 1:
		return true
	default:
		return false
	}
}
