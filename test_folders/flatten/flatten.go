package flatten

type Interface interface {
	Sum(x, y int) int
	SumString(x string, y string) string
}

type OneWay struct {
	i int
	s string
}

func (o *OneWay) Sum(x int, y int) int {
	return o.i + x + y
}

func (o *OneWay) SumString(x, y string) string {
	return o.s + x + y
}
