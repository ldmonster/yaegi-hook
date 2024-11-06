package main

import (
	"bytes"
	"os"
	"reflect"

	"github.com/deckhouse/deckhouse/pkg/log"
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

// for quick testing
var hooksSrc = `package hooks

import (
	"go_hook"
)

func DeleteOrphanEndpoints(input *go_hook.HookInput) error {
	snap := input.Snapshots["endpointslices"]

	for _, sn := range snap {
		endpointSliceName := sn.(string)
		input.Logger.Info(endpointSliceName)
		continue
	}

	return nil
}
`

func hookSrc() {
	i := interp.New(interp.Options{})

	//This allows use of standard libs
	i.Use(stdlib.Symbols)

	//This will make a 'custom' lib  available that can be imported and contains your Data struct
	custom := make(map[string]map[string]reflect.Value)
	custom["go_hook/go_hook"] = make(map[string]reflect.Value)
	custom["go_hook/go_hook"]["HookInput"] = reflect.ValueOf((*go_hook.HookInput)(nil))
	i.Use(custom)

	fbyte, err := os.ReadFile("./_pkg3/src/module/hooks/hook.go")
	if err != nil {
		panic(err)
	}

	fbyte = bytes.Replace(fbyte, []byte("github.com/flant/addon-operator/pkg/module_manager/go_hook"), []byte("go_hook"), 1)

	_, err = i.Eval(string(fbyte))
	if err != nil {
		panic(err)
	}

	v, err := i.Eval("hooks.DeleteOrphanEndpoints")
	if err != nil {
		panic(err)
	}

	bar := v.Interface().(func(*go_hook.HookInput) error)

	err = bar(&go_hook.HookInput{
		Logger: log.NewLogger(log.Options{}).With("hook", "first"),
		Snapshots: map[string][]go_hook.FilterResult{
			"endpointslices": {
				"first-slice",
				"second-slice",
			},
		},
	})
	if err != nil {
		panic(err)
	}

	err = bar(&go_hook.HookInput{
		Logger: log.NewLogger(log.Options{}).With("hook", "second"),
		Snapshots: map[string][]go_hook.FilterResult{
			"endpointslices": {
				"1-slice",
				"2-slice",
			},
		},
	})
	if err != nil {
		panic(err)
	}
}
