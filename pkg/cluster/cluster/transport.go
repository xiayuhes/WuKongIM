package cluster

import (
	replica "github.com/WuKongIM/WuKongIM/pkg/cluster/replica2"
	"github.com/WuKongIM/WuKongIM/pkg/wkserver/proto"
	wkproto "github.com/WuKongIM/WuKongIMGoProto"
)

type ITransport interface {
	// Send 发送消息
	Send(to uint64, m *proto.Message) error
	// OnMessage 收取消息
	OnMessage(f func(from uint64, m *proto.Message))
}

func NewMessage(shardNo string, msg replica.Message) (*proto.Message, error) {
	m := Message{
		ShardNo: shardNo,
		Message: msg,
	}
	data, err := m.Marshal()
	if err != nil {
		return nil, err
	}
	return &proto.Message{
		MsgType: MsgShardMsg,
		Content: data,
	}, nil

}

func NewMessageFromProto(m *proto.Message) (Message, error) {
	return UnmarshalMessage(m.Content)
}

type Message struct {
	ShardNo string
	replica.Message
}

func (m Message) Marshal() ([]byte, error) {
	enc := wkproto.NewEncoder()
	defer enc.End()

	enc.WriteString(m.ShardNo)
	msgData, err := m.Message.Marshal()
	if err != nil {
		return nil, err
	}
	enc.WriteBytes(msgData)
	return enc.Bytes(), nil
}

func UnmarshalMessage(data []byte) (Message, error) {
	m := Message{}
	dec := wkproto.NewDecoder(data)
	var err error
	if m.ShardNo, err = dec.String(); err != nil {
		return m, err
	}
	msgData, err := dec.BinaryAll()
	if err != nil {
		return m, err
	}
	m.Message, err = replica.UnmarshalMessage(msgData)
	if err != nil {
		return m, err
	}
	return m, nil
}

type MemoryTransport struct {
	nodeMessageListenerMap map[uint64]func(m *proto.Message)
}

func NewMemoryTransport() *MemoryTransport {
	return &MemoryTransport{
		nodeMessageListenerMap: make(map[uint64]func(m *proto.Message)),
	}
}

func (t *MemoryTransport) Send(to uint64, m *proto.Message) error {
	if f, ok := t.nodeMessageListenerMap[to]; ok {
		go f(m) // 模拟网络请求
	}
	return nil
}

func (t *MemoryTransport) OnMessage(f func(from uint64, m *proto.Message)) {

}

func (t *MemoryTransport) OnNodeMessage(nodeID uint64, f func(m *proto.Message)) {
	t.nodeMessageListenerMap[nodeID] = f
}