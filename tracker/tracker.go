// Responsibilities:
// - Implement the communication with the tracker specified in the .torrent file.
// - Send periodic requests to the tracker to get a list of peers.
// - Handle tracker response to get peer IPs and ports.
package tracker

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jackpal/bencode-go"

	"github.com/happen2me/tiny-torrent/peer"
	"github.com/happen2me/tiny-torrent/torrentfile"
)

// buildTrackerURL builds a URL for the tracker request.
// The peerID and port are used to identify our client to the tracker.
func buildTrackerURL(t *torrentfile.TorrentFile, peerID [20]byte, port uint16) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}
	base.RawQuery = url.Values{
		"info_hash":  []string{string(t.InfoHash[:])}, // convert rune to string
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(port))}, // Convert integer to decimal string
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"left":       []string{strconv.Itoa(int(t.Info.Length))},
		"compact":    []string{"1"},
	}.Encode()
	return base.String(), nil
}

func generatePeerID() [20]byte {
	// generate a random peer ID with rand
	peerID := [20]byte{}
	rand.Read(peerID[:])
	peerID[0] = '-'
	peerID[1] = 'T'
	peerID[2] = 'T'
	peerID[3] = '0'
	peerID[4] = '0'
	return peerID
}

type bencodeTracker struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func parsePeers(peersBin []byte) ([]peer.Peer, error) {
	// 1. Check if peers is a list
	const peerSize = 6
	btPeers := []byte(peersBin)
	if len(btPeers)%peerSize != 0 {
		err := fmt.Errorf("received invalid peers string of length %d", len(btPeers))
		return []peer.Peer{}, err
	}
	// 2. Iterate through list
	var peers []peer.Peer
	for i := 0; i < len(btPeers); i += peerSize {
		// 3. For each peer, parse IP and port
		peer := peer.Peer{
			IP:   net.IP(btPeers[i : i+4]),
			Port: binary.BigEndian.Uint16(btPeers[i+4 : i+6]),
		}
		peers = append(peers, peer)
	}
	// 4. Return list of peers
	return peers, nil
}

func requestPeers(t *torrentfile.TorrentFile, peerID [20]byte, port uint16) ([]peer.Peer, error) {
	// 1. Build the tracker URL
	url, err := buildTrackerURL(t, peerID, port)
	if err != nil {
		return []peer.Peer{}, err
	}
	// 2. Send GET request to tracker URL
	c := &http.Client{Timeout: 15 * time.Second}
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 3. Read and parse the response
	tracker := bencodeTracker{}
	err = bencode.Unmarshal(resp.Body, &tracker)
	if err != nil {
		return nil, err
	}
	// 4. Return list of peers
	return parsePeers([]byte(tracker.Peers))
}
