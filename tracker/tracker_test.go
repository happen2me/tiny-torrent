package tracker

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/happen2me/tiny-torrent/torrentfile"
	"github.com/stretchr/testify/assert"
)

func TestRequestPeers(t *testing.T) {
	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a successful response from the tracker
		response := []byte(
			"d" +
				"8:interval" + "i900e" +
				"5:peers" + "12:" +
				string([]byte{
					192, 0, 2, 123, 0x1A, 0xE1, // 0x1AE1 = 6881
					127, 0, 0, 1, 0x1A, 0xE9, // 0x1AE9 = 6889
				}) + "e")
		w.Write(response)
	}))
	defer mockServer.Close()

	// Create a mock torrent file
	mockTorrentFile := &torrentfile.TorrentFile{
		Announce: mockServer.URL,
		InfoHash: [20]byte{},
		Info:     torrentfile.TorrentInfo{},
	}

	// Generate a mock peer ID and port
	mockPeerID := generatePeerID()
	mockPort := uint16(12345)

	// Call the function under test
	peers, err := requestPeers(mockTorrentFile, mockPeerID, mockPort)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert the expected number of peers
	assert.Len(t, peers, 2)

	// Assert the expected peer IP and port values
	assert.Equal(t, "192.0.2.123", peers[0].IP.String())
	assert.Equal(t, uint16(6881), peers[0].Port)
	assert.Equal(t, "127.0.0.1", peers[1].IP.String())
	assert.Equal(t, uint16(6889), peers[1].Port)
}

func TestBuildTrackerURL(t *testing.T) {
	// Create a mock torrent file
	mockTorrentFile := &torrentfile.TorrentFile{
		Announce: "http://tracker.example.com/announce",
		InfoHash: [20]byte{},
		Info:     torrentfile.TorrentInfo{},
	}

	// Generate a mock peer ID and port
	mockPeerID := [20]byte{}
	mockPort := uint16(12345)

	// Call the function under test
	trackerURL, err := buildTrackerURL(mockTorrentFile, mockPeerID, mockPort)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert the expected tracker URL
	expectedURL := "http://tracker.example.com/announce?compact=1&downloaded=0&info_hash=%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00&left=0&peer_id=%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00%00&port=12345&uploaded=0"
	assert.Equal(t, expectedURL, trackerURL)
}
