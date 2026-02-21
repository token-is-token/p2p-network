package node

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/your-org/p2p-network/pkg/utils"
)

type Option func(*Node)

func WithLogger(logger *utils.Logger) Option {
	return func(n *Node) {
		n.logger = logger
	}
}

func WithConfig(cfg *Config) Option {
	return func(n *Node) {
		n.cfg = cfg
	}
}

func WithBootstrapPeers(peers []peer.AddrInfo) Option {
	return func(n *Node) {
		n.cfg.BootstrapPeers = nil
		for _, p := range peers {
			n.cfg.BootstrapPeers = append(n.cfg.BootstrapPeers, p.String())
		}
	}
}

func WithListenAddr(addr string) Option {
	return func(n *Node) {
		n.cfg.ListenPort = addr
	}
}

func WithNetworkName(name string) Option {
	return func(n *Node) {
		n.cfg.NetworkName = name
	}
}

func WithDataDir(dir string) Option {
	return func(n *Node) {
		n.cfg.DataDir = dir
	}
}

func EnableRelay() Option {
	return func(n *Node) {
		n.cfg.EnableRelay = true
	}
}

func DisableRelay() Option {
	return func(n *Node) {
		n.cfg.EnableRelay = false
	}
}

func EnableMDNS() Option {
	return func(n *Node) {
		n.cfg.DisableMDNS = false
	}
}

func DisableMDNS() Option {
	return func(n *Node) {
		n.cfg.DisableMDNS = true
	}
}
