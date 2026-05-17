package torrent

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io"
	"os"

	"github.com/devops-in-the-east/bittorrent-go/bencode"
)

// TorrentFile represents the contents of a .torrent file
type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

// Open reads a .torrent file and parses it
func Open(path string) (*TorrentFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return Parse(file)
}

// Parse extract the necessary metadata from the torrent file.
func Parse(r io.Reader) (*TorrentFile, error) {
	reader := bufio.NewReader(r)

	// Decode the bencoded data
	decoded, err := bencode.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode torrent file: %w", err)
	}

	// The root of a .torrent file must be a dictionay
	dict, ok := decoded.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Expected Root Element to be a Dictionary")

	}

	// 1. Get Announce URL
	announceStr := ""
	if announce, ok := dict["announce"].(string); ok {
		announceStr = announce

	}

	// 2. Get Info Dictionary
	infoRaw, ok := dict["info"]
	if !ok {
		return nil, fmt.Errorf("missing 'info' dictionary")
	}

	infoDict, ok := infoRaw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("'info' is not a dictionary")
	}

	// 3. Extract basic info fields
	name, _ := infoDict["name"].(string)
	length, _ := infoDict["length"].(int)
	pieceLength, _ := infoDict["piece length"].(int)
	pieces, _ := infoDict["pieces"].(string)

	// 4. Split the piece string into a slice of 20-byte hashes
	if len(pieces)%20 != 0 {
		return nil, fmt.Errorf("Malformed pieces string")
	}
	numHashes := len(pieces) / 20
	pieceHashes := make([][20]byte, numHashes)
	for i := 0; i < numHashes; i++ {
		copy(pieceHashes[i][:], pieces[i*20:(i+1)*20])
	}

	// 5. Calculate InfoHash (TODO: We need a bencode encoder for this)

	var infoHash [20]byte
	_ = sha1.Sum([]byte("dummy hash until we build the encoder"))

	return &TorrentFile{
		Announce:    announceStr,
		InfoHash:    infoHash,
		PieceHashes: pieceHashes,
		PieceLength: pieceLength,
		Length:      length,
		Name:        name,
	}, nil

}

// Qs answered by the agent

// Qx not answered by the agent

// Lines 39, 45, 52 consists of multiple new lines with comments under the same function.
// I would like to know what is the use of having sperate new lines wihtin a fuvntion.

// On a separate note;

// dict, ok := decoded.(map[string]interface{})dict, ok := decoded.(map[string]interface{})
// Let me try explain the the logivcal construct within this lin.
// we are delcaring a object that is dic called ok that the type of the dic is a map that constists of string that can be a bunch of key value pairs.
// the interface is just a fancy way of saying that the object can be callled anywhere within the entire file.
// now that obejct is again assinged to another object that is also acts like its twin.
// this line can be written as y = x1 + x2.
