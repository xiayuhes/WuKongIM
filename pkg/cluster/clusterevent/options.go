package clusterevent

import (
	"time"

	"github.com/WuKongIM/WuKongIM/pkg/cluster/icluster"
	"github.com/WuKongIM/WuKongIM/pkg/cluster/reactor"
)

type Options struct {
	NodeId                 uint64
	InitNodes              map[uint64]string
	SlotCount              uint32 // 槽位数量
	SlotMaxReplicaCount    uint32 // 每个槽位最大副本数量
	ChannelMaxReplicaCount uint32 // 每个频道最大副本数量
	ConfigDir              string
	ApiServerAddr          string // api服务地址
	Ready                  func(msgs []Message)
	Send                   func(m reactor.Message) // 发送消息
	// PongMaxTick 节点超过多少tick没有回应心跳就认为是掉线
	PongMaxTick int
	// 学习者检查间隔（每隔这个间隔时间检查下学习者的日志）
	LearnerCheckInterval time.Duration

	Cluster icluster.Cluster // 分布式接口
}

func NewOptions(opt ...Option) *Options {
	opts := &Options{
		SlotCount:              128,
		SlotMaxReplicaCount:    3,
		ChannelMaxReplicaCount: 3,
		ConfigDir:              "clusterconfig",
		PongMaxTick:            30,
		LearnerCheckInterval:   time.Second * 2,
	}
	for _, o := range opt {
		o(opts)
	}
	return opts
}

type Option func(opts *Options)

func WithNodeId(nodeId uint64) Option {
	return func(o *Options) {
		o.NodeId = nodeId
	}
}

func WithInitNodes(initNodes map[uint64]string) Option {
	return func(o *Options) {
		o.InitNodes = initNodes
	}

}
func WithSlotCount(slotCount uint32) Option {
	return func(o *Options) {
		o.SlotCount = slotCount
	}
}
func WithSlotMaxReplicaCount(slotMaxReplicaCount uint32) Option {
	return func(o *Options) {
		o.SlotMaxReplicaCount = slotMaxReplicaCount
	}
}

func WithChannelMaxReplicaCount(channelMaxReplicaCount uint32) Option {
	return func(o *Options) {
		o.ChannelMaxReplicaCount = channelMaxReplicaCount
	}

}

func WithReady(f func(msgs []Message)) Option {

	return func(o *Options) {
		o.Ready = f
	}
}

func WithSend(f func(m reactor.Message)) Option {
	return func(o *Options) {
		o.Send = f
	}
}

func WithConfigDir(configDir string) Option {
	return func(o *Options) {
		o.ConfigDir = configDir
	}
}

func WithApiServerAddr(apiServerAddr string) Option {
	return func(o *Options) {
		o.ApiServerAddr = apiServerAddr
	}
}

func WithCluster(cluster icluster.Cluster) Option {
	return func(o *Options) {
		o.Cluster = cluster
	}
}