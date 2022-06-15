package entity

type Chat struct {
	// Id of the telegram chat
	//
	// required: true
	// example: 1063099947
	Id int64 `json:"id" bson:"id"`

	// Type of the Chat
	//
	// required: true
	// example: `private`
	Type string `json:"type" bson:"type"`
}
