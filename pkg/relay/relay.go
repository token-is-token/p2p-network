package relay

import (
	"context"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	relayv1 "github.com/libp2p/go-libp2p-relay"
	"github.com/libp2p/go-libp2p/p2p/host/relay"
)

type RelayManager struct {
	host      host.Host
	relay     *relayv1.Relay
	enabled   bool
	reservations *ReservationManager
}

func NewRelayManager(h host.Host) *RelayManager {
	return &RelayManager{
		host:  h,
		enabled: true,
		reservations: NewReservationManager(),
	}
}

func (m *RelayManager) Enable() {
	m.enabled = true
}

func (m *RelayManager) Disable() {
	m.enabled = false
}

func (m *RelayManager) IsEnabled() bool {
	return m.enabled
}

func (m *RelayManager) Connect(ctx context.Context, p peer.ID) (peer.AddrInfo, error) {
	if !m.enabled {
		return peer.AddrInfo{}, nil
	}

	var addrInfo peer.AddrInfo
	var err error

	m.host.Peerstore().UpdateAddrs(p, m.host.Peerstore().Addrs(p))

	conns := m.host.Network().ConnsToPeer(p)
	for _, conn := range conns {
		relayAddrs := conn.RemoteMultiaddrs()
		for _, addr := range relayAddrs {
			_, isRelayAddr := addr.ValueForProtocol(multiaddr.P_Relay)
			if isRelayAddr {
				relayPID, _ := peer.IDFromBytes([]byte("relay"))
				addrInfo.ID = relayPID
				addrInfo.Addrs = append(addrInfo.Addrs, addr)
			}
		}
	}

	return addrInfo, err
}

func (m *RelayManager) Reservation(ctx context.Context, p peer.ID) (*Reservation, error) {
	return m.reservations.Get(p)
}

func (m *RelayManager) ListReservations() []*Reservation {
	return m.reservations.List()
}

func (m *RelayManager) AcceptReservations() bool {
	return true
}

func (m *RelayManager) MaxReservations() int {
	return 16
}

func (m *RelayManager) MaxCircuitSlots() int {
	return 1
}
