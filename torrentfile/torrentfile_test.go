package torrentfile

import (
	"os"
	"reflect"
	"testing"
)

func TestOpen(t *testing.T) {
	// Create a temporary test file
	file, err := os.CreateTemp("", "test_torrentfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	// Write test data to the file
	_, err = file.WriteString("d8:announce41:http://bttracker.debian.org:6969/announce7:comment35:\"Debian CD from cdimage.debian.org\"10:created by13:mktorrent 1.113:creation datei1696680176e4:infod6:lengthi658505728e4:name31:debian-12.2.0-amd64-netinst.iso12:piece lengthi262144e6:pieces20:01234567890123456789ee")
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	// Call the Open function
	tf, err := Open(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Verify the result
	expected := TorrentFile{
		Announce: "http://bttracker.debian.org:6969/announce",
		InfoHash: [20]uint8{141, 18, 124, 159, 106, 156, 110, 17, 248, 89, 151, 127, 208, 50, 71, 16, 94, 4, 190, 173},
		Info: struct {
			Length      int64
			Name        string
			PieceLength int64
			Pieces      [][20]byte
		}{
			Length:      658505728,
			Name:        "debian-12.2.0-amd64-netinst.iso",
			PieceLength: 262144,
			Pieces:      [][20]byte{[20]uint8{48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57}},
		},
	}
	if !reflect.DeepEqual(tf, expected) {
		t.Errorf("Open() = %v, expected %v", tf, expected)
	}
}

func TestOpen_InvalidFile(t *testing.T) {
	// Call the Open function with an invalid file path
	_, err := Open("nonexistent_file")
	if err == nil {
		t.Error("Expected an error, but got nil")
	}
}
