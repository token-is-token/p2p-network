package protocol

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

type Handler struct {
	host      host.Host
	node      interface{}
	mu        sync.RWMutex
	handlers  map[MessageType]MessageHandler
	responses map[string]chan *Message
}

type MessageHandler func(ctx context.Context, p peer.ID, msg *Message) (*Message, error)

func NewHandler(node interface{}) *Handler {
	h := &Handler{
		node:      node,
		handlers:  make(map[MessageType]MessageHandler),
		responses: make(map[string]chan *Message),
	}

	h.registerDefaultHandlers()

	return h
}

func (h *Handler) registerDefaultHandlers() {
	h.handlers[MsgTypeRequest] = h.handleRequest
	h.handlers[MsgTypeResponse] = h.handleResponse
	h.handlers[MsgTypeHeartbeat] = h.handleHeartbeat
	h.handlers[MsgTypePing] = h.handlePing
	h.handlers[MsgTypePong] = h.handlePong
}

func (h *Handler) HandleStream(stream network.Stream) {
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	decoder := json.NewDecoder(stream)
	encoder := json.NewEncoder(stream)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		var msg Message
		if err := decoder.Decode(&msg); err != nil {
			if err == io.EOF {
				return
			}
			continue
		}

		handler, ok := h.handlers[msg.Type]
		if !ok {
			continue
		}

		resp, err := handler(ctx, stream.Conn().RemotePeer(), &msg)
		if err != nil {
			continue
		}

		if resp != nil {
			if err := encoder.Encode(resp); err != nil {
				continue
			}
		}
	}
}

func (h *Handler) handleRequest(ctx context.Context, p peer.ID, msg *Message) (*Message, error) {
	var req Request
	if err := json.Unmarshal(msg.Payload, &req); err != nil {
		return nil, err
	}

	resp := &Response{
		ID:        req.ID,
		Result:    nil,
		Timestamp: time.Now().Unix(),
	}

	respData, _ := json.Marshal(resp)

	return &Message{
		Type:      MsgTypeResponse,
		RequestID: msg.RequestID,
		Payload:   respData,
	}, nil
}

func (h *Handler) handleResponse(ctx context.Context, p peer.ID, msg *Message) (*Message, error) {
	h.mu.RLock()
	ch, ok := h.responses[msg.RequestID]
	h.mu.RUnlock()

	if ok {
		select {
		case ch <- msg:
		default:
		}
	}

	return nil, nil
}

func (h *Handler) handleHeartbeat(ctx context.Context, p peer.ID, msg *Message) (*Message, error) {
	return nil, nil
}

func (h *Handler) handlePing(ctx context.Context, p peer.ID, msg *Message) (*Message, error) {
	return &Message{
		Type:      MsgTypePong,
		RequestID: msg.RequestID,
	}, nil
}

func (h *Handler) handlePong(ctx context.Context, p peer.ID, msg *Message) (*Message, error) {
	return nil, nil
}

func (h *Handler) RegisterHandler(msgType MessageType, handler MessageHandler) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.handlers[msgType] = handler
}

func (h *Handler) SendRequest(ctx context.Context, p peer.ID, msg *Message) (*Message, error) {
	stream, err := h.host.NewStream(ctx, p, ProtocolID)
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	if err := json.NewEncoder(stream).Encode(msg); err != nil {
		return nil, err
	}

	var resp Message
	if err := json.NewDecoder(stream).Decode(&resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (h *Handler) SendMessage(ctx context.Context, p peer.ID, msg *Message) error {
	stream, err := h.host.NewStream(ctx, p, ProtocolID)
	if err != nil {
		return err
	}
	defer stream.Close()

	return json.NewEncoder(stream).Encode(msg)
}

func (h *Handler) Broadcast(ctx context.Context, peers []peer.ID, msg *Message) error {
	var wg sync.WaitGroup
	errs := make([]error, len(peers))

	for i, p := range peers {
		wg.AddGo(func() {
			errs[i] = h.SendMessage(ctx, p, msg)
		})
	}

	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *Handler) SetHost(host host.Host) {
	h.host = host
}
