package domain

type Sort struct {
	// Example values: "name", "date", "price".
	SortBy string `json:"sort_by,omitempty"`
	// Example values: "asc" for ascending, "desc" for descending.
	Order string `json:"order,omitempty"`
	// distance by km
	// Example values: "5"
	Latitude float64 `json:"latitude,omitempty"`
	// Example values: "106.7247336079002"
	Longtitude float64 `json:"longtitude,omitempty"`
	// Example values: "10.73026809889869"
	Distance int64 `json:"distance,omitempty"`
}
