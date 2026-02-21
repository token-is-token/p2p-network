package protocol

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMessageTypes(t *testing.T) {
	assert.Equal(t, MessageType(0), MsgTypeRequest)
	assert.Equal(t, MessageType(1), MsgTypeResponse)
	assert.Equal(t, MessageType(2), MsgTypeHeartbeat)
	assert.Equal(t, MessageType(3), MsgTypePing)
	assert.Equal(t, MessageType(4), MsgTypePong)
}

func TestNewRequest(t *testing.T) {
	msg := NewRequest("test-id", []byte("test payload"))

	assert.Equal(t, MsgTypeRequest, msg.Type)
	assert.Equal(t, "test-id", msg.RequestID)
	assert.Equal(t, []byte("test payload"), msg.Payload)
}

func TestNewResponse(t *testing.T) {
	msg := NewResponse("test-id", []byte("test payload"))

	assert.Equal(t, MsgTypeResponse, msg.Type)
	assert.Equal(t, "test-id", msg.RequestID)
	assert.Equal(t, []byte("test payload"), msg.Payload)
}

func TestNewHeartbeat(t *testing.T) {
	msg := NewHeartbeat("peer-id")

	assert.Equal(t, MsgTypeHeartbeat, msg.Type)
	assert.Equal(t, "peer-id", msg.RequestID)
	assert.Nil(t, msg.Payload)
}

func TestMessageEncodeDecode(t *testing.T) {
	original := &Message{
		Type:      MsgTypeRequest,
		RequestID: "test-id",
		Payload:   []byte("test payload"),
		Timestamp: time.Now().Unix(),
	}

	encoded, err := original.Encode()
	assert.NoError(t, err)
	assert.NotEmpty(t, encoded)

	decoded, err := DecodeMessage(encoded)
	assert.NoError(t, err)
	assert.Equal(t, original.Type, decoded.Type)
	assert.Equal(t, original.RequestID, decoded.RequestID)
	assert.Equal(t, original.Payload, decoded.Payload)
}

func TestMessageEncodeTooShort(t *testing.T) {
	_, err := DecodeMessage([]byte{})
	assert.Error(t, err)
}

func TestProtocolID(t *testing.T) {
	assert.Equal(t, ProtocolIDStr, ProtocolID)
}
