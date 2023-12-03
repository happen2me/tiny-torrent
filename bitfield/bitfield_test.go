package bitfield

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasPiece(t *testing.T) {
	// Create a test bitfield
	bitfield := Bitfield{0b10101010, 0b01010101}

	// Test cases
	testCases := []struct {
		index         int
		expectedValue bool
	}{
		{0, true},   // First bit is set
		{1, false},  // Second bit is not set
		{7, false},  // Eighth bit of the first byte is not set
		{8, false},  // First bit of the second byte is not set
		{9, true},   // Second bit of the second byte is set
		{15, true},  // Eighth bit of the second byte is set
		{16, false}, // Index out of range
	}

	// Run the test cases
	for _, tc := range testCases {
		actualValue := bitfield.HasPiece(tc.index)
		assert.Equal(t, tc.expectedValue, actualValue)
	}
}
