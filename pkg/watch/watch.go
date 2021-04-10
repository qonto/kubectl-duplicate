package watch

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/qonto/kubectl-duplicate/pkg/config"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// Job waits until the Job from Job is created
func Job(ctx context.Context, clientset *kubernetes.Clientset, config config.Configuration, jobName string) {
	client := clientset.BatchV1().Jobs(*config.Namespace)
	options := metav1.ListOptions{
		// LabelSelector: "job-name=" + jobName,
	}
	watch, err := client.Watch(ctx, options)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := waitForJob(ctx, watch.ResultChan(), jobName); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func waitForJob(ctx context.Context, ch <-chan watch.Event, jobName string) error {
	for {
		select {
		case event := <-ch:
			// when client.Watch uses context with timeout it will return an empty object once it
			// reaches the deadline
			if event.Object == nil {
				return errors.New("Failed to retrieve the state of the pod")
			}
			j, _ := event.Object.(*batchv1.Job)
			if j.ObjectMeta.Name == jobName && j.Status.Active > 0 {
				return nil
			}
			if j.ObjectMeta.Name == jobName && j.Status.Failed > 0 {
				return errors.New("Job in error")
			}
		// we check for context being done also here to make sure the timeout is honored
		case <-ctx.Done():
			return fmt.Errorf("Timeout")
		}
	}
}

// Pod waits until the Pod from Job is running
func Pod(ctx context.Context, clientset *kubernetes.Clientset, config config.Configuration, jobName string) {
	client := clientset.CoreV1().Pods(*config.Namespace)
	options := metav1.ListOptions{
		LabelSelector: "job-name=" + jobName,
	}
	watch, err := client.Watch(ctx, options)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := waitForPod(ctx, watch.ResultChan()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func waitForPod(ctx context.Context, ch <-chan watch.Event) error {
	for {
		select {
		case event := <-ch:
			// when client.Watch uses context with timeout it will return an empty object once it
			// reaches the deadline
			if event.Object == nil {
				return errors.New("Failed to retrieve the state of the pod")
			}
			p, _ := event.Object.(*v1.Pod)
			switch p.Status.Phase {
			case v1.PodRunning:
				return nil
			case v1.PodFailed:
				return errors.New("Pod in error")
			}
			// we check for context being done also here to make sure the timeout is honored
		case <-ctx.Done():
			return fmt.Errorf("Timeout")
		}
	}
}
