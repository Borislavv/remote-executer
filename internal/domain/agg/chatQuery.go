package agg

type ChatQuery struct {
	// Filter by `ID`
	//
	// required: false
	// in: query
	// example: `62565ef869ece65cdec4a43c`
	ID string `json:"id" schema:"id"`

	// Filter by `ChatID`
	//
	// required: false
	// in: query
	// example: `62565ef869ece65cdec4a43c`
	ChatID string `json:"chatId" schema:"chatId"`

	// Filter by `Type`
	//
	// required: false
	// in: query
	// example: `private`
	Type string `json:"type" schema:"type"`

	// Sort by field
	//
	// required: false
	// in: query
	// example: text
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

func (q ChatQuery) GetOpts() OptsQuery {
	return OptsQuery{
		SortBy:  q.SortBy,
		OrderBy: q.OrderBy,
		Offset:  q.Offset,
		Limit:   q.Limit,
	}
}
