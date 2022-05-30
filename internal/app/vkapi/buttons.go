package vkapi

import "encoding/json"

type Keyboard struct {
	OneTime     bool            `json:"one_time"`
	Buttons     [][]interface{} `json:"buttons"`
	Inline      bool            `json:"inline"`
	LastButtons []interface{}   `json:"-"`
}

func NewKeyboard(oneTime bool, inline bool) *Keyboard {
	return &Keyboard{
		OneTime:     oneTime,
		Inline:      inline,
		Buttons:     make([][]interface{}, 0),
		LastButtons: make([]interface{}, 0),
	}
}

func (k *Keyboard) Add(button interface{}) {
	k.LastButtons = append(k.LastButtons, button)
}

func (k *Keyboard) NewLine() {
	k.Buttons = append(k.Buttons, k.LastButtons)
	k.LastButtons = make([]interface{}, 0)
}

func (k *Keyboard) GetJson() (string, error) {
	res, err := json.Marshal(k)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

type TextButton struct {
	Action *TextButtonAction `json:"action"`
	Color  string            `json:"color"`
}

type TextButtonAction struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
	Label   string `json:"label"`
}

func NewTextButton(label string, payload string, color string) *TextButton {
	return &TextButton{
		Action: &TextButtonAction{
			Type:    "text",
			Payload: payload,
			Label:   label,
		},
		Color: color,
	}
}

func NewCallbackButton(label string, payload string, color string) *TextButton {
	return &TextButton{
		Action: &TextButtonAction{
			Type:    "callback",
			Payload: payload,
			Label:   label,
		},
		Color: color,
	}
}
