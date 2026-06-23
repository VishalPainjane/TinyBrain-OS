package controllers

import (
	"context"

	memoryv1alpha1 "github.com/VishalPainjane/TinyBrain-OS/internal/k8s/api/memory/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// KVCacheReconciler reconciles a KVCache object
type KVCacheReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=memory.tinybrain.io,resources=kvcaches,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=memory.tinybrain.io,resources=kvcaches/status,verbs=get;update;patch

func (r *KVCacheReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var cache memoryv1alpha1.KVCache
	if err := r.Get(ctx, req.NamespacedName, &cache); err != nil {
		logger.Error(err, "unable to fetch KVCache")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.Info("Reconciling KVCache", "name", cache.Name, "size", cache.Spec.Size)

	// Update status as a simple test
	if cache.Status.UsedSize == "" {
		cache.Status.UsedSize = "0"
		if err := r.Status().Update(ctx, &cache); err != nil {
			logger.Error(err, "unable to update KVCache status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KVCacheReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&memoryv1alpha1.KVCache{}).
		Complete(r)
}

