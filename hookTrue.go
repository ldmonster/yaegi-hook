package main

import (
	"bytes"
	"fmt"
	"os"
	"reflect"

	"github.com/deckhouse/deckhouse/pkg/log"
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/pkg/module_manager/models/hooks/kind"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func hookTrue() {
	i := interp.New(interp.Options{})

	//This allows use of standard libs
	i.Use(stdlib.Symbols)

	//This will make a 'custom' lib  available that can be imported and contains your Data struct
	custom := map[string]map[string]reflect.Value{
		"go_hook/go_hook": {
			"HookInput":        reflect.ValueOf((*go_hook.HookInput)(nil)),
			"HookConfig":       reflect.ValueOf((*go_hook.HookConfig)(nil)),
			"KubernetesConfig": reflect.ValueOf((*go_hook.KubernetesConfig)(nil)),
			"FilterResult":     reflect.ValueOf((*go_hook.FilterResult)(nil)),
		},
		"metav1/metav1": {
			"LabelSelector": reflect.ValueOf((*metav1.LabelSelector)(nil)),
		},
		"addon-operator-sdk/addon-operator-sdk": {
			"RegisterFunc": reflect.ValueOf((func(config *go_hook.HookConfig, reconcileFunc kind.ReconcileFunc) bool)(nil)),
		},
		"unstructured/unstructured": {
			"Unstructured": reflect.ValueOf((*unstructured.Unstructured)(nil)),
		},
		"kind/kind": {
			"GoHook":        reflect.ValueOf((*kind.GoHook)(nil)),
			"NewGoHook":     reflect.ValueOf(kind.NewGoHook),
			"ReconcileFunc": reflect.ValueOf((kind.ReconcileFunc)(nil)),
		},
	}
	i.Use(custom)

	fbyte, err := os.ReadFile("./_pkg4/src/module/hooks/hook.go")
	if err != nil {
		panic(err)
	}

	fbyte = bytes.Replace(fbyte, []byte("github.com/flant/addon-operator/pkg/module_manager/go_hook"), []byte("go_hook"), 1)
	fbyte = bytes.Replace(fbyte, []byte("github.com/flant/addon-operator/sdk"), []byte("addon-operator-sdk"), 1)
	fbyte = bytes.Replace(fbyte, []byte("k8s.io/apimachinery/pkg/apis/meta/v1"), []byte("metav1"), 1)
	fbyte = bytes.Replace(fbyte, []byte("k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"), []byte("unstructured"), 1)
	fbyte = bytes.Replace(fbyte, []byte("github.com/flant/addon-operator/pkg/module_manager/models/hooks/kind"), []byte("kind"), 1)

	_, err = i.Eval(string(fbyte))
	if err != nil {
		panic(err)
	}

	prog, err := i.Compile(`hooks.GetHook`)
	if err != nil {
		panic(err)
	}

	fnv, err := i.Execute(prog)
	if err != nil {
		panic(err)
	}

	bar := fnv.Interface().(func() *kind.GoHook)

	gh := bar()

	fmt.Printf("%+v\n", gh.GetConfig())
	fmt.Printf("error: %+v\n", gh.Run(&go_hook.HookInput{
		Snapshots: go_hook.Snapshots{
			"ns": {
				"default",
				"some ns",
			},
		},
		Logger: log.NewLogger(log.Options{}),
	}))

	// register func like on webhook file

	prog, err = i.Compile(`hooks.RegisterFunc`)
	if err != nil {
		panic(err)
	}

	fnv, err = i.Execute(prog)
	if err != nil {
		panic(err)
	}

	regbar := fnv.Interface().(func() (*go_hook.HookConfig, kind.ReconcileFunc))

	hcfg, rf := regbar()
	gh = kind.NewGoHook(hcfg, rf)

	fmt.Printf("%+v\n", gh.GetConfig())
	fmt.Printf("error: %+v\n", gh.Run(&go_hook.HookInput{
		Snapshots: go_hook.Snapshots{
			"ns": {
				"default",
				"some ns",
			},
		},
		Logger: log.NewLogger(log.Options{}),
	}))
}
