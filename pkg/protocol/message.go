package protocol

import (
	"encoding/json"
	"time"
)

type Request struct {
	Method    string          `json:"method"`
	Model     string          `json:"model"`
	Params    json.RawMessage `json:"params,omitempty"`
	ID        string          `json:"id"`
	Timestamp int64           `json:"timestamp"`
}

type Response struct {
	ID        string          `json:"id"`
	Result    json.RawMessage `json:"result,omitempty"`
	Error     *Error          `json:"error,omitempty"`
	Timestamp int64           `json:"timestamp"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewRequest(method, model, id string, params interface{}) (*Request, error) {
	var paramsData json.RawMessage
	if params != nil {
		data, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		paramsData = data
	}

	return &Request{
		Method:    method,
		Model:     model,
		Params:    paramsData,
		ID:        id,
		Timestamp: time.Now().Unix(),
	}, nil
}

func NewResponse(id string, result interface{}) (*Response, error) {
	var resultData json.RawMessage
	if result != nil {
		data, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}
		resultData = data
	}

	return &Response{
		ID:        id,
		Result:    resultData,
		Timestamp: time.Now().Unix(),
	}, nil
}

func NewErrorResponse(id string, code int, message string) *Response {
	return &Response{
		ID:        id,
		Error:     &Error{Code: code, Message: message},
		Timestamp: time.Now().Unix(),
	}
}

type Heartbeat struct {
	PeerID    string `json:"peer_id"`
	Timestamp int64  `json:"timestamp"`
	Address   string `json:"address,omitempty"`
}

func NewHeartbeat(peerID, address string) *Heartbeat {
	return &Heartbeat{
		PeerID:    peerID,
		Timestamp: time.Now().Unix(),
		Address:   address,
	}
}

type ProviderInfo struct {
	PeerID     string                 `json:"peer_id"`
	Model      string                 `json:"model"`
	Address    string                 `json:"address"`
	Protocols  []string               `json:"protocols"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	LastSeen   int64                  `json:"last_seen"`
}

func NewProviderInfo(peerID, model, address string) *ProviderInfo {
	return &ProviderInfo{
		PeerID:    peerID,
		Model:     model,
		Address:   address,
		LastSeen:  time.Now().Unix(),
	}
}
