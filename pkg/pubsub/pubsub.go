package pubsub

import (
	"context"

	"github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p-pubsub/pb"
)

type MessageHandler func(ctx context.Context, msg *Message) error

type PubSubManager struct {
	pubsub     *pubsub.PubSub
	subs       map[string]*Subscription
	handlers   map[string]MessageHandler
}

func NewManager(ps *pubsub.PubSub) *PubSubManager {
	return &PubSubManager{
		pubsub:   ps,
		subs:     make(map[string]*Subscription),
		handlers: make(map[string]MessageHandler),
	}
}

func (m *PubSubManager) Subscribe(topic string, handler MessageHandler) (*Subscription, error) {
	if m.pubsub == nil {
		return nil, nil
	}

	sub, err := m.pubsub.Subscribe(topic)
	if err != nil {
		return nil, err
	}

	subscription := &Subscription{
		topic:    topic,
		sub:      sub,
		handler:  handler,
		messages: make(chan *Message, 100),
	}

	m.subs[topic] = subscription
	m.handlers[topic] = handler

	go subscription.readLoop()

	return subscription, nil
}

func (m *PubSubManager) Publish(topic string, data []byte) error {
	if m.pubsub == nil {
		return nil
	}

	return m.pubsub.Publish(topic, data)
}

func (m *PubSubManager) PublishWithOptions(topic string, data []byte, opts ...PublishOption) error {
	if m.pubsub == nil {
		return nil
	}

	msg := &pb.Message{
		Data: data,
	}

	for _, opt := range opts {
		opt(msg)
	}

	return m.pubsub.Publish(topic, msg)
}

func (m *PubSubManager) Unsubscribe(topic string) error {
	sub, ok := m.subs[topic]
	if !ok {
		return nil
	}

	delete(m.subs, topic)
	delete(m.handlers, topic)

	return sub.Cancel()
}

func (m *PubSubManager) GetTopics() []string {
	if m.pubsub == nil {
		return nil
	}

	return m.pubsub.GetTopics()
}

func (m *PubSubManager) ListPeers(topic string) []string {
	if m.pubsub == nil {
		return nil
	}

	return m.pubsub.ListPeers(topic)
}

func (m *PubSubManager) TopicScore(topic string) (*pubsub.TopicScoreSnapshot, error) {
	if m.pubsub == nil {
		return nil, nil
	}

	topicOpts, err := m.pubsub.Topic(topic)
	if err != nil {
		return nil, err
	}

	return topicOpts.Score(), nil
}

func (m *PubSubManager) SetTopicScore(topic string, params *pubsub.TopicScoreParams) error {
	if m.pubsub == nil {
		return nil
	}

	return m.pubsub.SetTopicScore(topic, params)
}

type Subscription struct {
	topic    string
	sub      *pubsub.Subscription
	handler  MessageHandler
	messages chan *Message
	ctx      context.Context
	cancel   context.CancelFunc
}

func (s *Subscription) readLoop() {
	s.ctx, s.cancel = context.WithCancel(context.Background())

	for {
		msg, err := s.sub.Next(s.ctx)
		if err != nil {
			if s.ctx.Err() != nil {
				return
			}
			continue
		}

		message := &Message{
			ID:          msg.ID,
			Data:        msg.Data,
			From:        msg.From,
			Seqno:       msg.Seqno,
			Topic:       s.topic,
			Signature:   msg.Signature,
			Key:         msg.Key,
			ReceivedAt:  msg.ReceivedAt,
		}

		if s.handler != nil {
			if err := s.handler(s.ctx, message); err != nil {
				continue
			}
		}

		select {
		case s.messages <- message:
		default:
		}
	}
}

func (s *Subscription) Messages() <-chan *Message {
	return s.messages
}

func (s *Subscription) Cancel() error {
	if s.cancel != nil {
		s.cancel()
	}

	if s.sub != nil {
		return s.sub.Cancel()
	}

	return nil
}

func (s *Subscription) Topic() string {
	return s.topic
}

type Message struct {
	ID         string
	Data       []byte
	From       []byte
	Seqno      []byte
	Topic      string
	Signature  []byte
	Key        []byte
	ReceivedAt interface{}
}
