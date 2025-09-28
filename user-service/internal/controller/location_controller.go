package controller

import (
	"context"

	"github.com/namnv2496/user-service/internal/service"
	userv1 "github.com/namnv2496/user-service/pkg/user_core/v1"
)

type LocationHanlder struct {
	userv1.UnimplementedLocationServiceServer
	trieService service.ITrie
}

func NewLocationHander(
	trieService service.ITrie,
) userv1.LocationServiceServer {
	return &LocationHanlder{
		trieService: trieService,
	}
}

func (c *LocationHanlder) GetLocationMapping(ctx context.Context, req *userv1.GetLocationMappingRequest) (*userv1.GetLocationMappingResponse, error) {
	// result := c.trieService.SearchSuggestion(req.word, req.limit)
	// resp := make(map)
	return nil, nil
}

func (c *LocationHanlder) SearchLocationSuggestion(ctx context.Context, req *userv1.SearchLocationSuggestionRequest) (*userv1.SearchLocationSuggestionResponse, error) {
	result := c.trieService.SearchSuggestion(req.Word, req.Limit)
	locations := make([]*userv1.LocationInfo, 0)

	for _, location := range result {
		locations = append(locations, &userv1.LocationInfo{
			CityV1:       location.CityV1,
			CityV1Name:   location.CityV1Name,
			CityV2:       location.CityV2,
			CityV2Name:   location.CityV2Name,
			District:     location.District,
			DistrictName: location.DistrictName,
			WardV1:       location.WardV1,
			WardV1Name:   location.WardV1Name,
			WardV2:       location.WardV2,
			WardV2Name:   location.WardV2Name,
		})
	}
	return &userv1.SearchLocationSuggestionResponse{
		Location: locations,
	}, nil
}
