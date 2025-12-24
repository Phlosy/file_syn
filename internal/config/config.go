package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config 配置结构
type Config struct {
	LeftDir       string `json:"left_dir"`
	RightDir      string `json:"right_dir"`
	ShowUnchanged bool   `json:"show_unchanged"`
	ConfigPath    string `json:"-"` // 实际使用的配置文件路径（不序列化）
}

// LoadConfig 从文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	// 如果配置文件路径为空，使用默认路径
	if configPath == "" {
		// 尝试从当前目录查找 config.json
		defaultPaths := []string{
			"config/config.json",
			"./config/config.json",
			"config.json",
		}

		for _, path := range defaultPaths {
			if _, err := os.Stat(path); err == nil {
				configPath = path
				break
			}
		}

		if configPath == "" {
			return nil, fmt.Errorf("未找到配置文件，请指定配置文件路径或确保 config/config.json 存在")
		}
	}

	// 转换为绝对路径
	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		absConfigPath = configPath
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("无法读取配置文件 %s: %v", configPath, err)
	}

	// 解析 JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("无法解析配置文件 %s: %v", configPath, err)
	}

	// 保存配置文件路径
	config.ConfigPath = absConfigPath

	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %v", err)
	}

	// 转换为绝对路径
	if err := config.NormalizePaths(); err != nil {
		return nil, fmt.Errorf("路径规范化失败: %v", err)
	}

	return &config, nil
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.LeftDir == "" {
		return fmt.Errorf("left_dir 不能为空")
	}

	if c.RightDir == "" {
		return fmt.Errorf("right_dir 不能为空")
	}

	// 检查目录是否存在
	if _, err := os.Stat(c.LeftDir); os.IsNotExist(err) {
		return fmt.Errorf("左侧目录不存在: %s", c.LeftDir)
	}

	if _, err := os.Stat(c.RightDir); os.IsNotExist(err) {
		return fmt.Errorf("右侧目录不存在: %s", c.RightDir)
	}

	return nil
}

// NormalizePaths 规范化路径为绝对路径
func (c *Config) NormalizePaths() error {
	leftAbs, err := filepath.Abs(c.LeftDir)
	if err != nil {
		return fmt.Errorf("无法获取左侧目录的绝对路径: %v", err)
	}

	rightAbs, err := filepath.Abs(c.RightDir)
	if err != nil {
		return fmt.Errorf("无法获取右侧目录的绝对路径: %v", err)
	}

	c.LeftDir = leftAbs
	c.RightDir = rightAbs

	return nil
}
