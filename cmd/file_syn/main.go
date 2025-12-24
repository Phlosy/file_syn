package main

import (
	"fmt"
	"os"

	"file_syn/internal/config"
	"file_syn/internal/diff"
	"file_syn/internal/reporter"
)

func main() {
	// 获取配置文件路径（如果通过命令行参数指定）
	configPath := ""
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	// 加载配置
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		fmt.Fprintf(os.Stderr, "\n用法: %s [配置文件路径]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "示例: %s\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "      %s config/config.json\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "      %s /path/to/custom-config.json\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n如果未指定配置文件路径，程序将按以下顺序查找:\n")
		fmt.Fprintf(os.Stderr, "  1. config/config.json\n")
		fmt.Fprintf(os.Stderr, "  2. ./config/config.json\n")
		fmt.Fprintf(os.Stderr, "  3. config.json\n")
		os.Exit(1)
	}

	// 显示使用的配置文件路径
	fmt.Printf("配置文件: %s\n", cfg.ConfigPath)
	fmt.Printf("左侧目录: %s\n", cfg.LeftDir)
	fmt.Printf("右侧目录: %s\n", cfg.RightDir)
	fmt.Println("正在扫描和对比...")

	// 执行对比
	comparer := diff.NewComparer()
	results, err := comparer.Compare(cfg.LeftDir, cfg.RightDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}

	// 打印结果
	reporter := reporter.NewReporter(cfg.ShowUnchanged)
	reporter.PrintResults(results)
}
