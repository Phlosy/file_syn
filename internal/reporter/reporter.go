package reporter

import (
	"fmt"
	"strings"

	"file_syn/pkg/models"
)

// Reporter 结果报告器
type Reporter struct {
	showUnchanged bool
}

// NewReporter 创建新的报告器
func NewReporter(showUnchanged bool) *Reporter {
	return &Reporter{
		showUnchanged: showUnchanged,
	}
}

// formatSize 格式化文件大小
func formatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// getStatusSymbol 获取状态符号
func getStatusSymbol(status string) string {
	switch status {
	case models.StatusAdded:
		return "➕"
	case models.StatusDeleted:
		return "➖"
	case models.StatusModified:
		return "🔄"
	case models.StatusUnchanged:
		return "✓"
	default:
		return "?"
	}
}

// getStatusText 获取状态文本
func getStatusText(status string) string {
	switch status {
	case models.StatusAdded:
		return "新增"
	case models.StatusDeleted:
		return "删除"
	case models.StatusModified:
		return "修改"
	case models.StatusUnchanged:
		return "未变更"
	default:
		return "未知"
	}
}

// getStatusDisplay 获取状态显示文本（带符号）
func getStatusDisplay(status string) string {
	symbol := getStatusSymbol(status)
	text := getStatusText(status)
	return fmt.Sprintf("%s %s", symbol, text)
}

// PrintResults 打印对比结果（表格格式）
func (r *Reporter) PrintResults(results []*models.DiffResult) {
	fmt.Println()
	fmt.Println("╔════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                        文件同步监测结果                                    ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// 统计信息
	addedCount := 0
	deletedCount := 0
	modifiedCount := 0
	unchangedCount := 0

	// 过滤需要显示的结果
	var displayResults []*models.DiffResult
	for _, result := range results {
		if result.Status == models.StatusUnchanged && !r.showUnchanged {
			unchangedCount++
			continue
		}
		displayResults = append(displayResults, result)
		switch result.Status {
		case models.StatusAdded:
			addedCount++
		case models.StatusDeleted:
			deletedCount++
		case models.StatusModified:
			modifiedCount++
		case models.StatusUnchanged:
			unchangedCount++
		}
	}

	if len(displayResults) == 0 {
		fmt.Println("  所有文件一致，无差异")
		fmt.Println()
	} else {
		// 打印表头
		fmt.Println("┌──────────┬──────────────────────────────────────────────────────────────┐")
		fmt.Printf("│ %-8s │ %-60s │\n", "状态", "文件路径")
		fmt.Println("├──────────┼──────────────────────────────────────────────────────────────┤")

		// 打印表格内容
		for _, result := range displayResults {
			status := getStatusDisplay(result.Status)

			// 处理长路径（超过60字符时截断）
			path := result.Path
			if len(path) > 60 {
				path = path[:57] + "..."
			}

			fmt.Printf("│ %-8s │ %-60s │\n", status, path)

			// 如果有差异详情，显示详细信息
			if len(result.Differences) > 0 {
				for _, diff := range result.Differences {
					// 格式化差异信息
					diffLines := formatDiffDetails(diff, result)
					for _, line := range diffLines {
						fmt.Printf("│          │   %-58s │\n", truncateString(line, 58))
					}
				}
			}

			// 添加分隔线（最后一个不添加）
			if result != displayResults[len(displayResults)-1] {
				fmt.Println("├──────────┼──────────────────────────────────────────────────────────────┤")
			}
		}

		fmt.Println("└──────────┴──────────────────────────────────────────────────────────────┘")
		fmt.Println()
	}

	// 打印统计信息表格
	fmt.Println("╔════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                              统计信息                                       ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	fmt.Println("┌──────────────────┬────────┐")
	fmt.Printf("│ %-16s │ %6d │\n", "新增文件", addedCount)
	fmt.Println("├──────────────────┼────────┤")
	fmt.Printf("│ %-16s │ %6d │\n", "删除文件", deletedCount)
	fmt.Println("├──────────────────┼────────┤")
	fmt.Printf("│ %-16s │ %6d │\n", "修改文件", modifiedCount)
	if r.showUnchanged {
		fmt.Println("├──────────────────┼────────┤")
		fmt.Printf("│ %-16s │ %6d │\n", "未变更文件", unchangedCount)
	}
	fmt.Println("├──────────────────┼────────┤")
	fmt.Printf("│ %-16s │ %6d │\n", "总计", len(results))
	fmt.Println("└──────────────────┴────────┘")
}

// formatDiffDetails 格式化差异详情
func formatDiffDetails(diff string, result *models.DiffResult) []string {
	var lines []string

	// 解析差异信息
	if strings.Contains(diff, "大小不同") {
		if result.LeftInfo != nil && result.RightInfo != nil {
			leftSize := formatSize(result.LeftInfo.Size)
			rightSize := formatSize(result.RightInfo.Size)
			lines = append(lines, fmt.Sprintf("大小: %s → %s", leftSize, rightSize))
		}
	} else if strings.Contains(diff, "修改时间不同") {
		if result.LeftInfo != nil && result.RightInfo != nil {
			leftTime := result.LeftInfo.ModTime.Format("2006-01-02 15:04:05")
			rightTime := result.RightInfo.ModTime.Format("2006-01-02 15:04:05")
			lines = append(lines, fmt.Sprintf("修改时间: %s → %s", leftTime, rightTime))
		}
	} else if strings.Contains(diff, "权限不同") {
		if result.LeftInfo != nil && result.RightInfo != nil {
			leftPerm := result.LeftInfo.Mode.Perm().String()
			rightPerm := result.RightInfo.Mode.Perm().String()
			lines = append(lines, fmt.Sprintf("权限: %s → %s", leftPerm, rightPerm))
		}
	} else if strings.Contains(diff, "仅存在于") {
		lines = append(lines, diff)
	} else {
		lines = append(lines, diff)
	}

	return lines
}

// truncateString 截断字符串
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
