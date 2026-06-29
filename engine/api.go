package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
)

type APIServer struct {
    k8sClient *kubernetes.Clientset
    store     *Store
    namespace string
}

type ServiceStatus struct {
    Name            string `json:"name"`
    DesiredReplicas int32  `json:"desired_replicas"`
    ReadyReplicas   int32  `json:"ready_replicas"`
    Image           string `json:"image"`
    Status          string `json:"status"`
}

func NewAPIServer(k8sClient *kubernetes.Clientset, store *Store, namespace string) *APIServer {
    return &APIServer{
        k8sClient: k8sClient,
        store:     store,
        namespace: namespace,
    }
}

func (a *APIServer) Start(port string) {
    mux := http.NewServeMux()
    mux.HandleFunc("/status",  a.handleStatus)
    mux.HandleFunc("/history", a.handleHistory)
    mux.HandleFunc("/config",  a.handleConfig)

    log.Printf("API server starting on :%s", port)
    if err := http.ListenAndServe(":"+port, mux); err != nil {
        log.Fatalf("API server failed: %v", err)
    }
}

func (a *APIServer) handleStatus(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")

    services := []string{"payments", "auth", "orders", "notifications", "api-gateway"}
    var statuses []ServiceStatus

    for _, name := range services {
        dep, err := a.k8sClient.AppsV1().Deployments(a.namespace).Get(
            context.Background(), name, metav1.GetOptions{},
        )
        if err != nil {
            statuses = append(statuses, ServiceStatus{
                Name:   name,
                Status: "unknown",
            })
            continue
        }

        status := "healthy"
        if dep.Status.ReadyReplicas < *dep.Spec.Replicas {
            status = "degraded"
        }

        statuses = append(statuses, ServiceStatus{
            Name:            name,
            DesiredReplicas: *dep.Spec.Replicas,
            ReadyReplicas:   dep.Status.ReadyReplicas,
            Image:           dep.Spec.Template.Spec.Containers[0].Image,
            Status:          status,
        })
    }

    json.NewEncoder(w).Encode(statuses)
}

func (a *APIServer) handleHistory(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    json.NewEncoder(w).Encode(a.store.GetEvents())
}

func (a *APIServer) handleConfig(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")

    raw, err := FetchConfig(repoOwner, repoName, configFile)
    if err != nil {
        http.Error(w, "failed to fetch config", http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(map[string]string{"config": raw})
}