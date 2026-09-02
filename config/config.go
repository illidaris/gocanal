package config

import "time"

type Config struct {
	EsName    string          `json:"es_name" yaml:"es_name" mapstructure:"es_name"`
	EsUrls    []string        `json:"es_urls" yaml:"es_urls" mapstructure:"es_urls"`
	ESUser    string          `json:"es_user" yaml:"es_user" mapstructure:"es_user"`
	ESPwd     string          `json:"es_pwd" yaml:"es_pwd" mapstructure:"es_pwd"`
	CanalIp   string          `json:"canal_ip" yaml:"canal_ip" mapstructure:"canal_ip"`
	CanalPort int             `json:"canal_port" yaml:"canal_port" mapstructure:"canal_port"`
	CanalUser string          `json:"canal_user" yaml:"canal_user" mapstructure:"canal_user"`
	CanalPwd  string          `json:"canal_pwd" yaml:"canal_pwd" mapstructure:"canal_pwd"`
	Syncs     []SyncConfig    `json:"syncs" yaml:"syncs" mapstructure:"syncs"`
	Migrates  []MigrateConfig `json:"migrates" yaml:"migrates" mapstructure:"migrates"`
}

type SyncConfig struct {
	Instance string        `json:"instance" yaml:"instance" mapstructure:"instance"`
	Index    string        `json:"index" yaml:"index" mapstructure:"index"`
	Filter   string        `json:"filter" yaml:"filter" mapstructure:"filter"`
	Batch    int32         `json:"batch" yaml:"batch" mapstructure:"batch"`
	Timeout  time.Duration `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
}

type MigrateConfig struct {
	DbIp        string        `json:"db_ip" yaml:"db_ip" mapstructure:"db_ip"`
	DbPort      int           `json:"db_port" yaml:"db_port" mapstructure:"db_port"`
	DbInstance  string        `json:"db_instance" yaml:"db_instance" mapstructure:"db_instance"`
	DbName      string        `json:"db_name" yaml:"db_name" mapstructure:"db_name"`
	DbUser      string        `json:"db_user" yaml:"db_user" mapstructure:"db_user"`
	DbPwd       string        `json:"db_pwd" yaml:"db_pwd" mapstructure:"db_pwd"`
	TableArgs   []int         `json:"table_args" yaml:"table_args" mapstructure:"table_args"`
	CursorField string        `json:"cursor_field" yaml:"cursor_field" mapstructure:"cursor_field"`
	CursorType  int8          `json:"cursor_type" yaml:"cursor_type" mapstructure:"cursor_type"` // 0- int 1- string
	CursorPos   string        `json:"cursor_pos" yaml:"cursor_pos" mapstructure:"cursor_pos"`
	Instance    string        `json:"instance" yaml:"instance" mapstructure:"instance"`
	Index       string        `json:"index" yaml:"index" mapstructure:"index"`
	Filter      string        `json:"filter" yaml:"filter" mapstructure:"filter"`
	Batch       int32         `json:"batch" yaml:"batch" mapstructure:"batch"`
	Timeout     time.Duration `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
}
