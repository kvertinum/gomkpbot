package vkapi

import (
	"encoding/json"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	baseApiUrl     = "https://api.vk.com/method/"
	baseApiVersion = "5.131"
)

type Api struct {
	Token   string
	Url     string
	Version string
	Client  *fasthttp.Client
}

func NewApi(token string) *Api {
	return &Api{
		Token:   token,
		Url:     baseApiUrl,
		Version: baseApiVersion,
		Client: &fasthttp.Client{
			ReadTimeout:              5 * time.Second,
			WriteTimeout:             5 * time.Second,
			MaxIdleConnDuration:      time.Minute,
			NoDefaultUserAgentHeader: true,
		},
	}
}

func (api *Api) Method(methodName string, params map[string]interface{}, response interface{}) error {
	params["access_token"] = api.Token
	params["v"] = api.Version

	urlParams := url.Values{}
	for key, value := range params {
		var strValue string
		switch t := value.(type) {
		case string:
			strValue = value.(string)
		case int:
			strValue = strconv.Itoa(value.(int))
		default:
			byteValue, err := json.Marshal(value)
			if err != nil {
				log.Printf("Bad type %v\n", t)
				return err
			}
			strValue = string(byteValue)
		}
		urlParams.Add(key, strValue)
	}

	urlEncoded := urlParams.Encode()
	reqEntityBytes := []byte(urlEncoded)

	return api.Post(
		api.Url+methodName, reqEntityBytes, response,
	)
}

func (api *Api) Post(url string, params []byte, response interface{}) error {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/x-www-form-urlencoded")
	req.SetRequestURI(url)
	req.SetBody(params)

	api.Client.Do(req, resp)
	body := resp.Body()

	if response == nil || len(body) == 0 {
		return nil
	}

	return json.Unmarshal(body, response)
}
