package main

import (
    "log"
    "time"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

const (
    repoOwner  = "myselfkunal"
    repoName   = "FlowOps"
    configFile = "flowops-config.yaml"
)

func main() {
	// load kubeconfig from default location (~/.kube/config)
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
    	loadingRules, configOverrides)

	restConfig, err := kubeConfig.ClientConfig()
	if err != nil {
		log.Fatalf("failed to load kubeconfig: %v", err)
	}
	k8sClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		log.Fatalf("failed to create k8s client: %v", err)
	}

	for {
		log.Println("reconciling...")
		
		// fetch config
		rawYAML, err := FetchConfig(repoOwner, repoName, configFile)
		if err != nil {
			log.Printf("failed to fetch config: %v", err)
			time.Sleep(30 * time.Second)
			continue
		}

		// parse config
		config, err := ParseConfig(rawYAML)
		if err != nil {
			log.Printf("failed to parse config: %v", err)
			time.Sleep(30 * time.Second)
			continue
		}

		// reconcile
		err = Reconcile(config, k8sClient, "default")
		if err != nil {
			log.Printf("failed to reconcile: %v", err)
		}
		
		time.Sleep(30 * time.Second)
	}

}