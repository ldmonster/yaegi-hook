package hooks

import (
	"fmt"
)

func Printstr() func(str fmt.Stringer) error {
	return func(str fmt.Stringer) error {
		fmt.Println(str.String())

		return nil
	}
}
