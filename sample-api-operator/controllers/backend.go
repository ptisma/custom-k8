package controllers

import (
	"context"
	"time"

	appv1 "github.com/ptisma/custom-k8/api/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const backendPort = 8080
const backendServicePort = 8080
const backendImage = "ptisma/simple-api-app"

func backendDeploymentName(v *appv1.SampleAPIApp) string {
	return v.Name + "-backend"
}

func backendServiceName(v *appv1.SampleAPIApp) string {
	return v.Name + "-backend-service"
}

func labels(v *appv1.SampleAPIApp, tier string) map[string]string {
	return map[string]string{
		"app":  "sample-app-api",
		"tier": tier,
	}
}

func (r *SampleAPIAppReconciler) backendDeployment(v *appv1.SampleAPIApp) *appsv1.Deployment {
	labels := labels(v, "backend")
	size := v.Spec.Size
	version := v.Spec.Version

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      backendDeploymentName(v),
			Namespace: v.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           backendImage,
						ImagePullPolicy: corev1.PullAlways,
						Name:            "sample-api-app-service",
						Ports: []corev1.ContainerPort{{
							ContainerPort: backendPort,
							Name:          "sample-api-app",
						}},
						Env: []corev1.EnvVar{
							{
								Name:  "VERSION",
								Value: version,
							},
						},
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(v, dep, r.Scheme)
	return dep
}

func (r *SampleAPIAppReconciler) backendService(v *appv1.SampleAPIApp) *corev1.Service {
	labels := labels(v, "backend")

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      backendServiceName(v),
			Namespace: v.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       backendPort,
				TargetPort: intstr.FromInt(backendPort),
			}},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	controllerutil.SetControllerReference(v, s, r.Scheme)
	return s
}

func (r *SampleAPIAppReconciler) updateBackendStatus(v *appv1.SampleAPIApp) error {
	v.Status.BackendImage = backendImage
	err := r.Client.Status().Update(context.TODO(), v)
	return err
}

func (r *SampleAPIAppReconciler) handleBackendChanges(v *appv1.SampleAPIApp) (*reconcile.Result, error) {
	found := &appsv1.Deployment{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      backendDeploymentName(v),
		Namespace: v.Namespace,
	}, found)
	if err != nil {
		// The deployment may not have been created yet, so requeue
		return &reconcile.Result{RequeueAfter: 5 * time.Second}, err
	}

	size := v.Spec.Size

	if size != *found.Spec.Replicas {
		found.Spec.Replicas = &size
		err = r.Client.Update(context.TODO(), found)
		if err != nil {
			log.Error(err, "Failed to update Deployment.", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return &reconcile.Result{}, err
		}
		// Spec updated - return and requeue
		return &reconcile.Result{Requeue: true}, nil
	}

	return nil, nil
}
