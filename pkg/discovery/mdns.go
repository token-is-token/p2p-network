package discovery

import (
	"context"
	"sync"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	mdns "github.com/libp2p/go-mdns"
)

type MDNSDiscovery struct {
	host        host.Host
	serviceName string
	peerStore   map[peer.ID]*peer.AddrInfo
	mu          sync.RWMutex
	service     *mdns.Service
	notif       chan peer.AddrInfo
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewMDNSDiscovery(h host.Host, serviceName string) (*MDNSDiscovery, error) {
	ctx, cancel := context.WithCancel(context.Background())

	m := &MDNSDiscovery{
		host:        h,
		serviceName: serviceName,
		peerStore:   make(map[peer.ID]*peer.AddrInfo),
		notif:       make(chan peer.AddrInfo, 100),
		ctx:         ctx,
		cancel:      cancel,
	}

	return m, nil
}

func (m *MDNSDiscovery) Start() error {
	service, err := mdns.NewMdnsService(m.ctx, m.host, m.serviceName)
	if err != nil {
		return err
	}

	m.service = service

	go m.handleUpdates()

	return service.Start()
}

func (m *MDNSDiscovery) Stop() {
	m.cancel()
	if m.service != nil {
		m.service.Close()
	}
}

func (m *MDNSDiscovery) handleUpdates() {
	for {
		select {
		case <-m.ctx.Done():
			return
		case info := <-m.service.Channel():
			m.mu.Lock()
			m.peerStore[info.ID] = &info
			m.mu.Unlock()

			select {
			case m.notif <- info:
			default:
			}
		}
	}
}

func (m *MDNSDiscovery) GetPeers() []peer.AddrInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	peers := make([]peer.AddrInfo, 0, len(m.peerStore))
	for _, info := range m.peerStore {
		peers = append(peers, *info)
	}

	return peers
}

func (m *MDNSDiscovery) PeerChan() <-chan peer.AddrInfo {
	return m.notif
}

func (m *MDNSDiscovery) HasPeer(id peer.ID) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, ok := m.peerStore[id]
	return ok
}
