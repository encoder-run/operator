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

package cloud

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1alpha1 "github.com/encoder-run/operator/api/cloud/v1alpha1"
)

// StorageReconciler reconciles a Storage object
type StorageReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cloud.encoder.run,resources=storages,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cloud.encoder.run,resources=storages/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cloud.encoder.run,resources=storages/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *StorageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	var storage v1alpha1.Storage
	if err := r.Get(ctx, req.NamespacedName, &storage); err != nil {
		log.Error(err, "unable to fetch Storage")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// First time around ensure that status is NOT_DEPLOYED
	if storage.Status.State == nil {
		state := v1alpha1.StorageStateNotDeployed
		storage.Status.State = &state
		if err := r.Status().Update(ctx, &storage); err != nil {
			return ctrl.Result{}, err
		}
		// Request a requeue to ensure that the status is update before proceeding.
		return ctrl.Result{Requeue: true}, nil
	}

	if err := r.ensureDeployment(ctx, &storage); err != nil {
		log.Error(err, "unable to ensure deployment")
		return ctrl.Result{}, err
	}

	if err := r.ensureDeploymentCleanup(ctx, &storage); err != nil {
		log.Error(err, "unable to ensure deployment cleanup")
		return ctrl.Result{}, err
	}

	if err := r.ensureStatus(ctx, &storage); err != nil {
		log.Error(err, "unable to ensure status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// ensureDeployment ensures
func (r *StorageReconciler) ensureDeployment(ctx context.Context, storage *v1alpha1.Storage) error {
	log := log.FromContext(ctx)
	// Check if the deployment spec exists and is enabled.
	if storage.Spec.Deployment == nil || !storage.Spec.Deployment.Enabled {
		return nil
	}

	// Get the deployment if it exists.
	deployment := &v1.Deployment{}
	if err := r.Get(ctx, client.ObjectKey{Name: storage.Name, Namespace: storage.Namespace}, deployment); err != nil {
		if client.IgnoreNotFound(err) == nil {
			// Create the deployment if it does not exist.
			log.Info("Creating Deployment", "Storage.Namespace", storage.Namespace, "Storage.Name", storage.Name)
			return r.createDeployment(ctx, storage)
		}
		return err
	}

	if deployment.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceCPU] != storage.Spec.Deployment.CPU ||
		deployment.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceMemory] != storage.Spec.Deployment.Memory {
		log.Info("Updating Deployment", "Storage.Namespace", storage.Namespace, "Storage.Name", storage.Name)
		deployment.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceCPU] = storage.Spec.Deployment.CPU
		deployment.Spec.Template.Spec.Containers[0].Resources.Limits[corev1.ResourceMemory] = storage.Spec.Deployment.Memory
		if err := r.Update(ctx, deployment); err != nil {
			return err
		}
		// Update the status of the storage.
		state := v1alpha1.StorageStateDeploying
		storage.Status.State = &state
		// Add condition to the storage.
		storage.Status.Conditions = append(storage.Status.Conditions, metav1.Condition{
			Type:               string(v1alpha1.StorageStateDeploying),
			Status:             metav1.ConditionTrue,
			Reason:             "DeploymentUpdated",
			Message:            "Deployment updated successfully",
			LastTransitionTime: metav1.Now(),
		})
		// Update the status of the storage.
		if err := r.Status().Update(ctx, storage); err != nil {
			return err
		}
	}
	return nil
}

// ensureDeploymentCleanup ensures that the deployment for the storage is deleted
// if deployment spec is disabled.
func (r *StorageReconciler) ensureDeploymentCleanup(ctx context.Context, storage *v1alpha1.Storage) error {
	log := log.FromContext(ctx)
	// Check if the deployment spec exists and is disabled.
	if storage.Spec.Deployment == nil || storage.Spec.Deployment.Enabled {
		return nil
	}

	// Get the deployment if it exists.
	deployment := &v1.Deployment{}
	if err := r.Get(ctx, client.ObjectKey{Name: storage.Name, Namespace: storage.Namespace}, deployment); err != nil {
		// If the deployment is not found, ignore the error.
		if client.IgnoreNotFound(err) != nil {
			return err
		}
	}

	// Delete the deployment.
	log.Info("Deleting Deployment", "Storage.Namespace", storage.Namespace, "Storage.Name", storage.Name)
	if err := r.Delete(ctx, deployment); err != nil {
		// If the deployment is not found, ignore the error.
		if client.IgnoreNotFound(err) != nil {
			return err
		}
	}

	// Get the service if it exists. Ignore if not.
	service := &corev1.Service{}
	if err := r.Get(ctx, client.ObjectKey{Name: storage.Name, Namespace: storage.Namespace}, service); err != nil {
		if client.IgnoreNotFound(err) != nil {
			return err
		}
	}

	// Delete the service.
	log.Info("Deleting Service", "Storage.Namespace", storage.Namespace, "Storage.Name", storage.Name)
	if err := r.Delete(ctx, service); err != nil {
		// If the service is not found, ignore the error.
		if client.IgnoreNotFound(err) != nil {
			return err
		}
	}

	// Update the status of the storage.
	state := v1alpha1.StorageStateNotDeployed
	storage.Status.State = &state
	// Add condition to the storage.
	storage.Status.Conditions = append(storage.Status.Conditions, metav1.Condition{
		Type:               string(v1alpha1.StorageStateNotDeployed),
		Status:             metav1.ConditionTrue,
		Reason:             "DeploymentDeleted",
		Message:            "Deployment deleted successfully",
		LastTransitionTime: metav1.Now(),
	})
	// Update the status of the storage.
	if err := r.Status().Update(ctx, storage); err != nil {
		return err
	}

	return nil
}

