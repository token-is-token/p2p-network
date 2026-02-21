package dht

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/peer"
)

type RoutingManager struct {
	dht *dht.IpfsDHT
}

func NewRoutingManager(dhtClient *dht.IpfsDHT) *RoutingManager {
	return &RoutingManager{
		dht: dhtClient,
	}
}

func (r *RoutingManager) RefreshRoutingTable(ctx context.Context) error {
	if r.dht == nil {
		return nil
	}

	return r.dht.RefreshRoutingTable()
}

func (r *RoutingManager) GetPeerInfos() []PeerInfo {
	if r.dht == nil || r.dht.RoutingTable() == nil {
		return nil
	}

	peers := r.dht.RoutingTable().GetPeers()
	infos := make([]PeerInfo, 0, len(peers))

	for _, p := range peers {
		infos = append(infos, PeerInfo{
			ID:       p.String(),
			LastSeen: time.Now(),
		})
	}

	return infos
}

type PeerInfo struct {
	ID       string
	LastSeen time.Time
}

func (r *RoutingManager) FindPeer(ctx context.Context, id peer.ID) (peer.AddrInfo, error) {
	if r.dht == nil {
		return peer.AddrInfo{}, nil
	}
	return r.dht.FindPeer(ctx, id)
}

func (r *RoutingManager) FindProviders(ctx context.Context, key string, limit int) ([]peer.AddrInfo, error) {
	if r.dht == nil {
		return nil, nil
	}

	if limit <= 0 {
		limit = 20
	}

	ch, err := r.dht.FindProvidersAsync(ctx, []byte(key), limit)
	if err != nil {
		return nil, err
	}

	var results []peer.AddrInfo
	for info := range ch {
		results = append(results, info)
		if len(results) >= limit {
			break
		}
	}

	return results, nil
}

func (r *RoutingManager) GetClosestPeers(ctx context.Context, key string) ([]peer.ID, error) {
	if r.dht == nil {
		return nil, nil
	}

	return r.dht.GetClosestPeers(ctx, []byte(key))
}

func (r *RoutingManager) IsOnline() bool {
	return r.dht != nil
}
