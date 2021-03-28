package create

import (
	"context"
	"time"

	"github.com/qonto/kubectl-duplicate/pkg/config"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Job creates the Job for Duplicata
func Job(clientset *kubernetes.Clientset, config config.Configuration, deployment appsv1.Deployment, container corev1.Container) (*batchv1.Job, error) {
	execAction := new(corev1.ExecAction)
	execAction.Command = []string{"true"}
	probe := new(corev1.Probe)
	probe.Handler.Exec = execAction
	probe.InitialDelaySeconds = int32(*config.TTL)

	container.Name = deployment.ObjectMeta.Name + "-exec"
	container.Command = []string{"/bin/sh", "-c", "--", "/bin/sh"}
	container.Ports = []corev1.ContainerPort{}
	container.VolumeMounts = []corev1.VolumeMount{}
	container.LivenessProbe = probe
	container.ReadinessProbe = probe
	container.Stdin = true
	container.TTY = true

	if *config.CPU != "" {
		container.Resources.Limits["cpu"] = resource.MustParse(*config.CPU)
		container.Resources.Requests["cpu"] = resource.MustParse(*config.CPU)
	}
	if *config.Memory != "" {
		container.Resources.Limits["memory"] = resource.MustParse(*config.Memory)
		container.Resources.Requests["memory"] = resource.MustParse(*config.Memory)
	}

	var backoffLimit int32
	var activeDeadlineSeconds, terminationGracePeriodSeconds int64
	var ttlAfterFinished int32
	activeDeadlineSeconds = int64(*config.TTL)
	ttlAfterFinished = int32(3600)
	backoffLimit = 0
	terminationGracePeriodSeconds = 1

	endAt := time.Now().Local().Add(time.Duration(*config.TTL) * time.Second)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: deployment.ObjectMeta.Name + "-duplicata-",
			Annotations: map[string]string{
				"end-at": endAt.Format("2006-01-02 15:04:05"),
			},
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "duplicate",
			},
		},
		Spec: batchv1.JobSpec{
			ActiveDeadlineSeconds:   &activeDeadlineSeconds,
			TTLSecondsAfterFinished: &ttlAfterFinished,
			BackoffLimit:            &backoffLimit,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: deployment.ObjectMeta.Name + "-duplicata-",
					Annotations: map[string]string{
						"end-at": endAt.Format("2006-01-02 15:04:05"),
					},
					Labels: map[string]string{
						"app.kubernetes.io/managed-by": "duplicate",
					},
				},
				Spec: corev1.PodSpec{
					Containers:                    []corev1.Container{container},
					RestartPolicy:                 corev1.RestartPolicyNever,
					TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
				},
			},
		},
	}

	client := clientset.BatchV1().Jobs(deployment.ObjectMeta.Namespace)
	result, err := client.Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	return result, nil
}
