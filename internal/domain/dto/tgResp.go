package dto

type TelegramResponseInterface interface {
	GetChatId() int64
	GetText() string
}

type TelegramResponse struct {
	chatId int64
	text   string
}

func NewTelegramResponse(chatId int64, text string) TelegramResponse {
	return TelegramResponse{
		chatId: chatId,
		text:   text,
	}
}

func (r TelegramResponse) GetChatId() int64 {
	return r.chatId
}

func (r TelegramResponse) GetText() string {
	return r.text
}
