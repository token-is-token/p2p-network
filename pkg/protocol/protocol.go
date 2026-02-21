package protocol

import (
	"encoding/binary"
	"fmt"
)

const (
	ProtocolIDStr = "/llm-share/1.0.0"
)

var (
	ProtocolID = ProtocolIDStr
)

type MessageType uint8

const (
	MsgTypeRequest MessageType = iota
	MsgTypeResponse
	MsgTypeHeartbeat
	MsgTypePing
	MsgTypePong
)

type Message struct {
	Type      MessageType
	RequestID string
	Payload   []byte
	Signature []byte
	Timestamp int64
}

func (m *Message) Encode() ([]byte, error) {
	data := make([]byte, 0, 4+len(m.RequestID)+len(m.Payload)+len(m.Signature))

	data = append(data, byte(m.Type))

	reqIDLen := uint16(len(m.RequestID))
	data = binary.LittleEndian.AppendUint16(data, reqIDLen)
	data = append(data, m.RequestID...)

	payloadLen := uint32(len(m.Payload))
	data = binary.LittleEndian.AppendUint32(data, payloadLen)
	data = append(data, m.Payload...)

	sigLen := uint16(len(m.Signature))
	data = binary.LittleEndian.AppendUint16(data, sigLen)
	data = append(data, m.Signature...)

	data = binary.LittleEndian.AppendUint64(data, uint64(m.Timestamp))

	return data, nil
}

func DecodeMessage(data []byte) (*Message, error) {
	if len(data) < 1 {
		return nil, fmt.Errorf("data too short")
	}

	msg := &Message{}
	offset := 0

	msg.Type = MessageType(data[offset])
	offset++

	if len(data) < offset+2 {
		return nil, fmt.Errorf("data too short for request ID length")
	}
	reqIDLen := binary.LittleEndian.Uint16(data[offset : offset+2])
	offset += 2

	if len(data) < offset+int(reqIDLen) {
		return nil, fmt.Errorf("data too short for request ID")
	}
	msg.RequestID = string(data[offset : offset+int(reqIDLen)])
	offset += int(reqIDLen)

	if len(data) < offset+4 {
		return nil, fmt.Errorf("data too short for payload length")
	}
	payloadLen := binary.LittleEndian.Uint32(data[offset : offset+4])
	offset += 4

	if len(data) < offset+int(payloadLen) {
		return nil, fmt.Errorf("data too short for payload")
	}
	msg.Payload = data[offset : offset+int(payloadLen)]
	offset += int(payloadLen)

	if len(data) < offset+2 {
		return nil, fmt.Errorf("data too short for signature length")
	}
	sigLen := binary.LittleEndian.Uint16(data[offset : offset+2])
	offset += 2

	if len(data) < offset+int(sigLen) {
		return nil, fmt.Errorf("data too short for signature")
	}
	msg.Signature = data[offset : offset+int(sigLen)]
	offset += int(sigLen)

	if len(data) < offset+8 {
		return nil, fmt.Errorf("data too short for timestamp")
	}
	msg.Timestamp = int64(binary.LittleEndian.Uint64(data[offset : offset+8]))

	return msg, nil
}

func NewRequest(requestID string, payload []byte) *Message {
	return &Message{
		Type:      MsgTypeRequest,
		RequestID: requestID,
		Payload:   payload,
		Timestamp: 0,
	}
}

func NewResponse(requestID string, payload []byte) *Message {
	return &Message{
		Type:      MsgTypeResponse,
		RequestID: requestID,
		Payload:   payload,
		Timestamp: 0,
	}
}

func NewHeartbeat(peerID string) *Message {
	return &Message{
		Type:      MsgTypeHeartbeat,
		RequestID: peerID,
		Payload:   nil,
		Timestamp: 0,
	}
}
