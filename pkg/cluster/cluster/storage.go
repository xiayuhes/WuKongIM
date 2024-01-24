package cluster

import (
	"path"

	replica "github.com/WuKongIM/WuKongIM/pkg/cluster/replica2"
	"github.com/cockroachdb/pebble"
)

// 日志分区存储
type IShardLogStorage interface {
	// AppendLog 追加日志
	AppendLog(shardNo string, logs []replica.Log) error
	// 截断日志
	TruncateLogTo(shardNo string, index uint64) error
	// 获取日志
	Logs(shardNo string, startLogIndex uint64, endLogIndex uint64, limit uint32) ([]replica.Log, error)
	// 最后一条日志的索引
	LastIndex(shardNo string) (uint64, error)
	// 设置成功被状态机应用的日志索引
	SetAppliedIndex(shardNo string, index uint64) error
	//	 获取最后一条日志的索引和追加时间
	LastIndexAndAppendTime(shardNo string) (uint64, uint64, error)

	// SetLeaderTermStartIndex 设置领导任期开始的第一条日志索引
	SetLeaderTermStartIndex(shardNo string, term uint32, index uint64) error
	// LeaderLastTerm 获取最新的本地保存的领导任期
	LeaderLastTerm(shardNo string) (uint32, error)
	// LeaderTermStartIndex 获取领导任期开始的第一条日志索引
	LeaderTermStartIndex(shardNo string, term uint32) (uint64, error)
	// 删除比传入的term大的的LeaderTermStartIndex记录
	DeleteLeaderTermStartIndexGreaterThanTerm(shardNo string, term uint32) error
}

type MemoryShardLogStorage struct {
	storage                 map[string][]replica.Log
	leaderTermStartIndexMap map[string]map[uint32]uint64
}

func NewMemoryShardLogStorage() *MemoryShardLogStorage {
	return &MemoryShardLogStorage{
		storage:                 make(map[string][]replica.Log),
		leaderTermStartIndexMap: make(map[string]map[uint32]uint64),
	}
}

func (m *MemoryShardLogStorage) AppendLog(shardNo string, logs []replica.Log) error {
	m.storage[shardNo] = append(m.storage[shardNo], logs...)
	return nil
}

func (m *MemoryShardLogStorage) TruncateLogTo(shardNo string, index uint64) error {
	logs := m.storage[shardNo]
	if len(logs) > 0 {
		m.storage[shardNo] = logs[:index]
	}
	return nil
}

func (m *MemoryShardLogStorage) Logs(shardNo string, startLogIndex uint64, endLogIndex uint64, limit uint32) ([]replica.Log, error) {
	logs := m.storage[shardNo]
	if len(logs) == 0 {
		return nil, nil
	}
	if endLogIndex == 0 {
		return logs[startLogIndex-1:], nil
	}
	if startLogIndex > uint64(len(logs)) {
		return nil, nil
	}
	if endLogIndex > uint64(len(logs)) {
		return logs[startLogIndex-1:], nil
	}
	return logs[startLogIndex-1 : endLogIndex-1], nil
}

func (m *MemoryShardLogStorage) LastIndex(shardNo string) (uint64, error) {
	logs := m.storage[shardNo]
	if len(logs) == 0 {
		return 0, nil
	}
	return uint64(len(logs) - 1), nil
}

func (m *MemoryShardLogStorage) SetAppliedIndex(shardNo string, index uint64) error {
	return nil
}

func (m *MemoryShardLogStorage) LastIndexAndAppendTime(shardNo string) (uint64, uint64, error) {
	return 0, 0, nil
}

func (m *MemoryShardLogStorage) SetLeaderTermStartIndex(shardNo string, term uint32, index uint64) error {
	if _, ok := m.leaderTermStartIndexMap[shardNo]; !ok {
		m.leaderTermStartIndexMap[shardNo] = make(map[uint32]uint64)
	}
	m.leaderTermStartIndexMap[shardNo][term] = index
	return nil
}

