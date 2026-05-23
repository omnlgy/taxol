package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	tripTypes "ride-sharing/services/trip-service/pkg/types"
	"ride-sharing/shared/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	repo domain.TripRepository
}

func NewService(repo domain.TripRepository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error) {
	t := &domain.TripModel{
		ID:       primitive.NewObjectID(),
		UserID:   fare.UserID,
		Status:   "panding",
		RideFare: *fare,
	}
	return s.repo.CreateTrip(ctx, t)
}

func (s *service) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*tripTypes.OsrmApiResponse, error) {
	url := fmt.Sprintf("http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson",
		pickup.Longitude, pickup.Latitude,
		destination.Longitude, destination.Latitude)

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from OSRM API: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response from OSRM API: %v", err)
	}

	var routeRes tripTypes.OsrmApiResponse
	if err := json.Unmarshal(body, &routeRes); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal response from OSRM API: %v", err)
	}

	return &routeRes, nil
}

func (s *service) EstimatePackagesPriceWithRoute(route *tripTypes.OsrmApiResponse) []*domain.RideFareModel {
	baseFare := getBaseFare()
	estimatedFare := make([]*domain.RideFareModel, len(baseFare))

	for i, fare := range estimatedFare {
		estimatedFare[i] = estimatedFareRoute(fare, route)
	}

	return estimatedFare
}

func (s *service) GenerateTripFares(ctx context.Context, rideFares []*domain.RideFareModel, userID string) ([]*domain.RideFareModel, error) {
	fares := make([]*domain.RideFareModel, len(rideFares))

	for i, rideFare := range rideFares {
		id := primitive.NewObjectID()

		fare := &domain.RideFareModel{
			ID:                id,
			UserID:            userID,
			TotalPriceInCents: rideFare.TotalPriceInCents,
			PackageSlug:       rideFare.PackageSlug,
		}

		if err := s.repo.SaveRiderFare(ctx, fare); err != nil {
			return nil, fmt.Errorf("Failed to save trip fare : %w", err)
		}

		fares[i] = fare
	}

	return fares, nil
}

func estimatedFareRoute(rideFare *domain.RideFareModel, route *tripTypes.OsrmApiResponse) *domain.RideFareModel {
	priceConfig := tripTypes.DefaultPriceConfig()
	carPackagePrice := rideFare.TotalPriceInCents

	distance := route.Routes[0].Distance
	duration := route.Routes[0].Duration

	distanceFare := distance * priceConfig.PricePerUnitOfDistance
	timeFare := duration * priceConfig.PriceingPerMinute
	totalFare := carPackagePrice + distanceFare + timeFare

	return &domain.RideFareModel{
		TotalPriceInCents: totalFare,
		PackageSlug:       rideFare.PackageSlug,
	}
}

func getBaseFare() []*domain.RideFareModel {
	return []*domain.RideFareModel{
		{
			PackageSlug:       "suv",
			TotalPriceInCents: 200,
		},
		{
			PackageSlug:       "sedan",
			TotalPriceInCents: 350,
		},
		{
			PackageSlug:       "van",
			TotalPriceInCents: 400,
		},
		{
			PackageSlug:       "luxury",
			TotalPriceInCents: 1000,
		},
	}
}
