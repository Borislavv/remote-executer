package agg

import (
	"time"

	agg "github.com/Borislavv/remote-executer/internal/domain/agg/mongo"
)

type MsgQuery struct {
	// Filter by `ID`
	//
	// required: false
	// in: query
	// example: `62565ef869ece65cdec4a43c`
	ID string `json:"id" schema:"id"`

	// Filter by `Title`
	//
	// required: false
	// in: query
	// example: `Hello world`
	Text string `json:"title" schema:"title"`

	// Filter by `UpdateId`
	//
	// required: false
	// in: query
	// example: `506233478d`
	UpdateId int64 `json:"updateId" schema:"updateId"`

	// (EQUALS)
	//	Filter by `Date of stat. from file`
	//
	// pattern: \d{4}-\d{2}-\d{2}
	// required: false
	// in: query
	// example: 2021-11-19
	Date time.Time `json:"date" schema:"date"`

	// (RANGE FROM)
	//	Filter by `Date of stat. from file` where it is greater than or equals to `from`
	//
	// pattern: \d{4}-\d{2}-\d{2}
	// required: false
	// in: query
	// example: 2021-11-19
	DateFrom time.Time `json:"dateFrom" schema:"dateFrom"`

	// (RANGE TO)
	//	Filter by `Date of stat. from file` where it is less than or equals to `to`
	//
	// pattern: \d{4}-\d{2}-\d{2}
	// required: false
	// in: query
	// example: 2021-11-21
	DateTo time.Time `json:"dateTo" schema:"dateTo"`

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

func (q MsgQuery) GetOpts() agg.OptsQuery {
	return agg.OptsQuery{
		SortBy:  q.SortBy,
		OrderBy: q.OrderBy,
		Offset:  q.Offset,
		Limit:   q.Limit,
	}
}
