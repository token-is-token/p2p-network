package discovery

import (
	"context"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p-discovery"
)

type DiscoveryManager struct {
	host       host.Host
	mdns       *MDNSDiscovery
	bootstrap  *BootstrapDiscovery
	routing    discovery.RoutingDiscovery
	advertise  *discovery.Advertiser
	peers      chan peer.AddrInfo
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewDiscoveryManager(h host.Host) *DiscoveryManager {
	ctx, cancel := context.WithCancel(context.Background())

	return &DiscoveryManager{
		host:   h,
		peers:  make(chan peer.AddrInfo, 100),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (m *DiscoveryManager) Start(ctx context.Context, rendezvous string) error {
	m.routing = discovery.NewRoutingDiscovery(m.host)

	m.advertise = m.routing.Advertise(ctx, rendezvous)

	return nil
}

func (m *DiscoveryManager) Stop() error {
	m.cancel()
	if m.mdns != nil {
		m.mdns.Stop()
	}
	return nil
}

func (m *DiscoveryManager) FindPeers(ctx context.Context, rendezvous string) (<-chan peer.AddrInfo, error) {
	return m.routing.FindPeers(ctx, rendezvous)
}

func (m *DiscoveryManager) Discover(ctx context.Context, rendezvous string, limit int) ([]peer.AddrInfo, error) {
	peerChan, err := m.FindPeers(ctx, rendezvous)
	if err != nil {
		return nil, err
	}

	var peers []peer.AddrInfo
	for p := range peerChan {
		if p.ID == m.host.ID() {
			continue
		}
		peers = append(peers, p)
		if limit > 0 && len(peers) >= limit {
			break
		}
	}

	return peers, nil
}

func (m *DiscoveryManager) AddMDNS(serviceName string) error {
	mdns, err := NewMDNSDiscovery(m.host, serviceName)
	if err != nil {
		return err
	}

	m.mdns = mdns
	return mdns.Start()
}

func (m *DiscoveryManager) SetBootstrapPeers(peers []peer.AddrInfo) {
	m.bootstrap = NewBootstrapDiscovery(peers)
}

func (m *DiscoveryManager) GetBootstrapPeers() []peer.AddrInfo {
	if m.bootstrap == nil {
		return nil
	}
	return m.bootstrap.GetPeers()
}
