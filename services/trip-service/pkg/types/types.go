package types

import (
	pb "ride-sharing/shared/proto/trip"
)

type OsrmApiResponse struct {
	Routes []struct {
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
		Geometry struct {
			Coordinates [][]float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"routes"`
}

func (o *OsrmApiResponse) ToProto() *pb.Route {
	route := o.Routes[0]
	geometry := route.Geometry.Coordinates
	coordinates := make([]*pb.Coordinate, len(geometry))

	for i, val := range geometry {
		coordinates[i] = &pb.Coordinate{
			Latitude:  val[0],
			Longitude: val[1],
		}
	}

	return &pb.Route{
		Geometry: []*pb.Geometry{
			{
				Coordinates: coordinates,
			},
		},
		Distance: route.Distance,
		Duration: route.Duration,
	}
}

type PriceConfig struct {
	PricePerUnitOfDistance float64
	PriceingPerMinute      float64
}

func DefaultPriceConfig() *PriceConfig {
	return &PriceConfig{
		PricePerUnitOfDistance: 1.5,
		PriceingPerMinute:      0.25,
	}
}
