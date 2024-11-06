package hooks

import (
	"fmt"
)

func Printlol() func(str string) error {
	return func(str string) error {
		fmt.Println("lol")

		return nil
	}
}
