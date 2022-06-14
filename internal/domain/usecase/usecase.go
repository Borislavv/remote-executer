package usecase

import (
	"context"

	telegramGateway "github.com/Borislavv/remote-executer/pkg/gateway/telegram"
)

type PollingUseCase interface {
	NewPolling(ctx context.Context, gateway telegramGateway.TelegramGateway) *Polling
	Do(messagesCh chan<- telegramGateway.ResponseGetMessagesInterface, errCh chan<- error)
}

type MessagesUseCase interface {
	Consume(ctx context.Context, messagesCh <-chan telegramGateway.ResponseGetMessagesInterface, errCh chan<- error)
}

type CommandUseCase interface {
	Exec(ctx context.Context, responseCh chan<- telegramGateway.RequestGetMessagesInterface, errCh chan<- error)
}

type ResponseUseCase interface {
	Send(ctx context.Context, responseCh <-chan telegramGateway.RequestGetMessagesInterface, errCh chan<- error)
}
