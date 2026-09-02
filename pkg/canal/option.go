package canal

import (
	"database/sql"
	"fmt"
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
		CanalIp:          "",
		CanalPort:        11111,
		CanalUser:        "admin",
		CanalPwd:         "",
		CanalInstance:    "",
		CanalSoTimeout:   60000,
		CanalIdleTimeout: 60 * 60 * 1000,
		BaseConnectorOption: BaseConnectorOption{
			Index:       "",
			TableFilter: ".*\\..*",
			Log:         DefaultLogger{},
		},
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
func WithSyncTableMap(v map[string]string) SyncConnectorOptionFunc {
	return func(o *SyncConnectorOption) { o.TableMap = v }
}
func WithSyncOuter(v IOuter) SyncConnectorOptionFunc {
	return func(o *SyncConnectorOption) { o.Outer = v }
}
func WithSyncLogger(v ILogger) SyncConnectorOptionFunc {
	return func(o *SyncConnectorOption) { o.Log = v }
}

type BaseConnectorOption struct {
	Outer       IOuter            `json:"-" yaml:"-"`
	Log         ILogger           `json:"-" yaml:"-"`
	Index       string            `json:"index" yaml:"index"`
	Batch       int32             `json:"batch" yaml:"batch"`
	Timeout     time.Duration     `json:"timeout" yaml:"timeout"`
	TableFilter string            `json:"table_filter" yaml:"table_filter"`
	TableMap    map[string]string `json:"table_map" yaml:"table_map"`
}
type SyncConnectorOption struct {
	BaseConnectorOption
	CanalIp          string `json:"canal_ip" yaml:"canal_ip"`
	CanalPort        int    `json:"canal_port" yaml:"canal_port"`
	CanalInstance    string `json:"canal_instance" yaml:"canal_instance"`
	CanalUser        string `json:"canal_user" yaml:"canal_user"`
	CanalPwd         string `json:"canal_pwd" yaml:"canal_pwd"`
	CanalSoTimeout   int32  `json:"canal_so_timeout" yaml:"canal_so_timeout"`
	CanalIdleTimeout int32  `json:"canal_idle_timeout" yaml:"canal_idle_timeout"`
}

func NewMigrateConnector(opts ...MigrateConnectorOptionFunc) (*MigrateConnector, error) {
	option := NewMigrateConnectorOption(opts...)
	sc := &MigrateConnector{
		MigrateConnectorOption: *option,
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=5s",
		option.DbUser,
		option.DbPwd,
		option.DbIp,
		option.DbPort,
		option.DbName,
	)
	client, err := sql.Open("mysql", dsn)
	if err != nil {
		return sc, err
	}
	sc.Client = client
	return sc, nil
}

type MigrateConnectorOption struct {
	BaseConnectorOption
	DbIp        string `json:"db_ip" yaml:"db_ip"`
	DbPort      int    `json:"db_port" yaml:"db_port"`
	DbInstance  string `json:"db_instance" yaml:"db_instance"`
	DbName      string `json:"db_name" yaml:"db_name"`
	DbUser      string `json:"db_user" yaml:"db_user"`
	DbPwd       string `json:"db_pwd" yaml:"db_pwd"`
	TableArgs   []int  `json:"table_args" yaml:"table_args"`
	CursorField string `json:"cursor_field" yaml:"cursor_field"`
	CursorType  int8   `json:"cursor_type" yaml:"cursor_type"` // 0- int 1- string
	CursorPos   string `json:"cursor_pos" yaml:"cursor_pos"`
}

// MigrateConnectorOptionFunc configures a MigrateConnectorOption.
type MigrateConnectorOptionFunc func(*MigrateConnectorOption)

// NewMigrateConnectorOption creates a migration connector option with defaults.
func NewMigrateConnectorOption(opts ...MigrateConnectorOptionFunc) *MigrateConnectorOption {
	option := &MigrateConnectorOption{
		DbPort:      3306,
		TableArgs:   []int{},
		CursorField: "id",
		CursorType:  0,
		CursorPos:   "0",
		BaseConnectorOption: BaseConnectorOption{
			Log:         DefaultLogger{},
			TableFilter: "%d",
		},
	}
	return option.WithOption(opts...)
}

// WithOption applies a set of migration connector options.
func (i *MigrateConnectorOption) WithOption(opts ...MigrateConnectorOptionFunc) *MigrateConnectorOption {
	for _, opt := range opts {
		if opt != nil {
			opt(i)
		}
	}
	return i
}

func WithMigrateDbIp(v string) MigrateConnectorOptionFunc {
	return func(o *MigrateConnectorOption) { o.DbIp = v }
}
func WithMigrateDbPort(v int) MigrateConnectorOptionFunc {
	return func(o *MigrateConnectorOption) { o.DbPort = v }
}
func WithMigrateDbInstance(v string) MigrateConnectorOptionFunc {
	return func(o *MigrateConnectorOption) { o.DbInstance = v }
}
func WithMigrateDbName(v string) MigrateConnectorOptionFunc {
	return func(o *MigrateConnectorOption) { o.DbName = v }
}
func WithMigrateDbUser(v string) MigrateConnectorOptionFunc {
	return func(o *MigrateConnectorOption) { o.DbUser = v }
}
func WithMigrateDbPwd(v string) MigrateConnectorOptionFunc {
	return func(o *MigrateConnectorOption) { o.DbPwd = v }
}
func WithMigrateTableArgs(v []int) MigrateConnectorOptionFunc {
	return func(o *MigrateConnectorOption) { o.TableArgs = v }
}
func WithMigrateCursorField(v string) MigrateConnectorOptionFunc {
	return func(o *MigrateConnectorOption) { o.CursorField = v }
}
func WithMigrateCursorType(v int8) MigrateConnectorOptionFunc {
	return func(o *MigrateConnectorOption) { o.CursorType = v }
}
func WithMigrateCursorPos(v string) MigrateConnectorOptionFunc {
	return func(o *MigrateConnectorOption) { o.CursorPos = v }
}
func WithMigrateIndex(v string) MigrateConnectorOptionFunc {
	return func(o *MigrateConnectorOption) { o.Index = v }
}
func WithMigrateTableFilter(v string) MigrateConnectorOptionFunc {
	return func(o *MigrateConnectorOption) { o.TableFilter = v }
}
func WithMigrateBatch(v int32) MigrateConnectorOptionFunc {
	return func(o *MigrateConnectorOption) { o.Batch = v }
}
func WithMigrateTimeout(v time.Duration) MigrateConnectorOptionFunc {
	return func(o *MigrateConnectorOption) { o.Timeout = v }
}
func WithMigrateTableMap(v map[string]string) MigrateConnectorOptionFunc {
	return func(o *MigrateConnectorOption) { o.TableMap = v }
}
func WithMigrateOuter(v IOuter) MigrateConnectorOptionFunc {
	return func(o *MigrateConnectorOption) { o.Outer = v }
}
func WithMigrateLogger(v ILogger) MigrateConnectorOptionFunc {
	return func(o *MigrateConnectorOption) { o.Log = v }
}
