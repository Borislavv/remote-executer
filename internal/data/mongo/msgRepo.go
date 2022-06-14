package mongoRepo

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	agg "github.com/Borislavv/remote-executer/internal/domain/agg/msg"
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
	if len(result.InsertedIDs) == 0 {
		return errors.New("MsgRepo.InsertMany: no one `Msg` document was created")
	}

	r.mu.Unlock()

	return err
}

func (r *MsgRepo) Find(ctx context.Context, q agg.MsgQuery) ([]agg.Msg, error) {
	response := []agg.Msg{}

	filter, err := r.GetFilter(q)
	if err != nil {
		return response, err
	}

	cursor, err := r.coll.Find(ctx, filter, GetOpts(q))
	if err != nil {
		return response, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &response); err != nil {
		return response, err
	}

	return response, nil
}

func (r *MsgRepo) GetOffset(ctx context.Context) (int64, error) {
	var msg agg.MsgQuery

	q := agg.MsgQuery{
		SortBy:  "updateId",
		OrderBy: "desc",
		Limit:   1,
	}

	f, err := r.GetFilter(q)
	if err != nil {
		return 0, err
	}

	c, err := r.coll.Find(ctx, f, GetOpts(q))
	if err != nil {
		return 0, err
	}
	defer c.Close(ctx)

	if err := c.Decode(msg); err != nil {
		return 0, err
	}

	return msg.UpdateId + 1, nil
}

func (r *MsgRepo) GetFilter(q agg.MsgQuery) (bson.M, error) {
	f := bson.M{}

	if q.ID != "" {
		id, err := primitive.ObjectIDFromHex(q.ID)
		if err != nil {
			return f, err
		}

		f["_id"] = bson.M{"$eq": id}
	}

	if q.Text != "" {
		f["text"] = bson.M{"$eq": q.Text}
	}

	if q.UpdateId != 0 {
		f["updateId"] = bson.M{"$eq": q.UpdateId}
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
