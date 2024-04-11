package clusterstore

import (
	"github.com/WuKongIM/WuKongIM/pkg/wkdb"
)

// AddSubscribers 添加订阅者
func (s *Store) AddSubscribers(channelID string, channelType uint8, subscribers []string) error {
	data := EncodeSubscribers(channelID, channelType, subscribers)
	cmd := NewCMD(CMDAddSubscribers, data)
	cmdData, err := cmd.Marshal()
	if err != nil {
		return err
	}
	_, err = s.opts.Cluster.ProposeChannelMeta(s.ctx, channelID, channelType, cmdData)
	return err
}

// RemoveSubscribers 移除订阅者
func (s *Store) RemoveSubscribers(channelID string, channelType uint8, subscribers []string) error {
	data := EncodeSubscribers(channelID, channelType, subscribers)
	cmd := NewCMD(CMDRemoveSubscribers, data)
	cmdData, err := cmd.Marshal()
	if err != nil {
		return err
	}
	_, err = s.opts.Cluster.ProposeChannelMeta(s.ctx, channelID, channelType, cmdData)
	return err
}

func (s *Store) RemoveAllSubscriber(channelId string, channelType uint8) error {
	data := EncodeChannel(channelId, channelType)
	cmd := NewCMD(CMDRemoveAllSubscriber, data)
	cmdData, err := cmd.Marshal()
	if err != nil {
		return err
	}
	_, err = s.opts.Cluster.ProposeChannelMeta(s.ctx, channelId, channelType, cmdData)
	return err
}

func (s *Store) GetSubscribers(channelID string, channelType uint8) ([]string, error) {
	return s.wdb.GetSubscribers(channelID, channelType)
}

// AddOrUpdateChannel add or update channel
func (s *Store) AddOrUpdateChannel(channelInfo wkdb.ChannelInfo) error {
	data, err := EncodeAddOrUpdateChannel(channelInfo)
	if err != nil {
		return err
	}
	cmd := NewCMD(CMDAddOrUpdateChannel, data)
	cmdData, err := cmd.Marshal()
	if err != nil {
		return err
	}
	_, err = s.opts.Cluster.ProposeChannelMeta(s.ctx, channelInfo.ChannelId, channelInfo.ChannelType, cmdData)
	return err
}

func (s *Store) DeleteChannel(channelId string, channelType uint8) error {
	data := EncodeChannel(channelId, channelType)
	cmd := NewCMD(CMDDeleteChannel, data)
	cmdData, err := cmd.Marshal()
	if err != nil {
		return err
	}
	_, err = s.opts.Cluster.ProposeChannelMeta(s.ctx, channelId, channelType, cmdData)
	return err
}

func (s *Store) GetChannel(channelID string, channelType uint8) (wkdb.ChannelInfo, error) {
	return s.wdb.GetChannel(channelID, channelType)
}

func (s *Store) ExistChannel(channelID string, channelType uint8) (bool, error) {
	return s.wdb.ExistChannel(channelID, channelType)
}

func (s *Store) AddDenylist(channelID string, channelType uint8, uids []string) error {
	data := EncodeSubscribers(channelID, channelType, uids)
	cmd := NewCMD(CMDAddDenylist, data)
	cmdData, err := cmd.Marshal()
	if err != nil {
		return err
	}
	_, err = s.opts.Cluster.ProposeChannelMeta(s.ctx, channelID, channelType, cmdData)
	return err

}

func (s *Store) GetDenylist(channelID string, channelType uint8) ([]string, error) {
	return s.wdb.GetDenylist(channelID, channelType)
}

func (s *Store) RemoveAllDenylist(channelID string, channelType uint8) error {
	cmd := NewCMD(CMDRemoveAllDenylist, nil)
	cmdData, err := cmd.Marshal()
	if err != nil {
		return err
	}
	_, err = s.opts.Cluster.ProposeChannelMeta(s.ctx, channelID, channelType, cmdData)
	return err
}

func (s *Store) RemoveDenylist(channelID string, channelType uint8, uids []string) error {
	data := EncodeSubscribers(channelID, channelType, uids)
	cmd := NewCMD(CMDRemoveDenylist, data)
	cmdData, err := cmd.Marshal()
	if err != nil {
		return err
	}
	_, err = s.opts.Cluster.ProposeChannelMeta(s.ctx, channelID, channelType, cmdData)
	return err
}

func (s *Store) AddAllowlist(channelID string, channelType uint8, uids []string) error {
	data := EncodeSubscribers(channelID, channelType, uids)
	cmd := NewCMD(CMDAddAllowlist, data)
	cmdData, err := cmd.Marshal()
	if err != nil {
		return err
	}
	_, err = s.opts.Cluster.ProposeChannelMeta(s.ctx, channelID, channelType, cmdData)
	return err
}

func (s *Store) GetAllowlist(channelID string, channelType uint8) ([]string, error) {
	return s.wdb.GetAllowlist(channelID, channelType)
}

func (s *Store) RemoveAllAllowlist(channelID string, channelType uint8) error {
	cmd := NewCMD(CMDRemoveAllAllowlist, nil)
	cmdData, err := cmd.Marshal()
	if err != nil {
		return err
	}
	_, err = s.opts.Cluster.ProposeChannelMeta(s.ctx, channelID, channelType, cmdData)
	return err
}

func (s *Store) RemoveAllowlist(channelID string, channelType uint8, uids []string) error {
	data := EncodeSubscribers(channelID, channelType, uids)
	cmd := NewCMD(CMDRemoveAllowlist, data)
	cmdData, err := cmd.Marshal()
	if err != nil {
		return err
	}
	_, err = s.opts.Cluster.ProposeChannelMeta(s.ctx, channelID, channelType, cmdData)
	return err
}

// func (s *Store) DeleteChannelClusterConfig(channelID string, channelType uint8) error {
// 	cmd := NewCMD(CMDChannelClusterConfigDelete, nil)
// 	cmdData, err := cmd.Marshal()
// 	if err != nil {
// 		return err
// 	}
// 	_, err = s.opts.Cluster.ProposeChannelMeta(s.ctx, channelID, channelType, cmdData)
// 	return err
// }