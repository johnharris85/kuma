package k8s

import (
	"github.com/Kong/konvoy/components/konvoy-control-plane/pkg/core"
	core_discovery "github.com/Kong/konvoy/components/konvoy-control-plane/pkg/core/discovery"
	"github.com/Kong/konvoy/components/konvoy-control-plane/pkg/plugins/discovery/k8s/controllers"

	kube_ctrl "sigs.k8s.io/controller-runtime"
)

func NewDiscoverySource(mgr kube_ctrl.Manager) (core_discovery.DiscoverySource, error) {
	// convert Pods into Dataplanes
	if err := addPodReconciler(mgr); err != nil {
		return nil, err
	}
	// discover Dataplanes
	return addDataplaneReconciler(mgr)
}

func addPodReconciler(mgr kube_ctrl.Manager) error {
	reconciler := &controllers.PodReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Log:    core.Log.WithName("controllers").WithName("Pod"),
	}
	return reconciler.SetupWithManager(mgr)
}

func addDataplaneReconciler(mgr kube_ctrl.Manager) (core_discovery.DiscoverySource, error) {
	reconciler := &controllers.DataplaneReconciler{
		Client: mgr.GetClient(),
		Log:    core.Log.WithName("controllers").WithName("Dataplane"),
	}
	if err := reconciler.SetupWithManager(mgr); err != nil {
		return nil, err
	}
	return reconciler, nil
}
