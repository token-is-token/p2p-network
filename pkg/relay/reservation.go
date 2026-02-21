package relay

import (
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
)

type Reservation struct {
	Peer      peer.ID
	ExpireAt  time.Time
	Slots     int
	Addr      string
}

type ReservationManager struct {
	mu           sync.RWMutex
	reservations map[peer.ID]*Reservation
	maxSlots     int
}

func NewReservationManager() *ReservationManager {
	return &ReservationManager{
		reservations: make(map[peer.ID]*Reservation),
		maxSlots:    16,
	}
}

func (m *ReservationManager) Add(p peer.ID, expireAt time.Time) *Reservation {
	m.mu.Lock()
	defer m.mu.Unlock()

	rsv := &Reservation{
		Peer:     p,
		ExpireAt: expireAt,
		Slots:    1,
	}

	m.reservations[p] = rsv
	return rsv
}

func (m *ReservationManager) Remove(p peer.ID) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.reservations, p)
}

func (m *ReservationManager) Get(p peer.ID) (*Reservation, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	rsv, ok := m.reservations[p]
	if !ok {
		return nil, false
	}

	if time.Now().After(rsv.ExpireAt) {
		return nil, false
	}

	return rsv, true
}

func (m *ReservationManager) List() []*Reservation {
	m.mu.RLock()
	defer m.mu.RUnlock()

	list := make([]*Reservation, 0, len(m.reservations))
	for _, rsv := range m.reservations {
		if time.Now().Before(rsv.ExpireAt) {
			list = append(list, rsv)
		}
	}

	return list
}

func (m *ReservationManager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	for _, rsv := range m.reservations {
		if time.Now().Before(rsv.ExpireAt) {
			count++
		}
	}

	return count
}

func (m *ReservationManager) IsFull() bool {
	return m.Count() >= m.maxSlots
}

func (m *ReservationManager) Cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for p, rsv := range m.reservations {
		if now.After(rsv.ExpireAt) {
			delete(m.reservations, p)
		}
	}
}
