package controllers

import (
	"context"

	corev1alpha1 "github.com/VishalPainjane/TinyBrain-OS/internal/k8s/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// AgentReconciler reconciles a Agent object
type AgentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=core.tinybrain.io,resources=agents,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core.tinybrain.io,resources=agents/status,verbs=get;update;patch

func (r *AgentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var agent corev1alpha1.Agent
	if err := r.Get(ctx, req.NamespacedName, &agent); err != nil {
		logger.Error(err, "unable to fetch Agent")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.Info("Reconciling Agent", "name", agent.Name, "model", agent.Spec.Model)

	// Here we would interact with the TinyBrain OS Agent Registry.
	// However, per instructions, we do not import k8s packages in runtime/scheduler.
	// The operator should call runtime APIs (e.g. gRPC/HTTP/Go interfaces) to manage agents.
	
	// Update status as a simple test
	if agent.Status.State == "" {
		agent.Status.State = "Pending"
		if err := r.Status().Update(ctx, &agent); err != nil {
			logger.Error(err, "unable to update Agent status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AgentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.Agent{}).
		Complete(r)
}

