package dto

type TelegramRequestInterface interface {
	GetOffset() int64
}

type TelegramRequest struct {
	offset int64
}

func NewTelegramRequest(offset int64) TelegramRequest {
	return TelegramRequest{
		offset: offset,
	}
}

func (r TelegramRequest) GetOffset() int64 {
	return r.offset
}
