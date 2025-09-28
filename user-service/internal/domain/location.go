package domain

type LocationInfo struct {
	CityV1       string `json:"city_v1,omitempty"`
	CityV1Name   string `json:"city_v1_name,omitempty"`
	CityV2       string `json:"city_v2,omitempty"`
	CityV2Name   string `json:"city_v2_name,omitempty"`
	District     string `json:"district,omitempty"`
	DistrictName string `json:"district_name,omitempty"`
	WardV1       string `json:"ward_v1,omitempty"`
	WardV1Name   string `json:"ward_v1_name,omitempty"`
	WardV2       string `json:"ward_v2,omitempty"`
	WardV2Name   string `json:"ward_v2_name,omitempty"`
}
