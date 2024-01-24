package clusterstore

import (
	"fmt"
	"strings"

	"github.com/WuKongIM/WuKongIM/pkg/cluster"
	"github.com/WuKongIM/WuKongIM/pkg/cluster/replica"
	"github.com/WuKongIM/WuKongIM/pkg/wklog"
	"github.com/WuKongIM/WuKongIM/pkg/wkstore"
	wkproto "github.com/WuKongIM/WuKongIMGoProto"
	"go.uber.org/zap"
)

func (s *Store) AppendMessages(channelID string, channelType uint8, msgs []wkstore.Message) error {
	if len(msgs) == 0 {
		return nil
	}
	var msgData [][]byte
	for _, msg := range msgs {
		msgData = append(msgData, msg.Encode())
	}

	lastIndexs, err := s.opts.Cluster.ProposeMessagesToChannel(channelID, channelType, msgData)
	if err != nil {
		return err
	}

	fmt.Println("lastIndexs------>", lastIndexs)
	for idx, msg := range msgs {
		msg.SetSeq(uint32(lastIndexs[idx]))
	}
	return nil
}

func (s *Store) LoadNextRangeMsgs(channelID string, channelType uint8, startMessageSeq, endMessageSeq uint32, limit int) ([]wkstore.Message, error) {
	return s.db.LoadNextRangeMsgs(channelID, channelType, startMessageSeq, endMessageSeq, limit)
}

func (s *Store) LoadMsg(channelID string, channelType uint8, seq uint32) (wkstore.Message, error) {
	return s.db.LoadMsg(channelID, channelType, seq)
}

func (s *Store) LoadLastMsgs(channelID string, channelType uint8, limit int) ([]wkstore.Message, error) {
	return s.db.LoadLastMsgs(channelID, channelType, limit)
}

func (s *Store) LoadLastMsgsWithEnd(channelID string, channelType uint8, end uint32, limit int) ([]wkstore.Message, error) {
	return s.db.LoadLastMsgsWithEnd(channelID, channelType, end, limit)
}

func (s *Store) LoadPrevRangeMsgs(channelID string, channelType uint8, start, end uint32, limit int) ([]wkstore.Message, error) {
	return s.db.LoadPrevRangeMsgs(channelID, channelType, start, end, limit)
}

func (s *Store) GetLastMsgSeq(channelID string, channelType uint8) (uint32, error) {
	return s.db.GetLastMsgSeq(channelID, channelType)
}

func (s *Store) GetMessagesOfNotifyQueue(count int) ([]wkstore.Message, error) {
	return s.db.GetMessagesOfNotifyQueue(count)
}

func (s *Store) AppendMessageOfNotifyQueue(messages []wkstore.Message) error {
	return s.db.AppendMessageOfNotifyQueue(messages)
}

func (s *Store) RemoveMessagesOfNotifyQueue(messageIDs []int64) error {
	return s.db.RemoveMessagesOfNotifyQueue(messageIDs)
}

func (s *Store) GetMessageShardLogStorage() *MessageShardLogStorage {
	return s.messageShardLogStorage
}

// SaveStreamMeta 保存消息流元数据
func (s *Store) SaveStreamMeta(meta *wkstore.StreamMeta) error {
	return nil
}

// StreamEnd 结束流
func (s *Store) StreamEnd(channelID string, channelType uint8, streamNo string) error {
	return nil
}

func (s *Store) GetStreamMeta(channelID string, channelType uint8, streamNo string) (*wkstore.StreamMeta, error) {
	return nil, nil
}

func (s *Store) GetStreamItems(channelID string, channelType uint8, streamNo string) ([]*wkstore.StreamItem, error) {
	return nil, nil
}

// AppendStreamItem 追加消息流
func (s *Store) AppendStreamItem(channelID string, channelType uint8, streamNo string, item *wkstore.StreamItem) (uint32, error) {
	return 0, nil
}

func (s *Store) UpdateMessageOfUserCursorIfNeed(uid string, messageSeq uint32) error {
	return nil
}

func (s *Store) GetMessageOfUserCursor(uid string) (uint32, error) {
	return 0, nil
}

func (s *Store) AppendMessagesOfUser(uid string, msgs []wkstore.Message) (err error) {
	return s.AppendMessages(uid, wkproto.ChannelTypePerson, msgs)
}

func (s *Store) SyncMessageOfUser(uid string, messageSeq uint32, limit uint32) ([]wkstore.Message, error) {
	return nil, nil
}

func (s *Store) DeleteChannelAndClearMessages(channelID string, channelType uint8) error {
	return nil
}

