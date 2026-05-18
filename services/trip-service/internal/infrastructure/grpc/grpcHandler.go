package grpc

import (
	"context"
	"fmt"
	"ride-sharing/services/trip-service/internal/domain"
	pb "ride-sharing/shared/proto/trip"
	"ride-sharing/shared/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcHandler struct {
	pb.UnimplementedTripServiceServer

	service domain.TripService
}

func NewGrpcHandler(server *grpc.Server, service domain.TripService) *grpcHandler {
	handler := &grpcHandler{
		service: service,
	}

	pb.RegisterTripServiceServer(server, handler)
	return handler
}

func (handler *grpcHandler) PreviewTrip(ctx context.Context, req *pb.PreviewTripRequest) (*pb.PreviewTripResponse, error) {
	pickup := &types.Coordinate{
		Latitude:  req.StartLocation.Latitude,
		Longitude: req.StartLocation.Longitude,
	}
	destination := &types.Coordinate{
		Latitude:  req.EndLocation.Latitude,
		Longitude: req.EndLocation.Longitude,
	}

	route, err := handler.service.GetRoute(ctx, pickup, destination)
	if err != nil {
		fmt.Println(err)
		return nil, status.Errorf(codes.Internal, "Error getting route %v", err)
	}

	return &pb.PreviewTripResponse{
		Route:     route.ToProto(),
		RideFares: []*pb.RideFare{},
	}, nil
}
