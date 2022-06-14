package entity

import "time"

type Msg struct {
	// Text of the message
	//
	// required: true
	// example: `Hello world`
	Text string `json:"text" bson:"text"`

	// UpdateId value, which reprecent telegram offset of the message
	//
	// required: true
	// example: 506233478
	UpdateId int64 `json:"updateId" bson:"updateId"`

	// Date of stat. from file
	//
	// pattern: `\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`
	// required: true
	// example: 2021-11-21 12:15:17
	Date time.Time `json:"date" bson:"date"`
}