type MessageShardLogStorage struct {
	db *wkstore.FileStore
	wklog.Log
}

func NewMessageShardLogStorage(db *wkstore.FileStore) *MessageShardLogStorage {
	return &MessageShardLogStorage{
		db:  db,
		Log: wklog.NewWKLog("MessageShardLogStorage"),
	}
}

func (s *MessageShardLogStorage) AppendLog(shardNo string, log replica.Log) error {
	channelID, channelType := cluster.GetChannelFromChannelKey(shardNo)
	s.Debug("AppendLog", zap.String("channelID", channelID), zap.Uint8("channelType", channelType), zap.Uint64("index", log.Index), zap.ByteString("data", log.Data))
	lastIndex, err := s.LastIndex(shardNo)
	if err != nil {
		s.Error("get last index err", zap.Error(err))
		return err
	}
	if log.Index <= lastIndex {
		s.Warn("log index is less than last index", zap.Uint64("logIndex", log.Index), zap.Uint64("lastIndex", lastIndex))
		return nil
	}
	msg, err := s.db.StoreConfig().DecodeMessageFnc(log.Data)
	if err != nil {
		s.Error("decode message err", zap.Error(err))
		return err
	}
	msg.SetSeq(uint32(log.Index))
	msg.SetTerm(uint64(log.Term))
	if IsUserOwnChannel(channelID, channelType) {
		_, err = s.db.AppendMessagesOfUser(channelID, []wkstore.Message{msg})
		if err != nil {
			s.Error("AppendMessagesOfUser err", zap.Error(err))
			return err
		}
	} else {
		_, err = s.db.AppendMessages(channelID, channelType, []wkstore.Message{msg})
		if err != nil {
			s.Error("AppendMessages err", zap.Error(err))
			return err
		}
	}

	return s.db.SaveChannelMaxMessageSeq(channelID, channelType, uint32(log.Index))
}

// 获取日志
func (s *MessageShardLogStorage) GetLogs(shardNo string, startLogIndex uint64, limit uint32) ([]replica.Log, error) {

	channelID, channelType := cluster.GetChannelFromChannelKey(shardNo)
	if startLogIndex > 0 {
		startLogIndex = startLogIndex - 1
	}
	var (
		messages []wkstore.Message
		err      error
	)
	if IsUserOwnChannel(channelID, channelType) {
		messages, err = s.db.SyncMessageOfUser(channelID, uint32(startLogIndex), int(limit))
		if err != nil {
			return nil, err
		}
	} else {
		messages, err = s.db.LoadNextRangeMsgs(channelID, channelType, uint32(startLogIndex), 0, int(limit))
		if err != nil {
			return nil, err
		}
	}

	if len(messages) == 0 {
		return nil, nil
	}
	logs := make([]replica.Log, len(messages))
	for i, msg := range messages {
		logs[i] = replica.Log{
			Index: uint64(msg.GetSeq()),
			Term:  uint32(msg.GetTerm()),
			Data:  msg.Encode(),
		}
	}
	return logs, nil
}

// 最后一条日志的索引
func (s *MessageShardLogStorage) LastIndex(shardNo string) (uint64, error) {
	channelID, channelType := cluster.GetChannelFromChannelKey(shardNo)

	var (
		msgSeq uint32
		err    error
	)
	if IsUserOwnChannel(channelID, channelType) {
		msgSeq, err = s.db.GetLastMsgSeqOfUser(channelID)
	} else {
		msgSeq, err = s.db.GetLastMsgSeq(channelID, channelType)
	}
	return uint64(msgSeq), err

}

// 获取第一条日志的索引
func (s *MessageShardLogStorage) FirstIndex(shardNo string) (uint64, error) {
	return 0, nil
}

// 设置成功被状态机应用的日志索引
func (s *MessageShardLogStorage) SetAppliedIndex(shardNo string, index uint64) error {
	return nil
}

func (s *MessageShardLogStorage) LastIndexAndAppendTime(shardNo string) (uint64, uint64, error) {
	channelID, channelType := cluster.GetChannelFromChannelKey(shardNo)
	lastMsgSeq, lastAppendTime, err := s.db.GetChannelMaxMessageSeq(channelID, channelType)
	if err != nil {
		return 0, 0, err
	}
	return uint64(lastMsgSeq), lastAppendTime, nil
}

// 是用户自己的频道
func IsUserOwnChannel(channelID string, channelType uint8) bool {
	if channelType == wkproto.ChannelTypePerson {
		return !strings.Contains(channelID, "@")
	}
	return false
}