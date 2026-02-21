package node

import (
	"time"
)

type Config struct {
	ListenPort      string
	BootstrapPeers  []string
	EnableRelay     bool
	DisableMDNS     bool
	NetworkName     string
	DataDir         string

	 KadDHTConfig
	PubSubConfig
	DiscoveryConfig
}

type KadDHTConfig struct {
	EnableDHT        bool
	RoutingDBDir     string
	Mode             string
	BootstrapTimeout time.Duration
}

type PubSubConfig struct {
	EnablePubSub      bool
	PubSubSignMessages bool
	PubSubValidateMessages bool
}

type DiscoveryConfig struct {
	EnableMDNS       bool
	MDNSServiceName string
	Rendezvous      string
}

func DefaultConfig() *Config {
	return &Config{
		ListenPort:     "0",
		BootstrapPeers: DefaultBootstrapPeers(),
		EnableRelay:    true,
		DisableMDNS:   false,
		NetworkName:   "llm-share",
		DataDir:        ".p2p-data",

		KadDHTConfig: KadDHTConfig{
			EnableDHT:        true,
			Mode:             "client",
			BootstrapTimeout: 30 * time.Second,
		},

		PubSubConfig: PubSubConfig{
			EnablePubSub:          true,
			PubSubSignMessages:    true,
			PubSubValidateMessages: true,
		},

		DiscoveryConfig: DiscoveryConfig{
			EnableMDNS:       true,
			MDNSServiceName: "_llm-share._tcp",
			Rendezvous:      "llm-share-p2p",
		},
	}
}

func DefaultBootstrapPeers() []string {
	return []string{}
}

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
