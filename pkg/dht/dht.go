package dht

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-kad-dht/providers"
	kb "github.com/libp2p/go-libp2p-kad-dht/routing"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/routing"
)

type DHTManager struct {
	dht           *dht.IpfsDHT
	providerStore *providers.ProviderStore
}

func NewManager(dhtClient *dht.IpfsDHT) *DHTManager {
	return &DHTManager{
		dht:           dhtClient,
		providerStore: providers.NewProviderStore(),
	}
}

func (m *DHTManager) Bootstrap(ctx context.Context) error {
	if m.dht == nil {
		return nil
	}

	bootstrapPeers := dht.DefaultBootstrapPeers()
	if len(bootstrapPeers) == 0 {
		return nil
	}

	for _, pi := range bootstrapPeers {
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		err := m.dht.BootstrapWithPeer(ctx, pi)
		cancel()
		if err != nil {
			continue
		}
	}

	return nil
}

func (m *DHTManager) PutProviderRecord(ctx context.Context, key string, record *ProviderRecord) error {
	if m.dht == nil {
		return nil
	}

	rec := record.ToPutRecord()
	return m.dht.ProviderStore().AddProvider(ctx, []byte(key), rec)
}

func (m *DHTManager) GetProviderRecords(ctx context.Context, key string) ([]*ProviderRecord, error) {
	if m.dht == nil {
		return nil, nil
	}

	providers, err := m.dht.ProviderStore().GetProviders(ctx, []byte(key))
	if err != nil {
		return nil, err
	}

	records := make([]*ProviderRecord, 0, len(providers))
	for _, p := range providers {
		records = append(records, ProviderRecordFromPeer(p))
	}

	return records, nil
}

func (m *DHTManager) PutNodeRecord(ctx context.Context, record *NodeRecord) error {
	if m.dht == nil {
		return nil
	}

	key := record.PeerID
	return m.dht.PutValue(ctx, []byte(key), record.Value)
}

func (m *DHTManager) GetNodeRecord(ctx context.Context, peerID string) (*NodeRecord, error) {
	if m.dht == nil {
		return nil, nil
	}

	value, err := m.dht.GetValue(ctx, []byte(peerID))
	if err != nil {
		if err == routing.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}

	record := &NodeRecord{
		PeerID: peerID,
		Value:  value,
	}

	return record, nil
}

func (m *DHTManager) FindPeer(ctx context.Context, peerID peer.ID) (peer.AddrInfo, error) {
	if m.dht == nil {
		return peer.AddrInfo{}, nil
	}

	return m.dht.FindPeer(ctx, peerID)
}

func (m *DHTManager) FindProviders(ctx context.Context, key string) ([]peer.AddrInfo, error) {
	if m.dht == nil {
		return nil, nil
	}

	ch, err := m.dht.FindProvidersAsync(ctx, []byte(key), 10)
	if err != nil {
		return nil, err
	}

	var providers []peer.AddrInfo
	for p := range ch {
		providers = append(providers, p)
	}

	return providers, nil
}

func (m *DHTManager) GetClosestPeers(ctx context.Context, key string) ([]peer.ID, error) {
	if m.dht == nil {
		return nil, nil
	}

	return m.dht.GetClosestPeers(ctx, []byte(key))
}

func (m *DHTManager) Provide(ctx context.Context, key string) error {
	if m.dht == nil {
		return nil
	}

	return m.dht.Provide(ctx, []byte(key), true)
}

func (m *DHTManager) RoutingTable() *kb.RoutingTable {
	if m.dht == nil {
		return nil
	}
	return m.dht.RoutingTable()
}

func (m *DHTManager) PeerCount() int {
	rt := m.RoutingTable()
	if rt == nil {
		return 0
	}
	return rt.Size()
}
