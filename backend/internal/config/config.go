// Package config 负责从 YAML 文件加载应用配置。
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 应用配置根结构，字段与 config.yaml 一一对应。
type Config struct {
	App      AppConfig      `yaml:"app"`
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

// AppConfig 应用级别配置。
type AppConfig struct {
	Env string `yaml:"env"`
}

// ServerConfig HTTP 服务配置。
type ServerConfig struct {
	Addr string `yaml:"addr"`
}

// DatabaseConfig 数据库连接配置。
type DatabaseConfig struct {
	DSN string `yaml:"dsn"`
}

// Load 从指定路径读取 YAML 配置文件并解析。
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}
	return &cfg, nil
}
