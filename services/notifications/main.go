package main

import (
	"context"
	"log"
	"net"
	"google.golang.org/grpc"
	pb "github.com/myselfkunal/FlowOps/proto/gen/notifications"
	"github.com/myselfkunal/FlowOps/shared"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
)

type NotificationServer struct {
	pb.UnimplementedNotificationsServiceServer
}

func (s *NotificationServer) SendNotifications(ctx context.Context, req *pb.NotificationsRequest) (*pb.NotificationsResponse, error) {
	tracer := otel.Tracer("notifications")
	_, span := tracer.Start(ctx, "SendNotifications")
	defer span.End()
	
	log.Printf("Received notification request: %v", req)

	return &pb.NotificationsResponse {
		Sent: true,
	}, nil
}

func main(){
	shutdown := shared.InitTracer("notifications")
	defer shutdown()
	
	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}	

	grpcServer := grpc.NewServer(
    	grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	pb.RegisterNotificationsServiceServer(grpcServer, &NotificationServer{})
	
	log.Println("Notifications service is running on :50054")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}	
}