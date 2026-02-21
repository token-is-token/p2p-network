package discovery

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

type BootstrapDiscovery struct {
	peers   []peer.AddrInfo
	refresh time.Duration
}

func NewBootstrapDiscovery(peers []peer.AddrInfo) *BootstrapDiscovery {
	return &BootstrapDiscovery{
		peers:   peers,
		refresh: 30 * time.Minute,
	}
}

func (b *BootstrapDiscovery) GetPeers() []peer.AddrInfo {
	return b.peers
}

func (b *BootstrapDiscovery) AddPeer(addr string) error {
	maddr, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		return err
	}

	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return err
	}

	b.peers = append(b.peers, *info)
	return nil
}

func (b *BootstrapDiscovery) RemovePeer(peerID peer.ID) {
	newPeers := make([]peer.AddrInfo, 0)
	for _, p := range b.peers {
		if p.ID != peerID {
			newPeers = append(newPeers, p)
		}
	}
	b.peers = newPeers
}

func (b *BootstrapDiscovery) Refresh() error {
	return nil
}

func (b *BootstrapDiscovery) Discover(ctx context.Context) <-chan peer.AddrInfo {
	ch := make(chan peer.AddrInfo, len(b.peers))

	go func() {
		defer close(ch)
		for _, p := range b.peers {
			select {
			case <-ctx.Done():
				return
			case ch <- p:
			}
		}
	}()

	return ch
}

func ParseBootstrapAddr(addr string) (peer.AddrInfo, error) {
	maddr, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		return peer.AddrInfo{}, err
	}

	return peer.AddrInfoFromP2pAddr(maddr)
}
