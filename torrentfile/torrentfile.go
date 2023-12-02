// Description: This package provides a struct that represents a torrent file.
package torrentfile

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/jackpal/bencode-go"
)

// bencodeTorrentFile is a struct that represents a torrent file.
type bencodeTorrentFile struct {
	Announce string             `bencode:"announce"`
	Info     bencodeTorrentInfo `bencode:"info"`
}

// bencodeTorrentInfo is a struct that represents the info section of a torrent file.
type bencodeTorrentInfo struct {
	Length      int64  `bencode:"length"`
	Name        string `bencode:"name"`
	PieceLength int64  `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
}

// TorrentFile is a struct that represents a torrent file, with pieces converted to byte and InfoHash.
type TorrentFile struct {
	Announce string
	InfoHash [20]byte
	Info     struct {
		Length      int64
		Name        string
		PieceLength int64
		Pieces      [][20]byte // List of hashes of each piece
	}
}

// splitPieceHashes splits the pieces string into a slice of 20-byte hashes.
func (bi *bencodeTorrentInfo) splitPieceHashes() ([][20]byte, error) {
	hashLen := 20 // SHA1 hashes are 20 bytes long
	piecesInBytes := []byte(bi.Pieces)
	if len(piecesInBytes)%hashLen != 0 {
		err := fmt.Errorf("received invalid pieces string of length %d", len(bi.Pieces))
		return [][20]byte{}, err
	}
	numHashes := len(piecesInBytes) / hashLen
	hashes := make([][20]byte, numHashes)
	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], piecesInBytes[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}

// hash returns the SHA1 hash of the bencoded Info struct.
func (bi *bencodeTorrentInfo) hash() ([20]byte, error) {
	buf := new(bytes.Buffer)
	err := bencode.Marshal(buf, *bi)
	if err != nil {
		return [20]byte{}, err
	}
	return sha1.Sum(buf.Bytes()), nil
}

// toTorrentFile converts a bencodeTorrentFile to a TorrentFile.
func (btf *bencodeTorrentFile) toTorrentFile() (TorrentFile, error) {
	tf := TorrentFile{}
	tf.Announce = btf.Announce
	tf.Info.Length = btf.Info.Length
	tf.Info.Name = btf.Info.Name
	tf.Info.PieceLength = btf.Info.PieceLength
	pieces, err := btf.Info.splitPieceHashes()
	if err != nil {
		return TorrentFile{}, err
	}
	tf.Info.Pieces = pieces
	infoHash, error := btf.Info.hash()
	if error != nil {
		return TorrentFile{}, error
	}
	tf.InfoHash = infoHash
	return tf, nil
}

// Open opens a torrent file and returns a TorrentFile.
func Open(path string) (TorrentFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return TorrentFile{}, err
	}
	defer file.Close()

	btf := bencodeTorrentFile{}
	err = bencode.Unmarshal(file, &btf)
	if err != nil {
		return TorrentFile{}, err
	}
	tf, err := btf.toTorrentFile()
	if err != nil {
		return TorrentFile{}, err
	}
	return tf, err
}
