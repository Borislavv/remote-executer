package builder

import (
	"time"

	agg "github.com/Borislavv/remote-executer/internal/domain/agg/msg"
	"github.com/Borislavv/remote-executer/internal/domain/dto"
	"github.com/Borislavv/remote-executer/internal/domain/entity"
	"github.com/Borislavv/remote-executer/internal/domain/vo"
)

func BuildMsgAggs(msgDTOs []dto.Msg) []agg.Msg {
	var msgAggs []agg.Msg

	for _, msg := range msgDTOs {
		msgAggs = append(msgAggs, buildMsgAgg(msg))
	}

	return msgAggs
}

func buildMsgAgg(msgDTO dto.Msg) agg.Msg {
	return agg.Msg{
		Msg: entity.Msg{
			Text:     msgDTO.Text,
			UpdateId: msgDTO.UpdateId,
			Date:     msgDTO.Date,
		},
		Executed: false,
		User:     msgDTO.User,
		Chat:     msgDTO.Chat,
		Timestamp: vo.Timestamp{
			CreatedAt: time.Now(),
		},
	}
}
