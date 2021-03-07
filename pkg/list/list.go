package list

import (
	"context"
	"strings"

	"github.com/qonto/kubectl-duplicate/pkg/config"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Pods lists Pods
func Pods(clientset *kubernetes.Clientset, config config.Configuration) []corev1.Pod {
	client := clientset.CoreV1().Pods(*config.Namespace)
	options := metav1.ListOptions{
		FieldSelector: "status.phase=Running",
	}
	result, err := client.List(context.TODO(), options)
	if err != nil {
		panic(err)
	}
	return sort(result.Items)
}

// Deployments lists Deployments
func Deployments(clientset *kubernetes.Clientset, config config.Configuration) []appsv1.Deployment {
	client := clientset.AppsV1().Deployments(*config.Namespace)
	result, err := client.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	return result.Items
}

func sort(pods []corev1.Pod) []corev1.Pod {
	var sortedPods []corev1.Pod
	for _, i := range pods {
		if strings.Contains(i.ObjectMeta.Name, "-duplicata-") {
			sortedPods = append(sortedPods, i)
		}
	}
	for _, i := range pods {
		if !strings.Contains(i.ObjectMeta.Name, "-duplicata-") {
			sortedPods = append(sortedPods, i)
		}
	}
	return sortedPods
}
