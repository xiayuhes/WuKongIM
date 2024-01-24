package clusterconfig

import (
	"io"
	"os"
	"path"
	"sync"

	"github.com/WuKongIM/WuKongIM/pkg/cluster/cluster/clusterconfig/pb"
	"github.com/WuKongIM/WuKongIM/pkg/wklog"
	"github.com/WuKongIM/WuKongIM/pkg/wkutil"
	"go.uber.org/zap"
)

type ConfigManager struct {
	sync.Mutex
	cfg     *pb.Config
	cfgFile *os.File
	opts    *Options
	wklog.Log
}

func NewConfigManager(opts *Options) *ConfigManager {

	cm := &ConfigManager{
		cfg:  &pb.Config{},
		opts: opts,
		Log:  wklog.NewWKLog("ConfigManager"),
	}

	configDir := path.Dir(opts.ConfigPath)
	if configDir != "" {
		err := os.MkdirAll(configDir, os.ModePerm)
		if err != nil {
			cm.Panic("create config dir error", zap.Error(err))
		}
	}
	err := cm.initConfigFromFile()
	if err != nil {
		cm.Panic("init cluster config from file error", zap.Error(err))
	}

	opts.AppliedConfigVersion = cm.cfg.Version

	return cm
}

func (c *ConfigManager) GetConfig() *pb.Config {
	return c.cfg
}

func (c *ConfigManager) AddOrUpdateNodes(nodes []*pb.Node, cfg *pb.Config) {

	for i, node := range nodes {
		if c.existNodeByCfg(node.Id, cfg) {
			cfg.Nodes[i] = node
			continue
		}
		cfg.Nodes = append(cfg.Nodes, node)
	}
}

func (c *ConfigManager) AddOrUpdateSlots(slots []*pb.Slot, cfg *pb.Config) {
	for i, slot := range slots {
		if c.existSlotByCfg(slot.Id, cfg) {
			cfg.Slots[i] = slot
			continue
		}
		cfg.Slots = append(cfg.Slots, slot)
	}
}

func (c *ConfigManager) Close() {
	c.cfgFile.Close()
}

func (c *ConfigManager) Version() uint64 {
	c.Lock()
	defer c.Unlock()
	return c.cfg.Version
}

func (c *ConfigManager) existNode(nodeId uint64) bool {
	for _, node := range c.cfg.Nodes {
		if node.Id == nodeId {
			return true
		}
	}
	return false
}

func (c *ConfigManager) existNodeByCfg(nodeId uint64, cfg *pb.Config) bool {
	for _, node := range cfg.Nodes {
		if node.Id == nodeId {
			return true
		}
	}
	return false
}

func (c *ConfigManager) existSlot(slotId uint32) bool {
	for _, slot := range c.cfg.Slots {
		if slot.Id == slotId {
			return true
		}
	}
	return false
}

func (c *ConfigManager) existSlotByCfg(slotId uint32, cfg *pb.Config) bool {
	for _, slot := range cfg.Slots {
		if slot.Id == slotId {
			return true
		}
	}
	return false
}

func (c *ConfigManager) saveConfig() error {
	data := c.getConfigData()
	if _, err := c.cfgFile.WriteAt(data, 0); err != nil {
		return err
	}
	return nil
}

func (c *ConfigManager) SaveConfig() error {
	c.Lock()
	defer c.Unlock()
	return c.saveConfig()
}

func (c *ConfigManager) UpdateConfig(cfg *pb.Config) error {
	c.Lock()
	defer c.Unlock()
	c.cfg = cfg
	return c.saveConfig()
}

func (c *ConfigManager) GetConfigData() []byte {
	c.Lock()
	defer c.Unlock()
	return c.getConfigData()
}

func (c *ConfigManager) getConfigData() []byte {
	return []byte(wkutil.ToJSON(c.cfg))
}

func (c *ConfigManager) GetConfigDataByCfg(cfg *pb.Config) []byte {
	return []byte(wkutil.ToJSON(cfg))
}

func (c *ConfigManager) UnmarshalConfigData(data []byte, cfg *pb.Config) error {
	c.Lock()
	defer c.Unlock()
	return wkutil.ReadJSONByByte(data, cfg)
}

func (c *ConfigManager) initConfigFromFile() error {
	clusterCfgPath := c.opts.ConfigPath
	var err error
	c.cfgFile, err = os.OpenFile(clusterCfgPath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		c.Panic("Open cluster config file failed!", zap.Error(err))
	}

	data, err := io.ReadAll(c.cfgFile)
	if err != nil {
		c.Panic("Read cluster config file failed!", zap.Error(err))
	}
	if len(data) > 0 {
		if err := wkutil.ReadJSONByByte(data, c.cfg); err != nil {
			c.Panic("Unmarshal cluster config failed!", zap.Error(err))
		}
	}
	return nil
}
