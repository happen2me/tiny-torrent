package message

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerialize(t *testing.T) {
	// Create a test message
	message := &Message{
		ID:      1,
		Payload: []byte("test payload"),
	}

	// Serialize the message
	serialized := message.Serialize()

	// Assert the length prefix
	expectedLength := uint32(len(message.Payload) + 1)
	actualLength := binary.BigEndian.Uint32(serialized[0:4])
	assert.Equal(t, expectedLength, actualLength)

	// Assert the message ID
	expectedID := byte(message.ID)
	actualID := serialized[4]
	assert.Equal(t, expectedID, actualID)

	// Assert the payload
	expectedPayload := message.Payload
	actualPayload := serialized[5:]
	assert.Equal(t, expectedPayload, actualPayload)
}

func TestRead(t *testing.T) {
	// Create a test reader with a serialized message
	serialized := []byte{
		0x00, 0x00, 0x00, 0x0D, // Length prefix: 13
		0x01,                                                                   // Message ID: 1
		0x74, 0x65, 0x73, 0x74, 0x20, 0x70, 0x61, 0x79, 0x6C, 0x6F, 0x61, 0x64, // Payload: "test payload"
	}
	reader := bytes.NewReader(serialized)

	// Call the Read function
	message, err := Read(reader)
	assert.NoError(t, err)

	// Assert the message ID
	expectedID := byte(1)
	actualID := byte(message.ID)
	assert.Equal(t, expectedID, actualID)

	// Assert the payload
	expectedPayload := []byte("test payload")
	actualPayload := message.Payload
	assert.Equal(t, expectedPayload, actualPayload)
}
