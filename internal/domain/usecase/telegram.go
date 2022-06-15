package usecase

import (
	"time"

	telegramGateway "github.com/Borislavv/TelegramApiGateway/telegram"

	"github.com/Borislavv/remote-executer/internal/domain/dto"
	"github.com/Borislavv/remote-executer/internal/domain/entity"
)

// Adapter of the Telegram gateway
// 	Repo: github.com/Borislavv/TelegramApiGateway
type Telegram struct {
	gateway *telegramGateway.TelegramGateway
}

func NewTelegram(endpoint string, token string) *Telegram {
	return &Telegram{
		gateway: telegramGateway.NewGateway(endpoint, token),
	}
}

func (t *Telegram) GetMessages(offset int64) ([]dto.Msg, error) {
	var msgDTOs []dto.Msg

	resp, err := t.gateway.GetMessages(telegramGateway.NewRequestGetMessages(offset))
	if err != nil {
		return msgDTOs, err
	}

	for _, respMsg := range resp.GetMessages() {
		msgDTOs = append(msgDTOs, dto.Msg{
			Text:     respMsg.GetText(),
			UpdateId: respMsg.GetQueueId(),
			Date:     time.Unix(respMsg.GetDate(), 0),
			User: entity.User{
				Firstname: respMsg.GetFirstName(),
				Lastname:  respMsg.GetLastName(),
				Username:  respMsg.GetUsername(),
			},
			Chat: entity.Chat{
				Id:   respMsg.GetChatId(),
				Type: respMsg.GetChatType(),
			},
		})
	}

	return msgDTOs, nil
}

func (t *Telegram) SendMessage(chatId int64, text string) error {
	return t.gateway.SendMessage(telegramGateway.NewRequestSendMessage(chatId, text))
}
