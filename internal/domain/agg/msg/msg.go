package agg

import (
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Borislavv/remote-executer/internal/domain/entity"
	"github.com/Borislavv/remote-executer/internal/domain/vo"
)

type Msg struct {
	// key
	ID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`

	// data
	Msg entity.Msg `bson:",inline"`

	// referencies
	User entity.User
	Chat entity.Chat

	// Volume object: created at...
	Timestamp vo.Timestamp `bson:",inline"`
}
