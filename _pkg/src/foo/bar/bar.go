package bar

import (
	"fmt"
)

type ILogger interface {
	Info(args ...interface{})
}

var version = "v1"

func NewSample() func(string, string, ILogger) func(string) string {
	fmt.Println("in NewSample")
	return func(val string, name string, logger ILogger) func(string) string {
		logger.Info("kek")
		fmt.Println("in function", version, val, name)
		return func(msg string) string {
			return fmt.Sprint("here", version, val, name, msg)
		}
	}
}
