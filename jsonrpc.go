package wsjsonrpc

import (
	"encoding/json"

	"golang.org/x/net/websocket"
)

type JsonRPC struct {
	version string
	url     string
	origin  string
	conn    *websocket.Conn
	codec   *websocket.Codec
	handler map[string]RecvHandler
}

type JsonRPCMessage struct {
	Version string          `json:"version"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      *int            `json:"id,omitempty"`
}

type RecvHandler func(msg json.RawMessage, id *int)

func NewJsonRPC(version string, url, origin string) (*JsonRPC, error) {
	c := &websocket.Codec{
		Marshal: func(v interface{}) ([]byte, byte, error) {
			b, err := json.Marshal(v)
			return b, websocket.TextFrame, err
		},
		Unmarshal: func(data []byte, payloadType byte, v interface{}) error {
			return json.Unmarshal(data, v)
		},
	}

	conn, err := websocket.Dial(
		url,
		"",
		origin)
	if err != nil {
		return nil, err
	}

	return &JsonRPC{version: version, url: url, origin: origin, conn: conn, codec: c, handler: map[string]RecvHandler{}}, nil
}

func (j *JsonRPC) Send(method string, msg interface{}, id *int) error {
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return j.codec.Send(j.conn, &JsonRPCMessage{
		Version: j.version,
		Method:  method,
		Params:  json.RawMessage(b),
		ID:      id,
	})
}

func (j *JsonRPC) Recv() (string, interface{}, *int, error) {
	msg := JsonRPCMessage{}
	if err := j.codec.Receive(j.conn, &msg); err != nil {
		return "", nil, nil, err
	}

	return msg.Method, msg.Params, msg.ID, nil
}

func (j *JsonRPC) Close() error {
	return j.conn.Close()
}
