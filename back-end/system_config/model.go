package system_config

type JSONData struct {
	Brokers []struct {
		Address string `json:"address"`
	} `json:"brokers"`
	Agents []AgentConfig `json:"agents"`
}

type AgentConfig struct {
	Name                 string `json:"name"`
	Address              string `json:"address"`
	BrokerPartition      int32  `json:"broker_partition"`
	ConsumerGroup        string `json:"consumer_group"`
	BrokerCommitInterval uint   `json:"broker_commit_interval"`
	HttpPort             int    `json:"http_port"`
	GrpcPort             uint   `json:"grpc_port"`
}
