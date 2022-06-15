package usecase

import (
	"context"
	"log"

	mongoRepo "github.com/Borislavv/remote-executer/internal/data/mongo"
	agg "github.com/Borislavv/remote-executer/internal/domain/agg/msg"
)

type Messages struct {
	ctx     context.Context
	msgRepo *mongoRepo.MsgRepo
}

func NewMessages(ctx context.Context, msgRepo *mongoRepo.MsgRepo) *Messages {
	return &Messages{
		ctx:     ctx,
		msgRepo: msgRepo,
	}
}

func (m *Messages) Consume(messagesCh <-chan []agg.Msg, errCh chan<- error) {
	log.Println("consuming messages has been started")

	for {
		select {
		case <-m.ctx.Done():
			log.Println("stop consuming messages due to context signal")
			return
		case msgAggs := <-messagesCh:
			if err := m.msgRepo.InsertMany(m.ctx, msgAggs); err != nil {
				errCh <- err
				continue
			}
		default:
			// don't block on awaiting messages from channels
		}
	}
}
