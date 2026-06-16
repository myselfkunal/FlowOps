package main

import (
	"context"
	"log"
	"net"
	"google.golang.org/grpc"
	pb "github.com/myselfkunal/FlowOps/proto/gen/payments"
	"github.com/myselfkunal/FlowOps/shared"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

type PaymentService struct {
	pb.UnimplementedPaymentServiceServer
}

func (s *PaymentService) ProcessPayment(ctx context.Context, req *pb.ProcessPaymentRequest) (*pb.ProcessPaymentResponse, error) {
    tracer := otel.Tracer("payments")
    _, span := tracer.Start(ctx, "ProcessPayment")
    defer span.End()

    log.Printf("Processing payment of amount: %v", req.Amount)
    return &pb.ProcessPaymentResponse{Success: true}, nil
}

func main() {
	shutdown := shared.InitTracer("payments")
	defer shutdown()

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	
	grpcServer := grpc.NewServer(
    	grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	pb.RegisterPaymentServiceServer(grpcServer, &PaymentService{})

	log.Println("Payments service is running on :50053")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}