func (m *MemoryShardLogStorage) LeaderLastTerm(shardNo string) (uint32, error) {
	if _, ok := m.leaderTermStartIndexMap[shardNo]; !ok {
		return 0, nil
	}
	var maxTerm uint32
	for term := range m.leaderTermStartIndexMap[shardNo] {
		if term > maxTerm {
			maxTerm = term
		}
	}
	return maxTerm, nil
}

func (m *MemoryShardLogStorage) LeaderTermStartIndex(shardNo string, term uint32) (uint64, error) {
	if _, ok := m.leaderTermStartIndexMap[shardNo]; !ok {
		return 0, nil
	}
	return m.leaderTermStartIndexMap[shardNo][term], nil
}

func (m *MemoryShardLogStorage) DeleteLeaderTermStartIndexGreaterThanTerm(shardNo string, term uint32) error {
	if _, ok := m.leaderTermStartIndexMap[shardNo]; !ok {
		return nil
	}
	for t := range m.leaderTermStartIndexMap[shardNo] {
		if t > term {
			delete(m.leaderTermStartIndexMap[shardNo], t)
		}
	}
	return nil
}

type proxyReplicaStorage struct {
	storage IShardLogStorage
	shardNo string
}

func newProxyReplicaStorage(shardNo string, storage IShardLogStorage) *proxyReplicaStorage {
	return &proxyReplicaStorage{
		storage: storage,
		shardNo: shardNo,
	}
}

func (p *proxyReplicaStorage) AppendLog(logs []replica.Log) error {
	return p.storage.AppendLog(p.shardNo, logs)
}

func (p *proxyReplicaStorage) TruncateLogTo(index uint64) error {
	return p.storage.TruncateLogTo(p.shardNo, index)
}

func (p *proxyReplicaStorage) Logs(startLogIndex uint64, endLogIndex uint64, limit uint32) ([]replica.Log, error) {
	return p.storage.Logs(p.shardNo, startLogIndex, endLogIndex, limit)
}

func (p *proxyReplicaStorage) LastIndex() (uint64, error) {
	return p.storage.LastIndex(p.shardNo)
}

func (p *proxyReplicaStorage) SetAppliedIndex(index uint64) error {
	return p.storage.SetAppliedIndex(p.shardNo, index)
}

func (p *proxyReplicaStorage) LastIndexAndAppendTime() (uint64, uint64, error) {
	return p.storage.LastIndexAndAppendTime(p.shardNo)
}

func (p *proxyReplicaStorage) SetLeaderTermStartIndex(term uint32, index uint64) error {
	return p.storage.SetLeaderTermStartIndex(p.shardNo, term, index)
}

func (p *proxyReplicaStorage) LeaderLastTerm() (uint32, error) {
	return p.storage.LeaderLastTerm(p.shardNo)
}

func (p *proxyReplicaStorage) LeaderTermStartIndex(term uint32) (uint64, error) {
	return p.storage.LeaderTermStartIndex(p.shardNo, term)
}

func (p *proxyReplicaStorage) DeleteLeaderTermStartIndexGreaterThanTerm(term uint32) error {
	return p.storage.DeleteLeaderTermStartIndexGreaterThanTerm(p.shardNo, term)
}

type localStorage struct {
	db    *pebble.DB
	opts  *Options
	dbDir string
}

func newLocalStorage(opts *Options) *localStorage {
	dbDir := path.Join(opts.DataDir, "wukongimdb")
	return &localStorage{
		opts:  opts,
		dbDir: dbDir,
	}
}

func (l *localStorage) open() error {
	var err error
	l.db, err = pebble.Open(l.dbDir, &pebble.Options{})
	return err
}

func (l *localStorage) close() error {
	return l.db.Close()
}

func (l *localStorage) saveChannelClusterInfo(channelID string, channelType uint8, clusterInfo *ChannelClusterConfig) error {
	return nil
}