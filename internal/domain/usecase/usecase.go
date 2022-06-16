package usecase

import (
	"context"

	mongoRepo "github.com/Borislavv/remote-executer/internal/data/mongo"
	"github.com/Borislavv/remote-executer/internal/domain/agg"
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
	Consuming(messagesCh <-chan []agg.Msg, errCh chan<- error)
}

type CommandsUseCase interface {
	NewCommands(ctx context.Context, msgRepo *mongoRepo.MsgRepo)
	// Exec - find and execute commands
	Executing(responseCh chan<- dto.TelegramResponseInterface, errCh chan<- error)
}

type ResponsesUseCase interface {
	NewResponse(ctx context.Context, gateway *Telegram) *Responses
	// Send - sending response of executed commands
	Sending(responseCh <-chan dto.TelegramResponseInterface, errCh chan<- error)
}
