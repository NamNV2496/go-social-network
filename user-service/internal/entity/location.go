package entity

type LocationRequest struct {
	CityV1   string `json:"city_v1,omitempty"`
	CityV2   string `json:"city_v2,omitempty"`
	District string `json:"district,omitempty"`
	WardV1   string `json:"ward_v1,omitempty"`
	WardV2   string `json:"ward_v2,omitempty"`
	IsNew    bool   `json:"is_new,omitempty"`
}

type LocationInfo struct {
	CityV1       string   `json:"city_v1,omitempty"`
	CityV1Name   string   `json:"city_v1_name,omitempty"`
	CityV2       string   `json:"city_v2,omitempty"`
	CityV2Name   string   `json:"city_v2_name,omitempty"`
	District     string   `json:"district,omitempty"`
	DistrictName string   `json:"district_name,omitempty"`
	WardV1       string   `json:"ward_v1,omitempty"`
	WardV1Name   string   `json:"ward_v1_name,omitempty"`
	WardV2       string   `json:"ward_v2,omitempty"`
	WardV2Name   string   `json:"ward_v2_name,omitempty"`
	BuilderName  string   `json:"builder_name,omitempty"`
	Variants     []string `json:"variants,omitempty"`
}

type LocationIndex struct {
	NormalizedName         map[string][]*LocationInfo // for all
	NormalizedWardName     map[string][]*LocationInfo
	NormalizedDistrictName map[string][]*LocationInfo
	NormalizedCityName     map[string][]*LocationInfo
}

type MatchResult struct {
	Location    *LocationInfo
	MatchedText string
	Score       float64
	MatchType   string // "dictionary" or "fuzzy"
}
