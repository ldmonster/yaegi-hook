package hooks

import (
	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
)

func SomeHook(input *go_hook.HookInput) error {
	input.Logger.Info("hook started")

	input.Logger.Info("hook ended")

	return nil
}

func DeleteOrphanEndpoints(input *go_hook.HookInput) error {
	snap := input.Snapshots["endpointslices"]

	for _, sn := range snap {
		endpointSliceName := sn.(string)
		input.Logger.Info(endpointSliceName)
	}

	return nil
}
