package mongoRepo

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Borislavv/remote-executer/internal/domain/agg"
	"github.com/Borislavv/remote-executer/internal/domain/errs"
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
		return errs.New(err).Interrupt()
	}
	if len(result.InsertedIDs) == 0 {
		return errs.New("MsgRepo.InsertMany: no one `Msg` document was created")
	}

	r.mu.Unlock()

	return nil
}

func (r *MsgRepo) Find(ctx context.Context, q agg.MsgQuery) ([]agg.Msg, error) {
	response := []agg.Msg{}

	filter, err := r.GetFilter(q)
	if err != nil {
		return response, errs.New(err).Interrupt()
	}

	cursor, err := r.coll.Find(ctx, filter, GetOpts(q))
	if err != nil {
		return response, errs.New(err).Interrupt()
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &response); err != nil {
		return response, errs.New(err).Interrupt()
	}

	return response, nil
}

func (r *MsgRepo) MarkAsExecuted(ctx context.Context, msg agg.Msg) error {
	_, err := r.coll.UpdateByID(ctx, msg.ID, bson.M{"$set": bson.M{"executed": true}})
	if err != nil {
		return errs.New(err).Interrupt()
	}
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
		return 0, errs.New(err).Interrupt()
	}

	c, err := r.coll.Find(ctx, f, GetOpts(q))
	if err != nil {
		return 0, errs.New(err).Interrupt()
	}
	defer c.Close(ctx)

	if !c.Next(ctx) {
		// handle the case when doc. was not found
		return 0, nil
	}
	if err := c.Decode(&resp); err != nil {
		return 0, errs.New(err).Interrupt()
	}

	return resp.Msg.UpdateId + 1, nil
}

func (r *MsgRepo) GetFilter(q agg.MsgQuery) (bson.M, error) {
	f := bson.M{}

	if q.ID != "" {
		id, err := primitive.ObjectIDFromHex(q.ID)
		if err != nil {
			return f, errs.New(err).Interrupt()
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
