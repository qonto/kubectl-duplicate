package main

import (
	"fmt"
	"os"
	"os/exec"

	cfg "github.com/qonto/kubectl-duplicate/pkg/config"
	"github.com/qonto/kubectl-duplicate/pkg/create"
	"github.com/qonto/kubectl-duplicate/pkg/find"
	"github.com/qonto/kubectl-duplicate/pkg/list"
	"github.com/qonto/kubectl-duplicate/pkg/selector"
	"github.com/qonto/kubectl-duplicate/pkg/watch"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/tools/clientcmd"
)

var config cfg.Configuration
var version string
var commit string

func init() {
	allnamespaces := kingpin.Flag("all-namespaces", "All Namespace").Bool()
	config.TTL = kingpin.Flag("ttl", "Time to live of pods is seconds").Short('t').Default("14400").Int32()
	config.Namespace = kingpin.Flag("namespace", "Namespace").Short('n').Default("default").String()
	config.Pod = kingpin.Flag("pod", "Pod").Short('p').String()
	config.Shell = kingpin.Flag("shell", "Shell to use").Short('s').Default("sh").String()
	config.CPU = kingpin.Flag("cpu", "cpu").Short('c').String()
	config.Memory = kingpin.Flag("memory", "Memory").Short('m').String()
	config.Kubeconfig = kingpin.Flag("kubeconfig", "Kube config file (override by env var KUBECONFIG").Short('k').Default(os.Getenv("HOME") + "/.kube/config").ExistingFile()
	v := kingpin.Flag("version", "Print version").Short('v').Bool()
	kingpin.Parse()

	if *v {
		fmt.Printf("Version: %s\nCommit: %s\n", version, commit)
		os.Exit(0)
	}

	if *allnamespaces {
		*config.Namespace = ""
	}

	if kc := os.Getenv("KUBECONFIG"); kc != "" {
		*config.Kubeconfig = kc
	}
}

func main() {
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", *config.Kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	pods := list.Pods(clientset, config)
	pod := selector.Pod(pods, *config.Pod)
	container := selector.Container(pod)

	if _, present := pod.ObjectMeta.Labels["job-name"]; !present {
		deployment := find.Deployment(list.Deployments(clientset, config), pod)
		result, err := create.Job(clientset, config, deployment, container)
		if err != nil {
			panic(err)
		}
		jobName := result.ObjectMeta.Name

		watch.Pod(clientset, config, jobName) //wait until Pod is running
		pod = selector.Pod(list.Pods(clientset, config), result.Name)
	}

	startShell(pod.Name, config.Namespace)
}

func startShell(pod string, namespace *string) {
	cmd := exec.Command(os.Getenv("SHELL"), "-c", "kubectl exec -n "+*namespace+" -ti "+pod+" -- /bin/"+*config.Shell)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil && *config.Shell != "sh" {
		fmt.Println("-- fallback to \"sh\" shell -- ")
		s := "sh"
		config.Shell = &s
		startShell(pod, namespace)
	}
}