// ensureStatus ensures that the status of the storage is updated based on the state of the inference service.
func (r *StorageReconciler) ensureStatus(ctx context.Context, storage *v1alpha1.Storage) error {
	log := log.FromContext(ctx)
	if storage.Spec.Deployment != nil && !storage.Spec.Deployment.Enabled {
		// Get the deployment if it exists.
		deployment := &v1.Deployment{}
		if err := r.Get(ctx, client.ObjectKey{Name: storage.Name, Namespace: storage.Namespace}, deployment); err != nil {
			if client.IgnoreNotFound(err) == nil {
				return nil
			}
			return err
		}

		// Update the status to ready if the ready replicas are greater than 0 and the storage state is not equal to ready.
		if deployment.Status.ReadyReplicas > 0 && storage.Status.State != nil && *storage.Status.State != v1alpha1.StorageStateReady {
			// Update the status of the storage.
			state := v1alpha1.StorageStateReady
			storage.Status.State = &state
			// Add condition to the storage.
			storage.Status.Conditions = append(storage.Status.Conditions, metav1.Condition{
				Type:               string(v1alpha1.StorageStateReady),
				Status:             metav1.ConditionTrue,
				Reason:             "DeploymentReady",
				Message:            "Deployment is ready",
				LastTransitionTime: metav1.Now(),
			})
			// Update the status of the storage.
			log.Info("Storage is ready", "Storage.Namespace", storage.Namespace, "Storage.Name", storage.Name)
			if err := r.Status().Update(ctx, storage); err != nil {
				return err
			}
		}
	} else if storage.Spec.Type == v1alpha1.StorageTypePostgres && storage.Spec.Postgres != nil {
		if storage.Spec.Postgres.External {
			state := v1alpha1.StorageStateReady
			storage.Status.State = &state
			// Add condition to the storage.
			storage.Status.Conditions = append(storage.Status.Conditions, metav1.Condition{
				Type:               string(v1alpha1.StorageStateReady),
				Status:             metav1.ConditionTrue,
				Reason:             "DeploymentReady",
				Message:            "Deployment is ready",
				LastTransitionTime: metav1.Now(),
			})
			// Update the status of the storage.
			log.Info("Storage is ready", "Storage.Namespace", storage.Namespace, "Storage.Name", storage.Name)
			if err := r.Status().Update(ctx, storage); err != nil {
				return err
			}
		}
	}

	return nil
}

// createPassword is a simple but random password used for redis.
func (r *StorageReconciler) createPassword() (string, error) {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	specials := "~=+%^*/()[]{}/!@#$?|"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		digits + specials
	length := 8
	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	buf[1] = specials[rand.Intn(len(specials))]
	for i := 2; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	return string(buf), nil
}

func (r *StorageReconciler) ensurePasswordSecret(ctx context.Context, storage *v1alpha1.Storage) error {
	secret := &corev1.Secret{}
	secretName := types.NamespacedName{Name: storage.Name, Namespace: storage.Namespace}
	if err := r.Get(ctx, secretName, secret); err != nil {
		if errors.IsNotFound(err) {
			// This should only run once unless the deployment is deleted.
			password, err := r.createPassword()
			if err != nil {
				return err
			}
			passwordArg := fmt.Sprintf("--requirepass %s", password)

			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      storage.Name,
					Namespace: storage.Namespace,
				},
				StringData: map[string]string{
					"password":       password,
					"redis-password": passwordArg,
				},
				Type: corev1.SecretTypeOpaque,
			}
			// Set Organization instance as the owner and controller of the secret.
			if err := controllerutil.SetControllerReference(storage, secret, r.Scheme); err != nil {
				return err
			}

			// Create the secret
			if err := r.Create(ctx, secret); err != nil {
				return err
			}
			return nil
		} else {
			// Error reading the secret
			return err
		}
	}
	return nil
}

