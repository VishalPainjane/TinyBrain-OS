package controllers

import (
	"context"

	memoryv1alpha1 "github.com/VishalPainjane/TinyBrain-OS/internal/k8s/api/memory/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// SwapPolicyReconciler reconciles a SwapPolicy object
type SwapPolicyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=memory.tinybrain.io,resources=swappolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=memory.tinybrain.io,resources=swappolicies/status,verbs=get;update;patch

func (r *SwapPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var policy memoryv1alpha1.SwapPolicy
	if err := r.Get(ctx, req.NamespacedName, &policy); err != nil {
		logger.Error(err, "unable to fetch SwapPolicy")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.Info("Reconciling SwapPolicy", "name", policy.Name, "threshold", policy.Spec.Threshold)

	// Update status as a simple test
	if !policy.Status.Active {
		policy.Status.Active = true
		if err := r.Status().Update(ctx, &policy); err != nil {
			logger.Error(err, "unable to update SwapPolicy status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SwapPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&memoryv1alpha1.SwapPolicy{}).
		Complete(r)
}

