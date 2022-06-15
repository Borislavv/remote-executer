package usecase

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"

	mongoRepo "github.com/Borislavv/remote-executer/internal/data/mongo"
	agg "github.com/Borislavv/remote-executer/internal/domain/agg"
	"github.com/Borislavv/remote-executer/internal/domain/dto"
	"github.com/Borislavv/remote-executer/internal/util"
)

type Commands struct {
	ctx     context.Context
	msgRepo *mongoRepo.MsgRepo
}

func NewCommands(ctx context.Context, msgRepo *mongoRepo.MsgRepo) *Commands {
	return &Commands{
		ctx:     ctx,
		msgRepo: msgRepo,
	}
}

func (c *Commands) Executing(responseCh chan<- dto.TelegramResponseInterface, errCh chan<- error) {
	log.Println("executing of commands has been started")

OUTER:
	for {
		select {
		case <-c.ctx.Done():
			log.Println("stop commands executing due to context signal")
			return
		default:
			msgAggs, err := c.findCommandsForExecute()
			if err != nil {
				errCh <- util.ErrWithTrace(err)
				continue
			}

			for _, msgAgg := range msgAggs {
				resp, err := c.exec(msgAgg)
				if err != nil {
					errCh <- err
					responseCh <- resp
					continue OUTER
				}

				responseCh <- resp

				log.Printf("Executed command: %s\n", msgAgg.Msg.Text)
			}
		}
	}
}

func (c *Commands) exec(msg agg.Msg) (dto.TelegramResponse, error) {
	fmt.Println(msg)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("bash", "-c", msg.Msg.Text)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if err := c.markAsExecuted(msg); err != nil {
			return c.formatResponse(msg, stdout, stderr, err), util.ErrWithTrace(err)
		}

		return c.formatResponse(msg, stdout, stderr, err), util.ErrWithTrace(err)
	}

	if err := c.markAsExecuted(msg); err != nil {
		return c.formatResponse(msg, stdout, stderr, err), util.ErrWithTrace(err)
	}

	return c.formatResponse(msg, stdout, stderr, nil), nil
}

func (c *Commands) formatResponse(msg agg.Msg, stdout bytes.Buffer, stderr bytes.Buffer, err error) dto.TelegramResponse {
	resp := ""

	if err != nil {
		resp = fmt.Sprintf("Sorry, we can't execute this command: [%s].", msg.Msg.Text)
	} else {
		if stderr.String() != "" {
			resp = fmt.Sprintf("Err: ```%s```", stderr.String())
		} else {
			resp = fmt.Sprintf("Out: ```%s```", stdout.String())
		}
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
		return util.ErrWithTrace(err)
	}
	return nil
}
