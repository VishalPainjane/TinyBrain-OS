package controllers

import (
	"context"

	corev1alpha1 "github.com/VishalPainjane/TinyBrain-OS/internal/k8s/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// TaskReconciler reconciles a Task object
type TaskReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=core.tinybrain.io,resources=tasks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core.tinybrain.io,resources=tasks/status,verbs=get;update;patch

func (r *TaskReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var task corev1alpha1.Task
	if err := r.Get(ctx, req.NamespacedName, &task); err != nil {
		logger.Error(err, "unable to fetch Task")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.Info("Reconciling Task", "name", task.Name, "agentRef", task.Spec.AgentRef)

	// Here we would interact with the TinyBrain OS Event Bus/Scheduler.
	
	if task.Status.State == "" {
		task.Status.State = "Pending"
		if err := r.Status().Update(ctx, &task); err != nil {
			logger.Error(err, "unable to update Task status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TaskReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.Task{}).
		Complete(r)
}

