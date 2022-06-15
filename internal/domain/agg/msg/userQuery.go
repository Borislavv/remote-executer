package agg

import agg "github.com/Borislavv/remote-executer/internal/domain/agg/mongo"

type UserQuery struct {
	// Filter by `ID`
	//
	// required: false
	// in: query
	// example: `62565ef869ece65cdec4a43c`
	ID string `json:"id" schema:"id"`

	// Filter by `Firstname`
	//
	// required: false
	// in: query
	// example: `Jared`
	Firstname string `json:"firstname" schema:"firstname"`

	// Filter by `Lastname`
	//
	// required: false
	// in: query
	// example: `Jackson`
	Lastname string `json:"lastname" schema:"lastname"`

	// Filter by `Username`
	//
	// required: false
	// in: query
	// example: `JaredsonUsername`
	Username string `json:"username" schema:"username"`

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

func (q UserQuery) GetOpts() agg.OptsQuery {
	return agg.OptsQuery{
		SortBy:  q.SortBy,
		OrderBy: q.OrderBy,
		Offset:  q.Offset,
		Limit:   q.Limit,
	}
}
