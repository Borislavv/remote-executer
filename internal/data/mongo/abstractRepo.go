package mongoRepo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	agg "github.com/Borislavv/remote-executer/internal/domain/agg/mongo"
)

func GetOpts(query agg.OptsableQuery) *options.FindOptions {
	q := query.GetOpts()

	opts := options.Find()

	if q.Offset > 0 {
		opts.SetSkip(q.Offset)
	}

	if q.Limit > 0 {
		opts.SetLimit(q.Limit)
	}

	if q.SortBy != "" {
		sort := bson.M{}
		switch q.OrderBy {
		case "", "asc":
			sort[q.SortBy] = 1
		case "desc":
			sort[q.SortBy] = -1
		default:
			return opts
		}
		opts.SetSort(sort)
	}

	return opts
}
