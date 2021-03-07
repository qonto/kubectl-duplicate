package watch

import (
	"context"
	"errors"

	"github.com/qonto/kubectl-duplicate/pkg/config"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Pod waits until the Pod from Job is running
func Pod(clientset *kubernetes.Clientset, config config.Configuration, jobName string) {
	client := clientset.CoreV1().Pods(*config.Namespace)
	options := metav1.ListOptions{
		LabelSelector: "job-name=" + jobName,
	}
	watch, err := client.Watch(context.TODO(), options)
	if err != nil {
		panic(err)
	}
L:
	for {
		event := <-watch.ResultChan()
		p, _ := event.Object.(*v1.Pod)
		switch p.Status.Phase {
		case v1.PodRunning:
			break L
		case v1.PodFailed:
			panic(errors.New("Pod in error"))
		default:
			continue
		}
	}
}
