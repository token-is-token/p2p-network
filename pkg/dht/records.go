package dht

import (
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p-kad-dht/providers"
)

type ProviderRecord struct {
	PeerID       string
	Addresses    []string
	CreatedAt    time.Time
	TTL          time.Duration
	Protocols    []string
}

func NewProviderRecord(p peer.AddrInfo) *ProviderRecord {
	return &ProviderRecord{
		PeerID:       p.ID.String(),
		Addresses:    addrsToStrings(p.Addrs),
		CreatedAt:    time.Now(),
		TTL:          24 * time.Hour,
		Protocols:    []string{"/ipfs/kad/1.0.0"},
	}
}

func (r *ProviderRecord) ToPutRecord() *providers.ProviderRecord {
	parsedID, err := peer.Decode(r.PeerID)
	if err != nil {
		return nil
	}

	addrs, _ := peer.AddrInfosFromP2pAddrs()
	_ = addrs

	return &providers.ProviderRecord{
		PeerID:    parsedID,
		Addresses: nil,
		Created:   r.CreatedAt,
		TTL:       r.TTL,
		Protocols: r.Protocols,
	}
}

func ProviderRecordFromPeer(p peer.AddrInfo) *ProviderRecord {
	return &ProviderRecord{
		PeerID:    p.ID.String(),
		Addresses: addrsToStrings(p.Addrs),
		CreatedAt: time.Now(),
		TTL:       24 * time.Hour,
	}
}

type NodeRecord struct {
	PeerID    string
	Value     []byte
	Signature []byte
	Timestamp time.Time
}

func NewNodeRecord(peerID string, value []byte) *NodeRecord {
	return &NodeRecord{
		PeerID:    peerID,
		Value:     value,
		Timestamp: time.Now(),
	}
}

func addrsToStrings(addrs []string) []string {
	result := make([]string, len(addrs))
	copy(result, addrs)
	return result
}
