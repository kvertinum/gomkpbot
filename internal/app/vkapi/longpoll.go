package vkapi

import (
	"encoding/json"
	"log"
	"net/url"
	"strconv"
)

const (
	baseLongpollWait = 15
)

type Longpoll struct {
	Api        *Api
	Params     url.Values
	Server     string
	NewEvent   chan LongpollMessage
	NewMessage chan Message
}

func NewLongpoll(api *Api, groupID int) (*Longpoll, error) {
	// Getting longpoll server and configs
	r := ResponseInit{}
	err := api.Method("groups.getLongPollServer", map[string]interface{}{
		"group_id": groupID,
	}, &r)
	if err != nil {
		return nil, err
	}

	// Convert LongpollWait value to string
	strWait := strconv.Itoa(baseLongpollWait)

	// Create url values
	urlParams := url.Values{}
	urlParams.Add("act", "a_check")
	urlParams.Add("ts", r.Response.Ts)
	urlParams.Add("key", r.Response.Key)
	urlParams.Add("wait", strWait)

	// Init Longpoll struct
	return &Longpoll{
		Api:        api,
		Params:     urlParams,
		Server:     r.Response.Server,
		NewEvent:   make(chan LongpollMessage),
		NewMessage: make(chan Message),
	}, nil
}

func (lp *Longpoll) Request() (*LongpollResponse, error) {
	// Create request to longpoll
	r := &LongpollResponse{}
	err := lp.Api.Post(lp.Server, []byte(lp.Params.Encode()), &r)
	if err != nil {
		return nil, err
	}

	if r.Ts == "" {
		return nil, nil
	}

	lp.Params.Set("ts", r.Ts)

	return r, nil
}

func (lp *Longpoll) ListenNewMessages() {
	// Listen new longpoll events
	for {
		// Create request
		event, err := lp.Request()
		if err != nil {
			log.Fatal(err)
		}
		if event == nil {
			continue
		}
		// Check updates
		for _, update := range event.Updates {
			if update.Type == "message_new" {
				// Init MessageJson struct
				m := MessageJson{}
				jsonString, err := json.Marshal(update.Object)
				if err != nil {
					log.Fatal(err)
				}
				if err := json.Unmarshal(jsonString, &m); err != nil {
					log.Fatal(err)
				}
				// Write message to lp.NewMessage
				lp.NewMessage <- m.CurrentMessage
			} else {
				// Write event to lp.NewEvent
				lp.NewEvent <- update
			}
		}
	}
}
