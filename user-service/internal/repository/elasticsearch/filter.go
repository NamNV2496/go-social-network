package elasticsearch

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"github.com/namnv2496/user-service/internal/domain"
)

var (
	CreateAt = "create_at"
	Distance = "distance"
)

func getQueryByValue[K int64 | float64 | string | bool](field string, value K) types.Query {
	return types.Query{
		Term: map[string]types.TermQuery{
			field: {
				Value: value,
			},
		},
	}
}

func getQueryByArray[K int64 | float64 | string | bool](field string, value []K) types.Query {
	if len(value) == 1 {
		return types.Query{
			Term: map[string]types.TermQuery{
				field: {
					Value: value[0],
				},
			},
		}
	}
	return types.Query{
		Terms: &types.TermsQuery{
			TermsQuery: map[string]types.TermsQueryField{
				field: value,
			},
		},
	}
}

func getQuerySort(sortOpt domain.Sort) (types.SortCombinations, error) {
	if sortOpt.Order != sortorder.Asc.Name && sortOpt.Order != sortorder.Desc.Name {
		return nil, fmt.Errorf("sort type is not support: %s", sortOpt.Order)
	}
	var order *sortorder.SortOrder
	if sortOpt.Order == sortorder.Asc.Name {
		order = &sortorder.Asc
	} else {
		order = &sortorder.Desc
	}
	switch sortOpt.SortBy {
	case CreateAt:
		return &types.SortOptions{
			SortOptions: map[string]types.FieldSort{
				sortOpt.SortBy: {
					Order: order,
				},
			},
		}, nil
	case Distance:
		if sortOpt.Distance == 0 || sortOpt.Longtitude == 0 || sortOpt.Latitude == 0 {
			return nil, nil
		}
		return &types.SortOptions{
			GeoDistance_: &types.GeoDistanceSort{
				GeoDistanceSort: map[string][]types.GeoLocation{
					"location": {fmt.Sprintf("%f,%f", sortOpt.Latitude, sortOpt.Longtitude)},
				},
			},
		}, nil
	default:
		return &types.SortOptions{
			SortOptions: map[string]types.FieldSort{
				sortOpt.SortBy: {
					Order: order,
				},
			},
		}, nil
	}
}
