package find

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// Deployments finds the Deployment  which has created a Pod
func Deployment(deployments []appsv1.Deployment, pod corev1.Pod) appsv1.Deployment {
	for _, i := range deployments {
		count := 0
		for j, k := range i.Spec.Selector.MatchLabels {
			if pod.ObjectMeta.Labels[j] == k {
				count++
			}
		}
		if len(i.Spec.Selector.MatchLabels) == count {
			return i
		}
	}
	fmt.Println("No deployment found")
	return appsv1.Deployment{}
}
