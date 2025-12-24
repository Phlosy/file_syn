package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// 创建临时配置文件
	tmpDir, err := os.MkdirTemp("", "file_syn_config_test_*")
	if err != nil {
		t.Fatalf("无法创建临时目录: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建测试目录
	leftDir := filepath.Join(tmpDir, "left")
	rightDir := filepath.Join(tmpDir, "right")
	if err := os.MkdirAll(leftDir, 0755); err != nil {
		t.Fatalf("无法创建左侧目录: %v", err)
	}
	if err := os.MkdirAll(rightDir, 0755); err != nil {
		t.Fatalf("无法创建右侧目录: %v", err)
	}

	// 创建配置文件
	configPath := filepath.Join(tmpDir, "config.json")
	configContent := `{
  "left_dir": "` + leftDir + `",
  "right_dir": "` + rightDir + `",
  "show_unchanged": true
}`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("无法创建配置文件: %v", err)
	}

	// 测试加载配置
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	// 验证配置
	if cfg.LeftDir == "" {
		t.Error("LeftDir 为空")
	}

	if cfg.RightDir == "" {
		t.Error("RightDir 为空")
	}

	if !cfg.ShowUnchanged {
		t.Error("ShowUnchanged 应该为 true")
	}

	// 验证路径已转换为绝对路径
	if !filepath.IsAbs(cfg.LeftDir) {
		t.Error("LeftDir 应该是绝对路径")
	}

	if !filepath.IsAbs(cfg.RightDir) {
		t.Error("RightDir 应该是绝对路径")
	}
}

func TestConfigValidate(t *testing.T) {
	// 测试空配置
	cfg := &Config{}
	if err := cfg.Validate(); err == nil {
		t.Error("空配置应该验证失败")
	}

	// 测试不存在的目录
	cfg = &Config{
		LeftDir:  "/nonexistent/left",
		RightDir: "/nonexistent/right",
	}
	if err := cfg.Validate(); err == nil {
		t.Error("不存在的目录应该验证失败")
	}
}
