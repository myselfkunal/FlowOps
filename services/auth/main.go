package main

import (
	"context"
	"log"
	"net"
	"google.golang.org/grpc"
	pb "github.com/myselfkunal/FlowOps/proto/gen/auth"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
}

func (s *AuthService) Validate(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	if req.AuthToken == "valid-token-123" {
		return &pb.AuthResponse{Success: true}, nil
	}
	return &pb.AuthResponse{Success: false}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, &AuthService{})
	log.Printf("Server listening on %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}