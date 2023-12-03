package bitfield

type Bitfield []byte

// HasPiece returns true if the bitfield indicates that the peer has the piece at the given index.
func (bf Bitfield) HasPiece(index int) bool {
	byteIndex := index / 8
	offset := index % 8
	if byteIndex >= len(bf) {
		return false
	}
	return bf[byteIndex]>>(7-offset)&1 != 0
}

// SetPiece sets the bitfield to indicate that the peer has the piece at the given index.
func (bf Bitfield) SetPiece(index int) {
	byteIndex := index / 8
	offset := index % 8
	if byteIndex >= len(bf) {
		return
	}
	bf[byteIndex] |= 1 << (7 - offset)
}
