package vkapi

import (
	"encoding/json"
	"time"

	// "log"
	// "strconv"

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
			ReadTimeout:         5 * time.Second,
			WriteTimeout:        5 * time.Second,
			MaxIdleConnDuration: time.Minute,
		},
	}
}

func (api *Api) Method(methodName string, params map[string]interface{}, response interface{}) error {
	params["access_token"] = api.Token
	params["v"] = api.Version

	reqEntityBytes, err := json.Marshal(params)
	if err != nil {
		return err
	}
	reqUrl := api.Url + methodName

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetRequestURI(reqUrl)
	req.SetBody(reqEntityBytes)

	api.Client.Do(req, resp)

	return json.Unmarshal(resp.Body(), response)

	// urlParams := url.Values{}
	// strParams := "?"
	// for key, value := range params {
	// 	var strValue string
	// 	switch t := value.(type) {
	// 	case string:
	// 		strValue = value.(string)
	// 	case int:
	// 		strValue = strconv.Itoa(value.(int))
	// 	default:
	// 		byteValue, err := json.Marshal(value)
	// 		if err != nil {
	// 			log.Printf("Bad type %v\n", t)
	// 			return err
	// 		}
	// 		strValue = string(byteValue)
	// 	}
	// 	// urlParams.Add(key, strValue)
	// 	strParams += key + "=" + strValue + "&"
	// }

	//apiAnswer, err := http.PostForm(api.Url+methodName, urlParams)

	// Using GET method for better perfomance
	//apiAnswer, err := http.Get(api.Url + methodName + strParams)
	// if err != nil {
	// 	return err
	// }
	// defer apiAnswer.Body.Close()

	// if response == nil {
	// 	return nil
	// }

	// return json.NewDecoder(apiAnswer.Body).Decode(response)
}
