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
	Payload               string        `json:"payload"`
	Text                  string        `json:"text"`
	ConversationMessageID int           `json:"conversation_message_id"`
	FwdMessages           []*Message    `json:"fwd_messages"`
	ReplyMessage          *Message      `json:"reply_message"`
	Important             bool          `json:"important"`
	RandomID              int           `json:"random_id"`
	Attachments           []Attachments `json:"attachments"`
	IsHidden              bool          `json:"is_hidden"`
	Action                *struct {
		MemberID int    `json:"member_id"`
		Type     string `json:"type"`
	} `json:"action"`
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

type Attachments struct {
	Photo struct {
		AlbumID int  `json:"album_id"`
		Date    int  `json:"date"`
		HasTags bool `json:"has_tags"`
		ID      int  `json:"id"`
		OwnerID int  `json:"owner_id"`
		Sizes   []struct {
			Height int    `json:"height"`
			Type   string `json:"type"`
			Url    string `json:"url"`
		} `json:"sizes"`
		Text string `json:"text"`
	} `json:"photo"`
	Type string `json:"type"`
}
