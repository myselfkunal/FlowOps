package main

import (
    //"context"
    "encoding/json"
    "log"
    "net/http"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    authpb   "github.com/myselfkunal/FlowOps/proto/gen/auth"
    orderspb "github.com/myselfkunal/FlowOps/proto/gen/orders"
    "github.com/myselfkunal/FlowOps/shared"
    "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

type OrderRequest struct {
    Token    string  `json:"token"`
    UserID   string  `json:"user_id"`
    ItemName string  `json:"item_name"`
    Amount   float32 `json:"amount"`
}

var (
    authClient   authpb.AuthServiceClient
    ordersClient orderspb.OrderServiceClient
)

func handleOrder(w http.ResponseWriter, r *http.Request) {
    
	// 1. decode request body
	var req OrderRequest
	json.NewDecoder(r.Body).Decode(&req)
    
    // 2. call auth service
    authRes, err := authClient.Validate(r.Context(), &authpb.AuthRequest{
        AuthToken: req.Token,
    })
    if err != nil || !authRes.Success {
        // return 401 unauthorized
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
        return
    }

    // 3. call orders service
    orderRes, err := ordersClient.CreateOrder(r.Context(), &orderspb.OrderRequest{
        UserId:   req.UserID,
        ItemName: req.ItemName,
        Amount:   req.Amount,
    })
	if err != nil {
		// return 500 internal server error
    	w.WriteHeader(http.StatusInternalServerError)
    	json.NewEncoder(w).Encode(map[string]string{"error": "order failed"})
    	return
	}

    // 4. return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
    "success": true,
    "order_id": orderRes.OrderId,
	})
}

func main() {
    shutdown := shared.InitTracer("api-gateway")
    defer shutdown()
    
	// connect to auth service
	authConn, err := grpc.NewClient("localhost:50051",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
    )
	if err != nil {
		log.Fatalf("Failed to connect to auth service: %v", err)
	}
	defer authConn.Close()

	// connect to orders service
	ordersConn, err := grpc.NewClient("localhost:50052",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
    )
	if err != nil {
		log.Fatalf("Failed to connect to orders service: %v", err)
	}
	defer ordersConn.Close()

	// create clients
	authClient = authpb.NewAuthServiceClient(authConn)
	ordersClient = orderspb.NewOrderServiceClient(ordersConn)

	// start HTTP server
	http.HandleFunc("/order", handleOrder)
	log.Println("API Gateway is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}