func (r *StorageReconciler) createDeployment(ctx context.Context, storage *v1alpha1.Storage) error {
	if storage.Spec.Type == v1alpha1.StorageTypeRedis {
		err := r.ensurePasswordSecret(ctx, storage)
		if err != nil {
			return err
		}

		// Define the Persistent Volume Claim name
		pvcName := types.NamespacedName{Name: storage.Name, Namespace: storage.Namespace}

		// Define the Persistent Volume Claim object
		pvc := &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pvcName.Name,
				Namespace: pvcName.Namespace,
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse("100Mi"),
					},
				},
			},
		}
		if err := controllerutil.SetControllerReference(storage, pvc, r.Scheme); err != nil {
			return err
		}

		// Create the Persistent Volume Claim
		if err := r.Create(ctx, pvc); err != nil {
			return err
		}

		// Define the Kubernetes Deployment
		deploy := &v1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      storage.Name,
				Namespace: storage.Namespace,
				Labels:    map[string]string{"app": "redis-stack"},
			},
		}

		// Set the storage as the owner of the deployment.
		if err := controllerutil.SetControllerReference(storage, deploy, r.Scheme); err != nil {
			return err
		}

		// Apply the Deployment to the cluster
		_, err = controllerutil.CreateOrUpdate(ctx, r.Client, deploy, func() error {
			deploy.Spec = v1.DeploymentSpec{
				Replicas: pointer.Int32Ptr(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": "redis-stack"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": "redis-stack"},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "redis-stack",
								Image: "redis/redis-stack:latest",
								Ports: []corev1.ContainerPort{
									{
										ContainerPort: 6379,
									},
									{
										ContainerPort: 8001,
									},
								},
								Env: []corev1.EnvVar{
									{
										Name: "REDIS_ARGS",
										ValueFrom: &corev1.EnvVarSource{
											SecretKeyRef: &corev1.SecretKeySelector{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: storage.Name,
												},
												Key: "redis-password",
											},
										},
									},
								},
								VolumeMounts: []corev1.VolumeMount{
									{
										Name:      "redis-data",
										SubPath:   "redis-data",
										MountPath: "/data",
									},
								},
								Resources: corev1.ResourceRequirements{
									Limits: corev1.ResourceList{
										"cpu":    storage.Spec.Deployment.CPU,
										"memory": storage.Spec.Deployment.Memory,
									},
								},
							},
						},
						Volumes: []corev1.Volume{
							{
								Name: "redis-data",
								VolumeSource: corev1.VolumeSource{
									PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
										ClaimName: pvcName.Name,
									},
								},
							},
						},
					},
				},
			}
			return nil
		})

		if err != nil {
			return err
		}

		// Define the Kubernetes Service for redis-stack
		svc := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      storage.Name,
				Namespace: storage.Namespace,
			},
		}

		// Set the storage as the owner of the service.
		if err := controllerutil.SetControllerReference(storage, svc, r.Scheme); err != nil {
			return err
		}

		// Apply the Service to the cluster
		_, err = controllerutil.CreateOrUpdate(ctx, r.Client, svc, func() error {
			svc.Spec = corev1.ServiceSpec{
				Selector: map[string]string{
					"app": "redis-stack",
				},
				Ports: []corev1.ServicePort{
					{
						Name:       "redis",
						Protocol:   corev1.ProtocolTCP,
						Port:       6379,
						TargetPort: intstr.FromInt(6379),
					},
					{
						Name:       "insight",
						Protocol:   corev1.ProtocolTCP,
						Port:       8001,
						TargetPort: intstr.FromInt(8001),
					},
				},
			}
			return nil
		})

		if err != nil {
			return err
		}

		// Update the status of the storage.
		state := v1alpha1.StorageStateDeploying
		storage.Status.State = &state
		// Add condition to the storage.
		storage.Status.Conditions = append(storage.Status.Conditions, metav1.Condition{
			Type:               string(v1alpha1.StorageStateDeploying),
			Status:             metav1.ConditionTrue,
			Reason:             "InferenceServiceCreated",
			Message:            "Inference service created successfully",
			LastTransitionTime: metav1.Now(),
		})
		// Update the status of the storage.
		if err := r.Status().Update(ctx, storage); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unsupported storage type: %s", storage.Spec.Type)
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *StorageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Storage{}).
		Owns(&v1.Deployment{}).
		Complete(r)
}
