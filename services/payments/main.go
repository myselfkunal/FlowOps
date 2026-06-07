package main

import (
	"context"
	"log"
	"net"
	"google.golang.org/grpc"
	pb "github.com/myselfkunal/FlowOps/proto/gen/payments"
)

type PaymentService struct {
	pb.UnimplementedPaymentServiceServer
}

func (s *PaymentService) ProcessPayment(ctx context.Context, req *pb.ProcessPaymentRequest) (*pb.ProcessPaymentResponse, error) {
	log.Printf("Received payment request: %v", req)

	return &pb.ProcessPaymentResponse {
		Success: true,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	pb.RegisterPaymentServiceServer(grpcServer, &PaymentService{})

	log.Println("Payments service is running on :50053")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}