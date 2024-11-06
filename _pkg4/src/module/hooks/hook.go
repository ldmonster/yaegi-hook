package hooks

import (
	"log/slog"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	kind "github.com/flant/addon-operator/pkg/module_manager/models/hooks/kind"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func GetHook() *kind.GoHook {
	return kind.NewGoHook(
		&go_hook.HookConfig{
			Kubernetes: []go_hook.KubernetesConfig{
				{
					Name:         "ns",
					ApiVersion:   "v1",
					Kind:         "Namespace",
					NameSelector: nil,
					LabelSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"foo": "bar"},
					},
					FilterFunc: filterNsName,
				},
			},
		},
		pendingReleaseHandler,
	)
}

func RegisterFunc() (*go_hook.HookConfig, kind.ReconcileFunc) {
	return &go_hook.HookConfig{
		Kubernetes: []go_hook.KubernetesConfig{
			{
				Name:         "ns",
				ApiVersion:   "v1",
				Kind:         "Namespace",
				NameSelector: nil,
				LabelSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"foo": "bar"},
				},
				FilterFunc: filterNsName,
			},
		},
	}, pendingReleaseHandler
}

func filterNsName(obj *unstructured.Unstructured) (go_hook.FilterResult, error) {
	return obj.GetName(), nil
}

func pendingReleaseHandler(input *go_hook.HookInput) error {
	sn := input.Snapshots["ns"]
	for _, s := range sn {
		name := s.(string)
		input.Logger.Info("Found ns with label", slog.String("name", name))
	}

	return nil
}
