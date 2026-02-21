package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeCreation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Skip("Requires libp2p setup - run manually")

	cfg := DefaultConfig()
	cfg.ListenPort = "0"

	node, err := NewNode(cfg)
	require.NoError(t, err)
	require.NotNil(t, node)

	err = node.Start(ctx)
	require.NoError(t, err)

	assert.NotEmpty(t, node.ID())
	assert.NotEmpty(t, node.Addrs())

	err = node.Stop(ctx)
	require.NoError(t, err)
}

func TestNodeConnection(t *testing.T) {
	t.Skip("Requires libp2p setup - run manually")
}

func TestPubSubPublishSubscribe(t *testing.T) {
	t.Skip("Requires libp2p setup - run manually")
}

func TestDHTFindPeer(t *testing.T) {
	t.Skip("Requires libp2p setup - run manually")
}

func TestDiscovery(t *testing.T) {
	t.Skip("Requires libp2p setup - run manually")
}

func TestRelayConnection(t *testing.T) {
	t.Skip("Requires libp2p setup - run manually")
}

func TestProtocolCommunication(t *testing.T) {
	t.Skip("Requires libp2p setup - run manually")
}

func TestNodeGracefulShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	cfg := DefaultConfig()
	cfg.ListenPort = "0"

	node, err := NewNode(cfg)
	require.NoError(t, err)

	err = node.Start(ctx)
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	cancel()

	err = node.Stop(ctx)
	require.NoError(t, err)
}
