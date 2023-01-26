package script

import "fmt"

const UNDEF = "undefined"

type Value float64

func (v Value) String() string {
	if v == -1 {
		return fmt.Sprint(UNDEF)
	}
	return fmt.Sprintf("%v", float64(v))
}
