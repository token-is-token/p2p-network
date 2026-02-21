package pubsub

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

type HandlerRegistry struct {
	mu       sync.RWMutex
	handlers map[string]MessageHandler
}

func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make(map[string]MessageHandler),
	}
}

func (r *HandlerRegistry) Register(topic string, handler MessageHandler) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.handlers[topic]; exists {
		return ErrHandlerAlreadyRegistered
	}

	r.handlers[topic] = handler
	return nil
}

func (r *HandlerRegistry) Unregister(topic string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.handlers, topic)
}

func (r *HandlerRegistry) GetHandler(topic string) (MessageHandler, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	h, ok := r.handlers[topic]
	return h, ok
}

func (r *HandlerRegistry) ListTopics() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	topics := make([]string, 0, len(r.handlers))
	for t := range r.handlers {
		topics = append(topics, t)
	}

	return topics
}

var (
	ErrHandlerAlreadyRegistered = &HandlerError{"handler already registered for topic"}
	ErrInvalidMessage           = &HandlerError{"invalid message"}
)

type HandlerError struct {
	msg string
}

func (e *HandlerError) Error() string {
	return e.msg
}

func NewProviderHandler() MessageHandler {
	return func(ctx context.Context, msg *Message) error {
		var provider ProviderMessage
		if err := json.Unmarshal(msg.Data, &provider); err != nil {
			return err
		}

		provider.ReceivedAt = time.Now()
		provider.From = string(msg.From)

		return nil
	}
}

func NewRequestHandler() MessageHandler {
	return func(ctx context.Context, msg *Message) error {
		var request RequestMessage
		if err := json.Unmarshal(msg.Data, &request); err != nil {
			return err
		}

		request.ReceivedAt = time.Now()
		request.From = string(msg.From)

		return nil
	}
}

func NewResponseHandler() MessageHandler {
	return func(ctx context.Context, msg *Message) error {
		var response ResponseMessage
		if err := json.Unmarshal(msg.Data, &response); err != nil {
			return err
		}

		response.ReceivedAt = time.Now()
		response.From = string(msg.From)

		return nil
	}
}

func NewHeartbeatHandler() MessageHandler {
	return func(ctx context.Context, msg *Message) error {
		var hb HeartbeatMessage
		if err := json.Unmarshal(msg.Data, &hb); err != nil {
			return err
		}

		hb.ReceivedAt = time.Now()

		return nil
	}
}

type ProviderMessage struct {
	Type      string    `json:"type"`
	Model     string    `json:"model"`
	Address   string    `json:"address"`
	Port      int       `json:"port"`
	Protocols []string  `json:"protocols"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	From      string    `json:"-"`
	ReceivedAt time.Time `json:"-"`
}

type RequestMessage struct {
	Type      string    `json:"type"`
	RequestID string    `json:"request_id"`
	Model     string    `json:"model"`
	Payload   []byte    `json:"payload"`
	From      string    `json:"-"`
	ReceivedAt time.Time `json:"-"`
}

type ResponseMessage struct {
	Type      string    `json:"type"`
	RequestID string    `json:"request_id"`
	Payload   []byte    `json:"payload"`
	Error     string    `json:"error,omitempty"`
	From      string    `json:"-"`
	ReceivedAt time.Time `json:"-"`
}

type HeartbeatMessage struct {
	Type      string    `json:"type"`
	PeerID    string    `json:"peer_id"`
	Timestamp int64     `json:"timestamp"`
	Address   string    `json:"address,omitempty"`
	ReceivedAt time.Time `json:"-"`
}
