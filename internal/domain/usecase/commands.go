package usecase

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"sync"
	"time"

	mongoRepo "github.com/Borislavv/remote-executer/internal/data/mongo"
	"github.com/Borislavv/remote-executer/internal/domain/agg"
	"github.com/Borislavv/remote-executer/internal/domain/dto"
	"github.com/Borislavv/remote-executer/internal/domain/errs"
)

type Commands struct {
	timeout int

	ctx     context.Context
	msgRepo *mongoRepo.MsgRepo
	wg      *sync.WaitGroup
}

func NewCommands(
	ctx context.Context,
	msgRepo *mongoRepo.MsgRepo,
	wg *sync.WaitGroup,
	mongoTimeout int,
) *Commands {
	return &Commands{
		ctx:     ctx,
		msgRepo: msgRepo,
		wg:      wg,
		timeout: mongoTimeout,
	}
}

func (c *Commands) Executing(responseCh chan<- dto.TelegramResponseInterface, errCh chan<- error) {
	log.Println("STARTED: executing of commands")
	defer c.wg.Done()

	for {
		select {
		case <-c.ctx.Done():
			log.Println("STOPPED: commands executing")
			return
		default:
			msgAggs, err := c.findCommandsForExecute()
			if err != nil {
				errCh <- errs.New(err).Interrupt()
				continue
			}

			for _, msgAgg := range msgAggs {
				resp, err := c.exec(msgAgg)
				if err != nil {
					errCh <- err
					responseCh <- resp
					continue
				}

				responseCh <- resp

				log.Printf("executed command: %s\n", msgAgg.Msg.Text)
			}

			// timeout before new request
			time.Sleep(time.Microsecond * time.Duration(c.timeout))
		}
	}
}

func (c *Commands) exec(msg agg.Msg) (dto.TelegramResponse, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("bash", "-c", msg.Msg.Text)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if err := c.markAsExecuted(msg); err != nil {
			return c.formatResponse(msg, stdout, stderr, err, false), errs.New(err).Interrupt()
		}

		return c.formatResponse(msg, stdout, stderr, err, true), nil
	}

	if err := c.markAsExecuted(msg); err != nil {
		return c.formatResponse(msg, stdout, stderr, err, false), errs.New(err).Interrupt()
	}

	return c.formatResponse(msg, stdout, stderr, nil, false), nil
}

func (c *Commands) formatResponse(
	msg agg.Msg,
	stdout bytes.Buffer,
	stderr bytes.Buffer,
	err error,
	rawErr bool,
) dto.TelegramResponse {
	resp := ""

	if stdout.String() == "" {
		if stderr.String() == "" {
			if err != nil && rawErr {
				resp = "Err: ```" + err.Error() + "```"
			} else {
				resp = fmt.Sprintf("Sorry, we can't execute this command: [%s].", msg.Msg.Text)
			}
		} else {
			resp = fmt.Sprintf("Err: ```%s```", stderr.String())
		}
	} else {
		resp = fmt.Sprintf("Out: ```%s```", stdout.String())
	}

	return dto.NewTelegramResponse(msg.Chat.Id, resp)
}

func (c *Commands) findCommandsForExecute() ([]agg.Msg, error) {
	return c.msgRepo.Find(c.ctx, agg.MsgQuery{
		ByExecuted: true,
		Executed:   false,
		OrderBy:    "updateId",
		SortBy:     "asc",
		Limit:      1,
	})
}

func (c *Commands) markAsExecuted(msg agg.Msg) error {
	if err := c.msgRepo.MarkAsExecuted(c.ctx, msg); err != nil {
		return errs.New(err).Interrupt()
	}
	return nil
}
