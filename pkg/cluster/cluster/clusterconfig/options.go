package clusterconfig

import "time"

type Options struct {
	NodeId               uint64
	ConfigPath           string        // 集群配置文件路径
	ElectionTimeoutTick  int           // 选举超时tick次数
	HeartbeatTimeoutTick int           // 心跳超时tick次数
	Replicas             []uint64      // 副本列表 (必须包含自己本身的id)
	Transport            ITransport    // 传输层
	AppliedConfigVersion uint64        // 已应用的配置版本
	ProposeTimeout       time.Duration // 提议超时时间
}

func NewOptions() *Options {
	return &Options{
		ConfigPath:           "clusterconfig.json",
		ElectionTimeoutTick:  10,
		HeartbeatTimeoutTick: 1,
		ProposeTimeout:       time.Second * 5,
	}
}

type Option func(opts *Options)

func WithNodeId(nodeId uint64) Option {
	return func(opts *Options) {
		opts.NodeId = nodeId
	}
}

func WithConfigPath(configPath string) Option {
	return func(opts *Options) {
		opts.ConfigPath = configPath
	}
}

func WithElectionTimeoutTick(electionTimeoutTick int) Option {
	return func(opts *Options) {
		opts.ElectionTimeoutTick = electionTimeoutTick
	}
}

func WithHeartbeatTimeoutTick(heartbeatTimeoutTick int) Option {
	return func(opts *Options) {
		opts.HeartbeatTimeoutTick = heartbeatTimeoutTick
	}
}

func WithReplicas(replicas []uint64) Option {
	return func(opts *Options) {
		opts.Replicas = replicas
	}
}

func WithTransport(transport ITransport) Option {
	return func(opts *Options) {
		opts.Transport = transport
	}
}
