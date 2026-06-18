package main

import (
	"context"
	"fmt"
	"time"
	"net"
	"log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "github.com/myselfkunal/FlowOps/proto/gen/orders"
	pb2 "github.com/myselfkunal/FlowOps/proto/gen/payments"
	pb3 "github.com/myselfkunal/FlowOps/proto/gen/notifications"
    "github.com/myselfkunal/FlowOps/shared"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

type OrderServer struct {
    pb.UnimplementedOrderServiceServer
    paymentsClient  pb2.PaymentServiceClient
    notificationsClient pb3.NotificationsServiceClient
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *pb.OrderRequest) (*pb.OrderResponse, error) {
    tracer := otel.Tracer("orders")
    _, span := tracer.Start(ctx, "CreateOrder")
    defer span.End()
    
    // 1. generate order id
    orderID := fmt.Sprintf("order-%d", time.Now().Unix())

    // 2. call payments
    paymentRes, err := s.paymentsClient.ProcessPayment(ctx, &pb2.ProcessPaymentRequest{
        Amount: req.Amount,
    })
    if err != nil || !paymentRes.Success {
        return &pb.OrderResponse{Success: false}, nil
    }

    // 3. call notifications
    s.notificationsClient.SendNotifications(ctx, &pb3.NotificationsRequest{
        UserId:  req.UserId,
        Success: true,
    })

    // 4. return success
    return &pb.OrderResponse{
        OrderId: orderID,
        Success: true,
    }, nil
}

func main() {
    shutdown := shared.InitTracer("orders")
    defer shutdown()

	// connect to payments
    paymentsConn, err := grpc.NewClient("payments:50053",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
    )
    if err != nil {
        log.Fatalf("Failed to connect to payments: %v", err)
    }
    defer paymentsConn.Close()

    // connect to notifications
    notifConn, err := grpc.NewClient("notifications:50054",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
    )
    if err != nil {
        log.Fatalf("Failed to connect to notifications: %v", err)
    }
    defer notifConn.Close()

	// create clients
	paymentsClient := pb2.NewPaymentServiceClient(paymentsConn)
	notificationsClient := pb3.NewNotificationsServiceClient(notifConn)

	// create your server with clients inside
    orderServer := &OrderServer{
        paymentsClient: paymentsClient,
        notificationsClient: notificationsClient,
    }

    // start gRPC server as usual
    lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
        grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	pb.RegisterOrderServiceServer(grpcServer, orderServer)

	log.Println("Orders service is running on :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

