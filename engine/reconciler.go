package main

import (
    "context"
    "fmt"
    "log"
	"time"
    "gopkg.in/yaml.v3"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
)

type ServiceConfig struct {
    Name     string `yaml:"name"`
    Image    string `yaml:"image"`
    Replicas int32  `yaml:"replicas"`
}

type FlowOpsConfig struct {
    Services []ServiceConfig `yaml:"services"`
}

// parses the raw YAML string into a FlowOpsConfig struct
func ParseConfig(rawYAML string) (*FlowOpsConfig, error) {
	var config FlowOpsConfig
	err := yaml.Unmarshal([]byte(rawYAML), &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	return &config, nil
}

// compares config against cluster and applies changes
func Reconcile(config *FlowOpsConfig, k8sClient *kubernetes.Clientset, namespace string, store *Store) error {
    for _, svc := range config.Services {
        // get deployment
        deployment, err := k8sClient.AppsV1().Deployments(namespace).Get(
            context.Background(),
            svc.Name,
            metav1.GetOptions{},
        )
        if err != nil {
            log.Printf("failed to get deployment %s: %v", svc.Name, err)
            continue
        }

        // check image
		currentImage := deployment.Spec.Template.Spec.Containers[0].Image
		if currentImage != svc.Image {
			log.Printf("updating image %s: %s -> %s", svc.Name, currentImage, svc.Image)
			deployment.Spec.Template.Spec.Containers[0].Image = svc.Image
			_, err = k8sClient.AppsV1().Deployments(namespace).Update(
				context.Background(),
				deployment,
				metav1.UpdateOptions{},
			)
			if err != nil {
				log.Printf("failed to update image %s: %v", svc.Name, err)
				continue
			}
			log.Printf("successfully updated image %s", svc.Name)
			store.AddEvent(ReconcileEvent{
				Timestamp:   time.Now(),
				ServiceName: svc.Name,
				WhatChanged: "image",
				OldValue:    currentImage,
				NewValue:    svc.Image,
			})
		} else {
			log.Printf("%s image is up to date", svc.Name)
		}

		// check replicas
		currentReplicas := *deployment.Spec.Replicas
		if currentReplicas != svc.Replicas {
			log.Printf("scaling %s: %d -> %d replicas", svc.Name, currentReplicas, svc.Replicas)
			deployment.Spec.Replicas = &svc.Replicas
			_, err = k8sClient.AppsV1().Deployments(namespace).Update(
				context.Background(),
				deployment,
				metav1.UpdateOptions{},
			)
			if err != nil {
				log.Printf("failed to scale deployment %s: %v", svc.Name, err)
				continue
			}
			log.Printf("successfully scaled %s", svc.Name)
			store.AddEvent(ReconcileEvent{
				Timestamp:   time.Now(),
				ServiceName: svc.Name,
				WhatChanged: "replicas",
				OldValue:    fmt.Sprintf("%d", currentReplicas),
				NewValue:    fmt.Sprintf("%d", svc.Replicas),
			})
		}else{
			log.Printf("%s replicas are up to date", svc.Name)
		}
    }
    return nil
}