package usecase

import (
	"context"
	"log"

	mongoRepo "github.com/Borislavv/remote-executer/internal/data/mongo"
	agg "github.com/Borislavv/remote-executer/internal/domain/agg/msg"
	"github.com/Borislavv/remote-executer/internal/domain/builder"
)

type Polling struct {
	// offset for request messages from telegram
	msgOffset int64

	ctx     context.Context
	gateway *Telegram
	msgRepo *mongoRepo.MsgRepo
}

func NewPolling(
	ctx context.Context,
	gateway *Telegram,
	msgRepo *mongoRepo.MsgRepo,
) *Polling {
	return &Polling{
		ctx:     ctx,
		gateway: gateway,
		msgRepo: msgRepo,
	}
}

func (p *Polling) Do(messagesCh chan<- []agg.Msg, errCh chan<- error) {
	log.Println("polling has been started")

	for {
		select {
		case <-p.ctx.Done():
			log.Println("stop polling due to context signal")
			return
		default:
			offset, err := p.getOffset()
			if err != nil {
				errCh <- err
				continue
			}

			msgDTOs, err := p.gateway.GetMessages(offset)
			if err != nil {
				errCh <- err
				continue
			}

			msgAggs := builder.BuildMsgAggs(msgDTOs)

			len := len(msgAggs)
			if len > 0 {
				messagesCh <- msgAggs

				p.updateOffset(len)
			}
		}
	}
}

func (p *Polling) getOffset() (int64, error) {
	if p.msgOffset != 0 {
		return p.msgOffset, nil
	}

	offset, err := p.msgRepo.GetOffset(p.ctx)
	if err != nil {
		return 0, err
	}

	p.msgOffset = offset

	return offset, nil
}

func (p *Polling) updateOffset(msgsSliceLen int) {
	if p.msgOffset != 0 {
		p.msgOffset = p.msgOffset + int64(msgsSliceLen)
	}
}
