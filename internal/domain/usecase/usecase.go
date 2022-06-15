package usecase

import (
	"context"

	mongoRepo "github.com/Borislavv/remote-executer/internal/data/mongo"
	agg "github.com/Borislavv/remote-executer/internal/domain/agg/msg"
	"github.com/Borislavv/remote-executer/internal/domain/dto"
)

type PollingUseCase interface {
	NewPolling(ctx context.Context, gateway *Telegram, msgRepo *mongoRepo.MsgRepo) *Polling
	// Do - polling telegram and send messages into channel
	Do(messagesCh chan<- []agg.Msg, errCh chan<- error)
}

type MessagesUseCase interface {
	NewMessages(ctx context.Context, msgRepo *mongoRepo.MsgRepo) *Messages
	// Consume - consuming messages from channel and store them
	Consume(messagesCh <-chan []agg.Msg, errCh chan<- error)
}

type CommandsUseCase interface {
	NewCommands(ctx context.Context, msgRepo *mongoRepo.MsgRepo)
	// Exec - find and execute commands
	Exec(responseCh chan<- dto.TelegramResponseInterface, errCh chan<- error)
}

type ResponseUseCase interface {
	NewResponse(ctx context.Context)
	// Send - sending response of executed commands
	Send(responseCh <-chan dto.TelegramResponseInterface, errCh chan<- error)
}
