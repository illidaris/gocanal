package config

import "time"

type Config struct {
	EsName    string       `json:"es_name" yaml:"es_name" mapstructure:"es_name"`
	EsUrls    []string     `json:"es_urls" yaml:"es_urls" mapstructure:"es_urls"`
	ESUser    string       `json:"es_user" yaml:"es_user" mapstructure:"es_user"`
	ESPwd     string       `json:"es_pwd" yaml:"es_pwd" mapstructure:"es_pwd"`
	CanalIp   string       `json:"canal_ip" yaml:"canal_ip" mapstructure:"canal_ip"`
	CanalPort int          `json:"canal_port" yaml:"canal_port" mapstructure:"canal_port"`
	CanalUser string       `json:"canal_user" yaml:"canal_user" mapstructure:"canal_user"`
	CanalPwd  string       `json:"canal_pwd" yaml:"canal_pwd" mapstructure:"canal_pwd"`
	Syncs     []SyncConfig `json:"syncs" yaml:"syncs" mapstructure:"syncs"`
}

type SyncConfig struct {
	Instance string        `json:"instance" yaml:"instance" mapstructure:"instance"`
	Index    string        `json:"index" yaml:"index" mapstructure:"index"`
	Filter   string        `json:"filter" yaml:"filter" mapstructure:"filter"`
	Batch    int32         `json:"batch" yaml:"batch" mapstructure:"batch"`
	Timeout  time.Duration `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
}
