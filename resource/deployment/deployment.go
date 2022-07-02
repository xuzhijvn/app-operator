package deployment

import (
	appv1 "github.com/xuzhijvn/app-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func New(app *appv1.App) *appsv1.Deployment {
	labels := map[string]string{"app.example.com/v1": app.Name}
	selector := &metav1.LabelSelector{MatchLabels: labels}
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(app, schema.GroupVersionKind{
					Group:   appv1.GroupVersion.Group,
					Version: appv1.GroupVersion.Version,
					Kind:    "App",
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: app.Spec.Replicas,
			Selector: selector,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: newContainers(app),
				},
			},
		},
	}
}

func newContainers(app *appv1.App) []corev1.Container {
	var containerPorts []corev1.ContainerPort
	for _, servicePort := range app.Spec.Ports {
		var cport corev1.ContainerPort
		cport.ContainerPort = servicePort.TargetPort.IntVal
		containerPorts = append(containerPorts, cport)
	}
	return []corev1.Container{
		{
			Name:            app.Name,
			Image:           app.Spec.Image,
			Ports:           containerPorts,
			Env:             app.Spec.Envs,
			Resources:       app.Spec.Resources,
			ImagePullPolicy: corev1.PullIfNotPresent,
		},
	}
}
