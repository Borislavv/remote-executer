package usecase

import (
	"context"
	"log"
	"sync"

	mongoRepo "github.com/Borislavv/remote-executer/internal/data/mongo"
	"github.com/Borislavv/remote-executer/internal/domain/agg"
)

type Messages struct {
	ctx     context.Context
	msgRepo *mongoRepo.MsgRepo
	wg      *sync.WaitGroup
}

func NewMessages(ctx context.Context, msgRepo *mongoRepo.MsgRepo, wg *sync.WaitGroup) *Messages {
	return &Messages{
		ctx:     ctx,
		msgRepo: msgRepo,
		wg:      wg,
	}
}

func (m *Messages) Consuming(messagesCh <-chan []agg.Msg, errCh chan<- error) {
	log.Println("STARTED: consuming messages")
	defer m.wg.Done()

	for {
		select {
		case <-m.ctx.Done():
			log.Println("STOPPED: consuming messages")
			return
		case msgAggs := <-messagesCh:
			if err := m.msgRepo.InsertMany(m.ctx, msgAggs); err != nil {
				errCh <- err
				continue
			}
		}
	}
}
