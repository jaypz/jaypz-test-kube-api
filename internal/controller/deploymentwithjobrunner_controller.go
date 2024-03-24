/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	tasksv1alpha1 "github.io/jaypz/api/v1alpha1"
)

// DeploymentWithJobRunnerReconciler reconciles a DeploymentWithJobRunner object
type DeploymentWithJobRunnerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tasks.github.io,resources=deploymentwithjobrunners,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tasks.github.io,resources=deploymentwithjobrunners/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tasks.github.io,resources=deploymentwithjobrunners/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DeploymentWithJobRunner object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *DeploymentWithJobRunnerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	deployWithJob := &tasksv1alpha1.DeploymentWithJobRunner{}
	if err := r.Get(ctx, req.NamespacedName, deployWithJob); err != nil {
		logger.Error(err, "Failed to find DeploymentWithJobRunner object")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	deployment := &deployWithJob.Spec.Deployment

	if err := r.Create(ctx, deployment); err != nil {
		logger.Error(err, "failed to apply Deployment")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	// Check if the Deployment is ready
	if !isDeploymentReady(deployment) {
		logger.Info("Deployment not yet ready. Requeuing...")
		return reconcile.Result{Requeue: true}, nil
	}

	job := &deployWithJob.Spec.Job

	if err := r.Create(ctx, job); err != nil {
		logger.Error(err, "Job failed to create")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !isJobSuccess(job) {
		logger.Info("Job not yet finished. Requeuing...")
		return reconcile.Result{Requeue: true}, nil
	}

	return ctrl.Result{}, nil
}

// isDeploymentReady checks if the Deployment is ready
func isDeploymentReady(deployment *appsv1.Deployment) bool {
	return deployment.Status.ReadyReplicas > 0
}

func isJobSuccess(job *batchv1.Job) bool {
	return job.Status.Succeeded > 0
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeploymentWithJobRunnerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tasksv1alpha1.DeploymentWithJobRunner{}).
		Complete(r)
}
