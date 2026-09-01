package canal

import (
	"time"

	"github.com/withlin/canal-go/client"
)

func NewSyncConnector(opts ...SyncConnectorOptionFunc) (*SyncConnector, error) {
	option := NewSyncConnectorOption(opts...)
	sc := &SyncConnector{
		SyncConnectorOption: *option,
		CanalConnector: client.NewSimpleCanalConnector(
			option.CanalIp,
			option.CanalPort,
			option.CanalUser,
			option.CanalPwd,
			option.CanalInstance,
			option.CanalSoTimeout,
			option.CanalIdleTimeout,
		),
	}
	return sc, nil
}

type SyncConnectorOptionFunc func(*SyncConnectorOption)

func NewSyncConnectorOption(opts ...SyncConnectorOptionFunc) *SyncConnectorOption {
	option := &SyncConnectorOption{
		Index:            "",
		CanalIp:          "",
		CanalPort:        11111,
		CanalUser:        "admin",
		CanalPwd:         "",
		CanalInstance:    "",
		CanalSoTimeout:   60000,
		CanalIdleTimeout: 60 * 60 * 1000,
		TableFilter:      ".*\\..*",
		Logger:           defaultLogger{},
	}
	return option.WithOption(opts...)
}

// WithOption 应用一组连接器配置选项。
func (i *SyncConnectorOption) WithOption(opts ...SyncConnectorOptionFunc) *SyncConnectorOption {
	for _, opt := range opts {
		if opt != nil {
			opt(i)
		}
	}
	return i
}

func WithSyncCanalIp(v string) SyncConnectorOptionFunc {
	return func(o *SyncConnectorOption) { o.CanalIp = v }
}
func WithSyncCanalPort(v int) SyncConnectorOptionFunc {
	return func(o *SyncConnectorOption) { o.CanalPort = v }
}
func WithSyncCanalInstance(v string) SyncConnectorOptionFunc {
	return func(o *SyncConnectorOption) { o.CanalInstance = v }
}
func WithSyncCanalUser(v string) SyncConnectorOptionFunc {
	return func(o *SyncConnectorOption) { o.CanalUser = v }
}
func WithSyncCanalPwd(v string) SyncConnectorOptionFunc {
	return func(o *SyncConnectorOption) { o.CanalPwd = v }
}
func WithSyncIndex(v string) SyncConnectorOptionFunc {
	return func(o *SyncConnectorOption) { o.Index = v }
}
func WithSyncTableFilter(v string) SyncConnectorOptionFunc {
	return func(o *SyncConnectorOption) { o.TableFilter = v }
}
func WithSyncBatch(v int32) SyncConnectorOptionFunc {
	return func(o *SyncConnectorOption) { o.Batch = v }
}
func WithSyncTimeout(v time.Duration) SyncConnectorOptionFunc {
	return func(o *SyncConnectorOption) { o.Timeout = v }
}
func WithSyncOuter(v IOuter) SyncConnectorOptionFunc {
	return func(o *SyncConnectorOption) { o.Outer = v }
}
func WithSyncLogger(v ILogger) SyncConnectorOptionFunc {
	return func(o *SyncConnectorOption) { o.Logger = v }
}

type SyncConnectorOption struct {
	CanalIp          string        `json:"canal_ip" yaml:"canal_ip"`
	CanalPort        int           `json:"canal_port" yaml:"canal_port"`
	CanalInstance    string        `json:"canal_instance" yaml:"canal_instance"`
	CanalUser        string        `json:"canal_user" yaml:"canal_user"`
	CanalPwd         string        `json:"canal_pwd" yaml:"canal_pwd"`
	CanalSoTimeout   int32         `json:"canal_so_timeout" yaml:"canal_so_timeout"`
	CanalIdleTimeout int32         `json:"canal_idle_timeout" yaml:"canal_idle_timeout"`
	Index            string        `json:"index" yaml:"index"`
	TableFilter      string        `json:"table_filter" yaml:"table_filter"`
	Batch            int32         `json:"batch" yaml:"batch"`
	Timeout          time.Duration `json:"timeout" yaml:"timeout"`
	Outer            IOuter        `json:"-" yaml:"-"`
	Logger           ILogger       `json:"-" yaml:"-"`
}
