package selector

import (
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	corev1 "k8s.io/api/core/v1"
)

// Pod selects a Pod
func Pod(pods []corev1.Pod, match string) corev1.Pod {
	if match != "" {
		for _, pod := range pods {
			if pod.ObjectMeta.Labels["job-name"] == match {
				return pod
			}
			if pod.Name == match {
				return pod
			}
		}
	}

	templates := &promptui.SelectTemplates{
		// Label: `		`,
		Active:   `{{ "> " | cyan | bold }}{{ .ObjectMeta.Name | cyan | bold }}{{if (index .ObjectMeta.Annotations "end-at")}}{{ " [duplicata]" }}{{ end }}`,
		Inactive: `  {{ .ObjectMeta.Name }}{{if (index .ObjectMeta.Annotations "end-at")}}{{ " [duplicata]" }}{{ end }}`,
		Details: `
{{if (index .ObjectMeta.Annotations "end-at")}}{{ " End: " }}{{ index .ObjectMeta.Annotations "end-at" | bold }}{{ end }}`,
	}

	searcher := func(input string, index int) bool {
		p := pods[index]
		Name := strings.ToLower(p.ObjectMeta.Name) + strings.ToLower(p.ObjectMeta.Labels["created-by"])
		input = strings.ToLower(input)

		return strings.Contains(Name, input)
	}

	prompt := promptui.Select{
		Label:             "Pods",
		Items:             pods,
		Templates:         templates,
		Size:              10,
		Searcher:          searcher,
		HideSelected:      true,
		StartInSearchMode: true,
	}

	selected, _, err := prompt.Run()
	if err != nil {
		os.Exit(0)
	}

	return pods[selected]
}

// Container selects a Container
func Container(pod corev1.Pod) corev1.Container {
	if len(pod.Spec.Containers) == 1 {
		return pod.Spec.Containers[0]
	}

	templates := &promptui.SelectTemplates{
		// Label: `		`,
		Active:   `{{ "> " | cyan | bold }}{{ .Name | cyan | bold }}`,
		Inactive: `  {{ .Name }}`,
	}

	searcher := func(input string, index int) bool {
		j := pod.Spec.Containers[index]
		Name := strings.ToLower(j.Name)
		input = strings.ToLower(input)

		return strings.Contains(Name, input)
	}

	prompt := promptui.Select{
		Label:             "Containers",
		Items:             pod.Spec.Containers,
		Templates:         templates,
		Size:              10,
		Searcher:          searcher,
		HideSelected:      true,
		StartInSearchMode: true,
	}

	selected, _, err := prompt.Run()
	if err != nil {
		os.Exit(0)
	}

	return pod.Spec.Containers[selected]
}
