<h1 align="center">
  MKP bot
</h1>

> Fast golang vk bot

## Example
```golang
package main

import "github.com/Kvertinum01/gomkpbot/internal/app/vkapi"

func main() {
  api := vkapi.NewApi("Token here")
  lp, err := vkapi.NewLongpoll(api, "int group id here")
  if err != nil {
    panic(err)
  }

  go lp.ListenNewMessages()

  for message := range lp.NewMessage {
    switch message.Text {
    case "hello":
      if err := api.Method("messages.send", map[string]interface{
        "peer_id":   message.PeerID,
        "random_id": 0,
        "message":   "hi!",
      }, nil); err != nil {
        panic(err)
      }
    }
  }
}
```