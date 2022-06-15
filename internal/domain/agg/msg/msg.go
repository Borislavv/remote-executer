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
	Msg      entity.Msg `bson:",inline"`
	Executed bool       `json:"executed" bson:"executed"`

	// referencies
	User entity.User `bson:",inline"`
	Chat entity.Chat `bson:",inline"`

	// Volume object: created at...
	Timestamp vo.Timestamp `bson:",inline"`
}
