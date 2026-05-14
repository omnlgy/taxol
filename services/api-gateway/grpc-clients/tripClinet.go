package grpc_clients

import (
	"os"
	pb "ride-sharing/shared/proto/trip"

	"google.golang.org/grpc"
)

type tripServiceClient struct {
	client pb.TripServiceClient
	conn   *grpc.ClientConn
}

func NewTripServiceClient() (*tripServiceClient, error) {
	tripServiceUrl := os.Getenv("TRIP_SERVICE_URL")
	if tripServiceUrl == "" {
		tripServiceUrl = "trip-service:9093"
	}

	conn, err := grpc.NewClient(tripServiceUrl)
	if err != nil {
		return nil, err
	}

	client := pb.NewTripServiceClient(conn)

	return &tripServiceClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *tripServiceClient) Close() error {
	return c.conn.Close()
}
