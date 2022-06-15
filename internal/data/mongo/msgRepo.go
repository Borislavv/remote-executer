package mongoRepo

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	agg "github.com/Borislavv/remote-executer/internal/domain/agg/msg"
	"github.com/Borislavv/remote-executer/internal/util"
)

type MsgRepo struct {
	coll *mongo.Collection
	mu   *sync.Mutex
	buf  []interface{}
}

func NewMsgRepo(collection *mongo.Collection) *MsgRepo {
	return &MsgRepo{
		coll: collection,
		mu:   &sync.Mutex{},
		buf:  []interface{}{},
	}
}

func (r *MsgRepo) InsertMany(ctx context.Context, msgs []agg.Msg) error {
	r.mu.Lock()

	r.buf = r.buf[:0]
	for _, fundStat := range msgs {
		r.buf = append(r.buf, fundStat)
	}
	result, err := r.coll.InsertMany(ctx, r.buf, options.InsertMany())
	if err != nil {
		return util.ErrWithTrace(err)
	}
	if len(result.InsertedIDs) == 0 {
		return util.ErrWithTrace(errors.New("MsgRepo.InsertMany: no one `Msg` document was created"))
	}

	r.mu.Unlock()

	return nil
}

func (r *MsgRepo) Find(ctx context.Context, q agg.MsgQuery) ([]agg.Msg, error) {
	response := []agg.Msg{}

	filter, err := r.GetFilter(q)
	if err != nil {
		return response, util.ErrWithTrace(err)
	}

	cursor, err := r.coll.Find(ctx, filter, GetOpts(q))
	if err != nil {
		return response, util.ErrWithTrace(err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &response); err != nil {
		return response, util.ErrWithTrace(err)
	}

	return response, nil
}

func (r *MsgRepo) MarkAsExecuted(ctx context.Context, msg agg.Msg) error {
	res, err := r.coll.UpdateByID(ctx, msg.ID, bson.M{"$set": bson.M{"executed": true}})
	if err != nil {
		return util.ErrWithTrace(err)
	}
	fmt.Println(res.MatchedCount, res.ModifiedCount)
	return nil
}

func (r *MsgRepo) GetOffset(ctx context.Context) (int64, error) {
	var resp agg.Msg

	q := agg.MsgQuery{
		SortBy:  "updateId",
		OrderBy: "desc",
		Limit:   1,
	}

	f, err := r.GetFilter(q)
	if err != nil {
		return 0, util.ErrWithTrace(err)
	}

	c, err := r.coll.Find(ctx, f, GetOpts(q))
	if err != nil {
		return 0, util.ErrWithTrace(err)
	}
	defer c.Close(ctx)

	if !c.Next(ctx) {
		// handle the case when doc. was not found
		return 0, nil
	}
	if err := c.Decode(&resp); err != nil {
		return 0, util.ErrWithTrace(err)
	}

	return resp.Msg.UpdateId + 1, nil
}

func (r *MsgRepo) GetFilter(q agg.MsgQuery) (bson.M, error) {
	f := bson.M{}

	if q.ID != "" {
		id, err := primitive.ObjectIDFromHex(q.ID)
		if err != nil {
			return f, util.ErrWithTrace(err)
		}

		f["_id"] = bson.M{"$eq": id}
	}

	if q.Text != "" {
		f["text"] = bson.M{"$eq": q.Text}
	}

	if q.UpdateId != 0 {
		f["updateId"] = bson.M{"$eq": q.UpdateId}
	}

	if q.ByExecuted {
		f["executed"] = bson.M{"$eq": q.Executed}
	}

	emptyTamestamp := (time.Time{})
	if q.Date != emptyTamestamp || (q.DateFrom != emptyTamestamp || q.DateTo != emptyTamestamp) {
		if q.Date != emptyTamestamp {
			f["date"] = bson.M{"$eq": q.Date}
		} else {
			dateFilter := bson.M{}

			if q.DateFrom != emptyTamestamp {
				dateFilter["$gte"] = q.DateFrom
			}

			if q.DateTo != emptyTamestamp {
				dateFilter["$lt"] = q.DateTo
			}

			f["date"] = dateFilter
		}
	}

	return f, nil
}
