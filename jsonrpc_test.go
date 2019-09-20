package wsjsonrpc

import (
	"encoding/json"
	"fmt"
)

func ExampleJsonRPC() {
	rpc, _ := NewJsonRPC("2.0", "wss://ws.lightstream.bitflyer.com/json-rpc", "https://ws.lightstream.bitflyer.com/json-rpc")

	type channel struct {
		Channel string `json:"channel"`
	}

	rpc.Open()

	rpc.OnRecv("channelMessage", func(msg json.RawMessage, id *int) {
		rpc.Close()
	})

	go rpc.Recv()

	rpc.Send("subscribe", &channel{Channel: "lightning_board_BTC_JPY"}, nil)

	fmt.Println("close")
	// Output: close
}
