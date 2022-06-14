package usecase

import (
	"time"

	"github.com/Borislavv/remote-executer/internal/domain/dto"
	"github.com/Borislavv/remote-executer/internal/domain/entity"
	tg "github.com/Borislavv/remote-executer/pkg/gateway/telegram"
)

// Adapter of the Telegram gateway
// 	Repo: github.com/Borislavv/remote-executer/pkg/gateway/telegram
type Telegram struct {
	gateway *tg.TelegramGateway
}

func NewTelegram(endpoint string, token string) *Telegram {
	return &Telegram{
		gateway: tg.NewGateway(endpoint, token),
	}
}

func (t *Telegram) GetMessages(offset int64) ([]dto.Msg, error) {
	var msgDTOs []dto.Msg

	resp, err := t.gateway.GetMessages(tg.NewRequestGetMessages(offset))
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
	return t.gateway.SendMessage(tg.NewRequestSendMessage(chatId, text))
}
