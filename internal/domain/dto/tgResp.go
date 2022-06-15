package dto

const ParseMode = "Markdown"

type TelegramResponseInterface interface {
	GetChatId() int64
	GetText() string
	GetParseMode() string
}

type TelegramResponse struct {
	chatId    int64
	text      string
	parseMode string
}

func NewTelegramResponse(chatId int64, text string) TelegramResponse {
	return TelegramResponse{
		chatId:    chatId,
		text:      text,
		parseMode: ParseMode,
	}
}

func (r TelegramResponse) GetChatId() int64 {
	return r.chatId
}

func (r TelegramResponse) GetText() string {
	return r.text
}

func (r TelegramResponse) GetParseMode() string {
	return r.parseMode
}
