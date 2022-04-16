package telegramGateway

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	getMessagesMethod = "getUpdates"
	sendMessageMethod = "sendMessage"
)

type TelegramGateway struct {
	endpointPattern string
	token           string
}

// NewGateway is a constructor of TelegramGateway struct.
func NewGateway(endpointPattern string, token string) *TelegramGateway {
	return &TelegramGateway{
		endpointPattern: endpointPattern,
		token:           token,
	}
}

func (gateway *TelegramGateway) GetMessages(reqMsgs RequestGetMessagesInterface) (ResponseGetMessagesInterface, error) {
	// getting message by network from telegram api
	response, err := http.Get(
		fmt.Sprintf(
			fmt.Sprintf(gateway.endpointPattern, gateway.token, getMessagesMethod),
			fmt.Sprint(reqMsgs.GetOffset()),
		),
	)
	if err != nil {
		return nil, err
	}

	// reading the received body to slice of bytes
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// unmarshaling slice of bytes to structure
	messages := NewResponseGetMessages()
	if err := json.Unmarshal(body, messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func (gateway *TelegramGateway) SendMessage(reqMsg RequestSendMessageInterface) error {
	// sending message by network to telegram api
	response, err := http.Post(
		fmt.Sprintf(
			fmt.Sprintf(
				gateway.endpointPattern, gateway.token, sendMessageMethod,
			),
			reqMsg.GetChatId(),
			url.QueryEscape(reqMsg.GetMessage()),
		),
		"application/json",
		strings.NewReader(url.Values{}.Encode()),
	)
	if err != nil {
		return err
	}

	// reading the received body to slice of bytes
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// unmarshaling slice of bytes to structure
	tgResponse := NewResponseSendMessage()
	if err := json.Unmarshal(body, tgResponse); err != nil {
		return err
	}
	if !tgResponse.IsOK() {
		return tgResponse.ToError()
	}

	return nil
}
