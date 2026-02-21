package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.NotNil(t, cfg)
	assert.Equal(t, "0", cfg.ListenPort)
	assert.True(t, cfg.EnableRelay)
	assert.False(t, cfg.DisableMDNS)
	assert.Equal(t, "llm-share", cfg.NetworkName)
}

func TestConfigOptions(t *testing.T) {
	cfg := DefaultConfig()

	assert.NotNil(t, cfg)

	cfg.ListenPort = "4001"
	cfg.EnableRelay = false

	assert.Equal(t, "4001", cfg.ListenPort)
	assert.False(t, cfg.EnableRelay)
}

func TestKadDHTConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.True(t, cfg.EnableDHT)
	assert.Equal(t, "client", cfg.Mode)
}

func TestPubSubConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.True(t, cfg.EnablePubSub)
	assert.True(t, cfg.PubSubSignMessages)
	assert.True(t, cfg.PubSubValidateMessages)
}

func TestDiscoveryConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.True(t, cfg.EnableMDNS)
	assert.Equal(t, "_llm-share._tcp", cfg.MDNSServiceName)
	assert.Equal(t, "llm-share-p2p", cfg.Rendezvous)
}
