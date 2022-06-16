package usecase

import (
	"context"
	"log"
	"sync"
	"time"

	mongoRepo "github.com/Borislavv/remote-executer/internal/data/mongo"
	"github.com/Borislavv/remote-executer/internal/domain/agg"
	"github.com/Borislavv/remote-executer/internal/domain/builder"
	"github.com/Borislavv/remote-executer/internal/domain/errs"
)

type Polling struct {
	// offset for request messages from telegram
	msgOffset int64
	timeout   int

	ctx     context.Context
	gateway *Telegram
	msgRepo *mongoRepo.MsgRepo
	wg      *sync.WaitGroup
}

func NewPolling(
	ctx context.Context,
	gateway *Telegram,
	msgRepo *mongoRepo.MsgRepo,
	wg *sync.WaitGroup,
	pollingTimeout int,
) *Polling {
	return &Polling{
		ctx:     ctx,
		gateway: gateway,
		msgRepo: msgRepo,
		wg:      wg,
		timeout: pollingTimeout,
	}
}

func (p *Polling) Do(messagesCh chan<- []agg.Msg, errCh chan<- error) {
	log.Println("STARTED: polling")
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			log.Println("STOPPED: polling")
			return
		default:
			offset, err := p.getOffset()
			if err != nil {
				errCh <- errs.New(err).Interrupt()
				continue
			}

			msgDTOs, err := p.gateway.GetMessages(offset)
			if err != nil {
				errCh <- errs.New(err).Interrupt()
				continue
			}

			msgAggs := builder.BuildMsgAggs(msgDTOs)

			len := len(msgAggs)
			if len > 0 {
				messagesCh <- msgAggs

				p.updateOffset(len)
			}

			// timeout 0.25 before new request
			time.Sleep(time.Millisecond * time.Duration(p.timeout))
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
