package location

import "github.com/namnv2496/user-service/internal/entity"

type IDictionaryMatcher interface {
	Match(normalizedText string, level string) ([]*entity.LocationInfo, bool)
}

type dictionaryMatcher struct {
	index *entity.LocationIndex
}

func NewDictionaryMatcher(index *entity.LocationIndex) IDictionaryMatcher {
	return &dictionaryMatcher{index: index}
}

func (_self *dictionaryMatcher) Match(normalizedText string, level string) ([]*entity.LocationInfo, bool) {
	switch level {
	case "city":
		if _self.index == nil || _self.index.NormalizedCityName == nil {
			return nil, false
		}

		locs, found := _self.index.NormalizedCityName[normalizedText]
		return locs, found
	case "district":
		if _self.index == nil || _self.index.NormalizedDistrictName == nil {
			return nil, false
		}

		locs, found := _self.index.NormalizedDistrictName[normalizedText]
		return locs, found
	default:
		if _self.index == nil || _self.index.NormalizedWardName == nil {
			return nil, false
		}

		locs, found := _self.index.NormalizedWardName[normalizedText]
		if found {
			return locs, found
		}
		if _self.index == nil || _self.index.NormalizedDistrictName == nil {
			return nil, false
		}

		locs, found = _self.index.NormalizedDistrictName[normalizedText]
		if found {
			return locs, found
		}
		if _self.index == nil || _self.index.NormalizedCityName == nil {
			return nil, false
		}

		locs, found = _self.index.NormalizedCityName[normalizedText]
		if found {
			return locs, found
		}
		return nil, false
	}
}
