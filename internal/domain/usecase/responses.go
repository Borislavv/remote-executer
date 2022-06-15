package usecase

import (
	"context"
	"log"

	"github.com/Borislavv/remote-executer/internal/domain/dto"
	"github.com/Borislavv/remote-executer/internal/util"
)

type Responses struct {
	ctx     context.Context
	gateway *Telegram
}

func NewResponses(ctx context.Context, gateway *Telegram) *Responses {
	return &Responses{
		ctx:     ctx,
		gateway: gateway,
	}
}

func (r *Responses) Sending(responseCh <-chan dto.TelegramResponseInterface, errCh chan<- error) {
	log.Println("sending responses has been started")

	for {
		select {
		case <-r.ctx.Done():
			log.Println("stop sending responses due to context signal")
		case resp := <-responseCh:
			if err := r.gateway.SendMessage(resp); err != nil {
				errCh <- util.ErrWithTrace(err)
				continue
			}
		default:
			// don't block on awaiting messages from channels
		}
	}
}
