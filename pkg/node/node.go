package node

import (
	"context"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	ws "github.com/libp2p/go-libp2p/p2p/transport/websocket"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	dhtopts "github.com/libp2p/go-libp2p-kad-dht/opts"
	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/your-org/p2p-network/pkg/dht"
	"github.com/your-org/p2p-network/pkg/pubsub"
	"github.com/your-org/p2p-network/pkg/protocol"
	"github.com/your-org/p2p-network/pkg/utils"
)

type Node struct {
	host   host.Host
	dht    *dht.DHTManager
	pubsub *pubsub.PubSubManager
	proto  *protocol.Handler

	ctx    context.Context
	cancel context.CancelFunc
	cfg    *Config
	logger *utils.Logger
}

func NewNode(cfg *Config, opts ...Option) (*Node, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	n := &Node{
		cfg: cfg,
	}

	for _, opt := range opts {
		opt(n)
	}

	if n.logger == nil {
		logger, err := utils.NewLogger("node", utils.LogLevelInfo)
		if err != nil {
			return nil, fmt.Errorf("create logger: %w", err)
		}
		n.logger = logger
	}

	host, err := n.createHost()
	if err != nil {
		return nil, fmt.Errorf("create host: %w", err)
	}
	n.host = host

	dhtMgr, err := n.createDHT()
	if err != nil {
		host.Close()
		return nil, fmt.Errorf("create DHT: %w", err)
	}
	n.dht = dhtMgr

	pubSubMgr, err := n.createPubSub()
	if err != nil {
		host.Close()
		return nil, fmt.Errorf("create PubSub: %w", err)
	}
	n.pubsub = pubSubMgr

	n.proto = protocol.NewHandler(n)

	n.ctx, n.cancel = context.WithCancel(context.Background())

	return n, nil
}

func (n *Node) createHost() (host.Host, error) {
	var opts []libp2p.Option

	opts = append(opts, libp2p.ListenAddrStrings(
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", n.cfg.ListenPort),
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%s/ws", n.cfg.ListenPort),
	))

	if n.cfg.EnableRelay {
		opts = append(opts, libp2p.EnableRelay())
	}

	if n.cfg.DisableMDNS {
		opts = append(opts, libp2p.DisableDiscovery())
	}

	opts = append(opts, libp2p.Transport(tcp.NewTCPTransport))
	opts = append(opts, libp2p.Transport(ws.New))

	host, err := libp2p.New(opts...)
	if err != nil {
		return nil, err
	}

	return host, nil
}

func (n *Node) createDHT() (*dht.DHTManager, error) {
	var opts []dhtopts.Option
	opts = append(opts, dhtopts.Client(true))

	dhtClient, err := dht.New(n.ctx, n.host, opts...)
	if err != nil {
		return nil, err
	}

	return dht.NewManager(dhtClient), nil
}

func (n *Node) createPubSub() (*pubsub.PubSubManager, error) {
	ps, err := pubsub.NewPubSub(n.ctx, n.host,
		pubsub.WithMessageSigning(true),
		pubsub.WithStrictSignatureVerification(true),
	)
	if err != nil {
		return nil, err
	}

	return pubsub.NewManager(ps), nil
}

func (n *Node) Start(ctx context.Context) error {
	n.ctx, n.cancel = context.WithCancel(ctx)

	if err := n.dht.Bootstrap(n.ctx); err != nil {
		return fmt.Errorf("bootstrap DHT: %w", err)
	}

	for _, addr := range n.cfg.BootstrapPeers {
		pi, err := peer.AddrInfoFromString(addr)
		if err != nil {
			n.logger.Warn("Invalid bootstrap peer", "addr", addr, "error", err)
			continue
		}

		if err := n.host.Connect(n.ctx, *pi); err != nil {
			n.logger.Warn("Failed to connect to bootstrap peer", "peer", pi.ID, "error", err)
			continue
		}
		n.logger.Info("Connected to bootstrap peer", "peer", pi.ID)
	}

	n.host.SetStreamHandler(protocol.ProtocolID, n.proto.HandleStream)

	n.logger.Info("Node started", "peerID", n.ID(), "addrs", n.Addrs())
	return nil
}

func (n *Node) Stop(ctx context.Context) error {
	n.cancel()

	if n.host != nil {
		n.host.Close()
	}

	n.logger.Info("Node stopped")
	return nil
}

func (n *Node) ID() string {
	return n.host.ID().String()
}

func (n *Node) Addrs() []string {
	addrs := make([]string, 0)
	for _, addr := range n.host.Addrs() {
		addrs = append(addrs, addr.String())
	}
	return addrs
}

func (n *Node) Host() host.Host {
	return n.host
}

func (n *Node) DHT() *dht.DHTManager {
	return n.dht
}

func (n *Node) PubSub() *pubsub.PubSubManager {
	return n.pubsub
}

func (n *Node) Context() context.Context {
	return n.ctx
}

func (n *Node) Connect(ctx context.Context, pi peer.AddrInfo) error {
	return n.host.Connect(ctx, pi)
}

func (n *Node) Disconnect(ctx context.Context, p peer.ID) error {
	return n.host.ClosePeer(p)
}

func (n *Node) OpenStream(ctx context.Context, p peer.ID) (network.Stream, error) {
	return n.host.NewStream(ctx, p, protocol.ProtocolID)
}
