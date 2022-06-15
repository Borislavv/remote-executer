package agg

type OptsQuery struct {
	// Sort by field
	//
	// required: false
	// in: query
	// example: host
	SortBy string `json:"sortBy" schema:"sortBy"`

	// Order by asc/desc
	//
	// pattern: (asc|desc)
	// required: false
	// in: query
	// example: asc
	OrderBy string `json:"orderBy" schema:"orderBy"`

	// Offset items from start by value
	//
	// min: 0
	// required: false
	// in: query
	// example: 100
	Offset int64 `json:"offset" schema:"offset"`

	// Sets the number of items to return
	//
	// required: false
	// in: query
	// example: 50
	Limit int64 `json:"limit" schema:"limit"`
}
