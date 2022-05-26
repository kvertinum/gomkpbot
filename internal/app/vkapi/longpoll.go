package vkapi

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const (
	baseLongpollWait = 15
)

type Longpoll struct {
	Params      url.Values
	Server      string
	LastEvent   chan LongpollMessage
	LastMessage chan Message
}

func NewLongpoll(api *Api, groupID int) (*Longpoll, error) {
	r := ResponseInit{}
	err := api.Method("groups.getLongPollServer", map[string]interface{}{
		"group_id": groupID,
	}, &r)
	if err != nil {
		return nil, err
	}
	strWait := strconv.Itoa(baseLongpollWait)
	urlParams := url.Values{}
	urlParams.Add("act", "a_check")
	urlParams.Add("ts", r.Response.Ts)
	urlParams.Add("key", r.Response.Key)
	urlParams.Add("wait", strWait)
	return &Longpoll{
		Params:      urlParams,
		Server:      r.Response.Server,
		LastEvent:   make(chan LongpollMessage),
		LastMessage: make(chan Message),
	}, nil
}

func (lp *Longpoll) Request() (*LongpollResponse, error) {
	response, err := http.PostForm(lp.Server, lp.Params)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	r := &LongpollResponse{}
	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		return nil, err
	}

	lp.Params.Set("ts", r.Ts)

	return r, nil
}

func (lp *Longpoll) ListenNewMessages() {
	for {
		event, err := lp.Request()
		if err != nil {
			log.Fatal(err)
		}
		for _, update := range event.Updates {
			if update.Type == "message_new" {
				m := MessageJson{}
				jsonString, err := json.Marshal(update.Object)
				if err != nil {
					log.Fatal(err)
				}
				if err := json.Unmarshal(jsonString, &m); err != nil {
					log.Fatal(err)
				}
				lp.LastMessage <- m.CurrentMessage
			} else {
				lp.LastEvent <- update
			}
		}
	}
}
