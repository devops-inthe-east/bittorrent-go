package tracker

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/devops-in-the-east/bittorrent-go/bencode"
)

type TrackerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

type Peer struct {
	IP   net.IP
	Port uint16
}

// UnmarshalPeers chops the compact binary string into a list of Peers
func UnmarshalPeers(peersBin string) ([]Peer, error) {
	const peerSize = 6 // 4 bytes for ip and 2 bytes for Port

	if len(peersBin)%peerSize != 0 {
		return nil, fmt.Errorf("received malformed")

	}

	numPeers := len(peersBin) / peerSize
	peers := make([]Peer, numPeers)

	for i := 0; i < numPeers; i++ {
		offset := i * peerSize

		// Slice out the 4 bytes for the IP Address
		peers[i].IP = net.IP([]byte(peersBin[offset : offset+4]))

		// Slice out the 4 bytes for the IP Address
		peers[i].Port = binary.BigEndian.Uint16([]byte(peersBin[offset+4 : offset+6]))

	}

	return peers, nil
}

func buildTrackerURL(announce string, infoHash [20]byte, peerID [20]byte, port uint16, length int) (string, error) {
	base, err := url.Parse(announce)
	if err != nil {
		return "", err
	}

	params := url.Values{
		"info_hash":  []string{string(infoHash[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"left":       []string{strconv.Itoa(length)},
		"compact":    []string{"1"},
	}

	base.RawQuery = params.Encode()
	return base.String(), nil
}

func requestPeers(trackerURL string) (*TrackerResponse, error) {
	resp, err := http.Get(trackerURL)
	if err != nil {
		return nil, err
	}

	// To an untrained eye, how can I get the intuition that resp method can be extended to
	// accommodate other functions like Body.Close() // I am confused about the body.close() method. When is it used?
	// I am also clusless that are we extending a object to accomdate a method?
	// are we using a variable that extends a methods to an object? if so how?
	// Which is which resp and body and Close?
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tracker responeded with status: %d",
			resp.StatusCode)

	}

	// 1. Unpack the Bencoded payload from the Cargo Box (Body)
	reader := bufio.NewReader(resp.Body) // Wrap the raw body!
	decoded, err := bencode.Decode(reader)
	if err != nil {
		return nil, err
	}

	// 2. We knmow the root of the tracker response should be a dictionary
	dict, ok := decoded.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Expected tracker response to be a dictionary")
	}

	// 3. Extract Interval and Peers from the dict

	interval, ok := dict["interval"].(int)
	if !ok {
		return nil, fmt.Errorf("tracker response mising interval")
	}
	peers, ok := dict["peers"].(string)
	if !ok {
		return nil, fmt.Errorf("tracker response missing peers")
	}

	// Here why do we have a prefix '&' to the TrackerResponse.
	//

	return &TrackerResponse{
		Interval: interval,
		Peers:    peers,
	}, nil

}
