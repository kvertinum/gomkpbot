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
	Key    string
	Server string
	Ts     string
	Wait   string
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
	return &Longpoll{
		Key:    r.Response.Key,
		Server: r.Response.Server,
		Ts:     r.Response.Ts,
		Wait:   strWait,
	}, nil
}

func (lp *Longpoll) Request() (*LongpollResponse, error) {
	urlParams := url.Values{}
	urlParams.Add("act", "a_check")
	urlParams.Add("ts", lp.Ts)
	urlParams.Add("key", lp.Key)
	urlParams.Add("wait", lp.Wait)

	response, err := http.PostForm(lp.Server, urlParams)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	r := &LongpollResponse{}
	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		return nil, err
	}

	lp.Ts = r.Ts

	return r, nil
}

func (lp *Longpoll) ListenNewMessages(message chan Message) {
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
				message <- m.CurrentMessage
			}
		}
	}
}
