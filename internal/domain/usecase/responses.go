package usecase

import (
	"context"
	"log"
	"sync"

	"github.com/Borislavv/remote-executer/internal/domain/dto"
	"github.com/Borislavv/remote-executer/internal/domain/errs"
)

type Responses struct {
	ctx     context.Context
	gateway *Telegram
	wg      *sync.WaitGroup
}

func NewResponses(ctx context.Context, gateway *Telegram, wg *sync.WaitGroup) *Responses {
	return &Responses{
		ctx:     ctx,
		gateway: gateway,
		wg:      wg,
	}
}

func (r *Responses) Sending(responseCh <-chan dto.TelegramResponseInterface, errCh chan<- error) {
	log.Println("STARTED: sending responses")
	defer r.wg.Done()

	for {
		select {
		case <-r.ctx.Done():
			log.Println("STOPPED: sending responses")
			return
		case resp := <-responseCh:
			if err := r.gateway.SendMessage(resp); err != nil {
				errCh <- errs.New(err).Interrupt()
				continue
			}
		}
	}
}
