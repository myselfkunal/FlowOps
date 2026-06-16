package main

import (
	"context"
	"log"
	"net"
	"google.golang.org/grpc"
	pb "github.com/myselfkunal/FlowOps/proto/gen/auth"
	"github.com/myselfkunal/FlowOps/shared"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
}

func (s *AuthService) Validate(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	tracer := otel.Tracer("auth")
	_, span := tracer.Start(ctx, "Validate")
	defer span.End()

	log.Printf("Received auth request: %v", req)
	
	if req.AuthToken == "valid-token-123" {
		return &pb.AuthResponse{Success: true}, nil
	}
	return &pb.AuthResponse{Success: false}, nil
}

func main() {
	shutdown := shared.InitTracer("auth")
	defer shutdown()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
    grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	pb.RegisterAuthServiceServer(grpcServer, &AuthService{})
	log.Printf("Server listening on %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}