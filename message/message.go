// A message is composed of a 4-byte length prefix, a 1-byte ID, and a payload. The length prefix is the length of the message ID and payload.
package message

import (
	"encoding/binary"
	"io"
)

type messageID uint8

const (
	MsgChoke         messageID = 0
	MsgUnchoke       messageID = 1
	MsgInterested    messageID = 2
	MsgNotInterested messageID = 3
	MsgHave          messageID = 4
	MsgBitfield      messageID = 5
	MsgRequest       messageID = 6
	MsgPiece         messageID = 7
	MsgCancel        messageID = 8
)

type Message struct {
	ID      messageID
	Payload []byte
}

// Serialize serializes the message to a byte slice.
func (m *Message) Serialize() []byte {
	length := uint32(len(m.Payload) + 1) // 1 for message ID
	lengthBufLen := uint32(4)
	buf := make([]byte, lengthBufLen+length)
	// Set the length prefix.
	binary.BigEndian.PutUint32(buf[0:lengthBufLen], length)
	// Set the message ID.
	buf[lengthBufLen] = byte(m.ID)
	// Copy the payload.
	copy(buf[lengthBufLen+1:], m.Payload)
	return buf
}

// Read reads and parses a message from a reader stream.
// It returns nil if the keep-alive message is read.
func Read(r io.Reader) (*Message, error) {
	// Read the length prefix.
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lengthBuf)

	// Keep-alive message.
	if length == 0 {
		return nil, nil
	}

	// Read the message ID.
	idBuf := make([]byte, 1)
	_, err = io.ReadFull(r, idBuf)
	if err != nil {
		return nil, err
	}
	id := messageID(idBuf[0])
	// Read the payload.
	payload := make([]byte, length-1)
	_, err = io.ReadFull(r, payload)
	if err != nil {
		return nil, err
	}
	m := &Message{
		ID:      id,
		Payload: payload,
	}
	return m, nil
}
