package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
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
	allnamespaces := kingpin.Flag("all-namespaces", "All Namespaces").Bool()
	config.TTL = kingpin.Flag("ttl", "Time to live of pod in seconds").Short('t').Default("14400").Int32()
	config.Namespace = kingpin.Flag("namespace", "Namespace of the pod we want to duplicate").Short('n').Default("default").String()
	config.Pod = kingpin.Flag("pod", "Pod to duplicate").Short('p').String()
	config.CPU = kingpin.Flag("cpu", "CPU Request for the duplicated Pod").Short('c').String()
	config.Memory = kingpin.Flag("memory", "Memory Request for the duplicated Pod").Short('m').String()
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

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if _, present := pod.ObjectMeta.Labels["job-name"]; !present {
		deployment := find.Deployment(list.Deployments(clientset, config), pod)
		result, err := create.Job(clientset, config, deployment, container)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		jobName := result.ObjectMeta.Name

		s := spinner.New(spinner.CharSets[26], 500*time.Millisecond)
		s.Start()
		watch.Job(ctx, clientset, config, jobName) //wait until Job is created
		watch.Pod(ctx, clientset, config, jobName) //wait until Pod is running
		pod = selector.Pod(list.PodsForJob(clientset, config, jobName), result.Name)
		s.Stop()
	}

	startShell(pod.Name, container.Name, config.Namespace)
}

func startShell(pod, container string, namespace *string) {
	cmd := exec.Command(os.Getenv("SHELL"), "-c", "kubectl attach "+pod+" -n "+*namespace+" -t -i -c "+*namespace+"-"+container+"-exec")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
