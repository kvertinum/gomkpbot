package vkapi

type ResponseInit struct {
	Response struct {
		Key    string `json:"key"`
		Server string `json:"server"`
		Ts     string `json:"ts"`
	} `json:"response"`
}

type Message struct {
	Date                  int           `json:"date"`
	FromID                int           `json:"from_id"`
	ID                    int           `json:"id"`
	Out                   int           `json:"out"`
	PeerID                int           `json:"peer_id"`
	Text                  string        `json:"text"`
	ConversationMessageID int           `json:"conversation_message_id"`
	FwdMessages           []*Message    `json:"fwd_messages"`
	ReplyMessage          *Message      `json:"reply_message"`
	Important             bool          `json:"important"`
	RandomID              int           `json:"random_id"`
	Attachments           []interface{} `json:"attachments"`
	IsHidden              bool          `json:"is_hidden"`
}

type MessageJson struct {
	CurrentMessage Message `json:"message"`
}

type LongpollMessage struct {
	Type    string                 `json:"type"`
	Object  map[string]interface{} `json:"object"`
	GroupID int                    `json:"group_id"`
}

type LongpollResponse struct {
	Ts      string            `json:"ts"`
	Updates []LongpollMessage `json:"updates"`
}
