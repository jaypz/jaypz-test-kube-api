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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	tasksv1alpha1 "github.io/jaypz/api/v1alpha1"
)

var _ = Describe("DeploymentWithJobRunner Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default", // TODO(user):Modify as needed
		}

		deploymentwithjobrunner := &tasksv1alpha1.DeploymentWithJobRunner{
			ObjectMeta: metav1.ObjectMeta{
				Name:      resourceName,
				Namespace: "default",
			},
		}

		BeforeEach(func() {
			By("creating the custom resource for the Kind DeploymentWithJobRunner")
			err := k8sClient.Get(ctx, typeNamespacedName, deploymentwithjobrunner)
			replicas := int32(1)

			if err != nil && errors.IsNotFound(err) {
				resource := &tasksv1alpha1.DeploymentWithJobRunner{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					Spec: tasksv1alpha1.DeploymentWithJobRunnerSpec{
						Deployment: tasksv1alpha1.Deployment{
							Replicas: &replicas,
							Selector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									"app": resourceName,
								},
							},
							Metadata: tasksv1alpha1.Metadata{
								Labels: map[string]string{
									"app": resourceName,
								},
								Namespace: "default",
								Name:      resourceName,
							},
							PodSpec: v1.PodSpec{
								Containers: []v1.Container{
									{
										Name:  "my-container",
										Image: "nginx",
									},
								},
							},
						},
						Job: tasksv1alpha1.Job{
							Metadata: tasksv1alpha1.Metadata{Name: "hello-world-pod", Namespace: "default"},
							PodSpec: v1.PodSpec{
								Containers: []v1.Container{
									{
										Name:    "hello-world",
										Image:   "busybox",
										Command: []string{"/bin/sh", "-c", "echo 'Hello, World!'"},
									},
								},
								RestartPolicy: v1.RestartPolicyNever,
							},
						},
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			// TODO(user): Cleanup logic after each test, like removing the resource instance.
			resource := &tasksv1alpha1.DeploymentWithJobRunner{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance DeploymentWithJobRunner")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})
		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &DeploymentWithJobRunnerReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
			// Example: If you expect a certain status condition after reconciliation, verify it here.
		})
	})
})
