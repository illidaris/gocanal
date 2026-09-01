package config

import "time"

var myConfig Config

type Config struct {
	EsName    string       `json:"es_name" yaml:"es_name"`
	EsUrls    []string     `json:"es_urls" yaml:"es_urls"`
	ESUser    string       `json:"es_user" yaml:"es_user"`
	ESPwd     string       `json:"es_pwd" yaml:"es_pwd"`
	CanalIp   string       `json:"canal_ip" yaml:"canal_ip"`
	CanalPort int          `json:"canal_port" yaml:"canal_port"`
	CanalUser string       `json:"canal_user" yaml:"canal_user"`
	CanalPwd  string       `json:"canal_pwd" yaml:"canal_pwd"`
	Syncs     []SyncConfig `json:"syncs" yaml:"syncs"`
}

type SyncConfig struct {
	Instance string        `json:"instance" yaml:"instance"`
	Index    string        `json:"index" yaml:"index"`
	Filter   string        `json:"filter" yaml:"filter"`
	Batch    int32         `json:"batch" yaml:"batch"`
	Timeout  time.Duration `json:"timeout" yaml:"timeout"`
}
