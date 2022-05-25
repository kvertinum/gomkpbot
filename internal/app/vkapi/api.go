package vkapi

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const (
	baseApiUrl     = "https://api.vk.com/method/"
	baseApiVersion = "5.131"
)

type Api struct {
	Token   string
	Url     string
	Version string
}

func NewApi(token string) *Api {
	return &Api{
		Token:   token,
		Url:     baseApiUrl,
		Version: baseApiVersion,
	}
}

func (api *Api) Method(methodName string, params map[string]interface{}, response interface{}) error {
	params["access_token"] = api.Token
	params["v"] = api.Version

	urlParams := url.Values{}
	for k, v := range params {
		var strValue string
		switch t := v.(type) {
		case string:
			strValue = v.(string)
		case int:
			strValue = strconv.Itoa(v.(int))
		default:
			byteValue, err := json.Marshal(v)
			if err != nil {
				log.Printf("Unknown %v type\n", t)
				return err
			}
			strValue = string(byteValue)
		}
		urlParams.Add(k, strValue)
	}

	apiAnswer, err := http.PostForm(api.Url+methodName, urlParams)
	if err != nil {
		return err
	}
	defer apiAnswer.Body.Close()

	if response == nil {
		return nil
	}

	return json.NewDecoder(apiAnswer.Body).Decode(response)
}
