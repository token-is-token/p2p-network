package pubsub

const (
	TopicProviders  = "llm-share.providers"
	TopicRequests   = "llm-share.requests"
	TopicResponses  = "llm-share.responses"
	TopicHeartbeat  = "llm-share.heartbeat"
	TopicDiscovery  = "llm-share.discovery"
	TopicSync       = "llm-share.sync"
	TopicBroadcast  = "llm-share.broadcast"
)

var Topics = []string{
	TopicProviders,
	TopicRequests,
	TopicResponses,
	TopicHeartbeat,
	TopicDiscovery,
	TopicSync,
	TopicBroadcast,
}

type TopicConfig struct {
	Name      string
	Score     *TopicScoreConfig
	Validator MessageValidator
}

type TopicScoreConfig struct {
	ByTopicScoreWeight          float64
	IPColocationFactorWeight   float64
	IPColocationFactorThreshold float64
	BehaviourPenaltyWeight     float64
	BehaviourPenaltyThreshold  float64
	TimeInMeshWeight           float64
	TimeInMeshQuantum          float64
	FirstMessageDeliveriesWeight float64
	MessageDeliveriesWeight    float64
	MeshMessageDeliveriesWeight float64
	InvalidMessageDeliveriesWeight float64
}

func DefaultTopicConfigs() map[string]*TopicConfig {
	return map[string]*TopicConfig{
		TopicProviders: {
			Name: TopicProviders,
			Score: &TopicScoreConfig{
				ByTopicScoreWeight:          0.5,
				TimeInMeshWeight:           0.5,
				FirstMessageDeliveriesWeight: 1.0,
				InvalidMessageDeliveriesWeight: -100.0,
			},
		},
		TopicRequests: {
			Name: TopicRequests,
			Score: &TopicScoreConfig{
				ByTopicScoreWeight:            0.3,
				TimeInMeshWeight:            0.3,
				MessageDeliveriesWeight:       1.0,
				InvalidMessageDeliveriesWeight: -100.0,
			},
		},
		TopicResponses: {
			Name: TopicResponses,
			Score: &TopicScoreConfig{
				ByTopicScoreWeight:            0.3,
				TimeInMeshWeight:            0.3,
				MessageDeliveriesWeight:       1.0,
				InvalidMessageDeliveriesWeight: -100.0,
			},
		},
		TopicHeartbeat: {
			Name: TopicHeartbeat,
			Score: &TopicScoreConfig{
				TimeInMeshWeight:            0.1,
				FirstMessageDeliveriesWeight: 0.5,
			},
		},
	}
}

type MessageValidator func(msg *Message) bool

func NoOpValidator(msg *Message) bool {
	return len(msg.Data) > 0
}

func ProviderValidator(msg *Message) bool {
	return len(msg.Data) > 0
}

func RequestValidator(msg *Message) bool {
	return len(msg.Data) > 0 && len(msg.Data) < 1024*1024
}

func HeartbeatValidator(msg *Message) bool {
	return len(msg.Data) > 0 && len(msg.Data) < 1024
}
