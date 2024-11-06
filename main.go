package main

import (
	"fmt"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/syscall"
	"github.com/traefik/yaegi/stdlib/unsafe"
)

type ILogger interface {
	Info(msg string, args ...any)
}

func main() {
	hookTrue()
}

func pkg() {

	i := interp.New(interp.Options{GoPath: "./_pkg"})
	if err := i.Use(stdlib.Symbols); err != nil {
		panic(err)
	}
	if err := i.Use(syscall.Symbols); err != nil {
		panic(err)
	}
	if err := i.Use(unsafe.Symbols); err != nil {
		panic(err)
	}

	_, err := i.Eval(`import "foo/bar"`)
	if err != nil {
		panic(err)
	}

	prog, err := i.Compile(`bar.NewSample()`)
	if err != nil {
		panic(err)
	}

	fnv, err := i.Execute(prog)
	if err != nil {
		panic(err)
	}

	fn, ok := fnv.Interface().(func(string, string, ILogger) func(string) string)
	if !ok {
		panic("conversion failed")
	}

	closure := fn("arg1", "arg2", nil)
	fmt.Printf(closure("arg3"))
}

type HookInput struct {
	Logger  ILogger
	Message string
}

func pkg2StringerInterface() {
	i := interp.New(interp.Options{GoPath: "./_pkg2"})
	if err := i.Use(stdlib.Symbols); err != nil {
		panic(err)
	}
	if err := i.Use(syscall.Symbols); err != nil {
		panic(err)
	}
	if err := i.Use(unsafe.Symbols); err != nil {
		panic(err)
	}

	_, err := i.Eval(`import "module/hooks"`)
	if err != nil {
		panic(err)
	}

	prog, err := i.Compile(`hooks.Printstr()`)
	if err != nil {
		panic(err)
	}

	fnv, err := i.Execute(prog)
	if err != nil {
		panic(err)
	}

	fn, ok := fnv.Interface().(func(str fmt.Stringer) error)
	if !ok {
		panic("conversion failed")
	}

	err = fn(StringerTest{Str: "stub stringer"})
	if err != nil {
		panic(err)
	}
}

var _ fmt.Stringer = (*StringerTest)(nil)

type StringerTest struct {
	Str string
}

func (s StringerTest) String() string {
	return s.Str
}
