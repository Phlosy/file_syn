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

// formatFileInfo 格式化文件信息
func formatFileInfo(info *models.FileInfo) []string {
	if info == nil {
		return []string{"-"}
	}

	var lines []string
	if info.IsDir {
		lines = append(lines, fmt.Sprintf("📁 %s", info.Path))
		lines = append(lines, "   [目录]")
	} else {
		lines = append(lines, fmt.Sprintf("📄 %s", info.Path))
		lines = append(lines, fmt.Sprintf("   大小: %s", formatSize(info.Size)))
		lines = append(lines, fmt.Sprintf("   时间: %s", info.ModTime.Format("2006-01-02 15:04:05")))
		lines = append(lines, fmt.Sprintf("   权限: %s", info.Mode.Perm().String()))
	}
	return lines
}

// getStatusDisplay 获取状态显示文本（带符号）
func getStatusDisplay(status string) string {
	var symbol, text string
	switch status {
	case models.StatusAdded:
		symbol = "➕"
		text = "新增"
	case models.StatusDeleted:
		symbol = "➖"
		text = "删除"
	case models.StatusModified:
		symbol = "🔄"
		text = "修改"
	case models.StatusUnchanged:
		symbol = "✓"
		text = "未变更"
	default:
		symbol = "?"
		text = "未知"
	}
	return fmt.Sprintf("%s %s", symbol, text)
}

// wrapText 文本换行处理，将长文本按指定宽度换行
func wrapText(text string, width int) []string {
	if len(text) <= width {
		return []string{text}
	}

	var lines []string
	words := strings.Fields(text)
	if len(words) == 0 {
		// 如果没有空格，直接按字符截断
		for i := 0; i < len(text); i += width {
			end := i + width
			if end > len(text) {
				end = len(text)
			}
			lines = append(lines, text[i:end])
		}
		return lines
	}

	currentLine := ""
	for _, word := range words {
		if len(currentLine)+len(word)+1 <= width {
			if currentLine != "" {
				currentLine += " " + word
			} else {
				currentLine = word
			}
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			// 如果单个词就超过宽度，需要截断
			if len(word) > width {
				for i := 0; i < len(word); i += width {
					end := i + width
					if end > len(word) {
						end = len(word)
					}
					lines = append(lines, word[i:end])
				}
				currentLine = ""
			} else {
				currentLine = word
			}
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}
	return lines
}

// maxInt 返回两个整数中的较大值
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// PrintResults 打印对比结果（表格格式：左侧目录 | 右侧目录 | 状态）
func (r *Reporter) PrintResults(results []*models.DiffResult) {
	fmt.Println()
	fmt.Println("╔════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                                                                  文件同步监测结果                                                                              ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════════╝")
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
		// 列宽度定义
		const leftColWidth = 50
		const rightColWidth = 50
		const statusColWidth = 20

		// 打印表头
		fmt.Println("┌" + strings.Repeat("─", leftColWidth+2) + "┬" + strings.Repeat("─", rightColWidth+2) + "┬" + strings.Repeat("─", statusColWidth+2) + "┐")
		fmt.Printf("│ %-*s │ %-*s │ %-*s │\n", leftColWidth, "左侧目录", rightColWidth, "右侧目录", statusColWidth, "状态")
		fmt.Println("├" + strings.Repeat("─", leftColWidth+2) + "┼" + strings.Repeat("─", rightColWidth+2) + "┼" + strings.Repeat("─", statusColWidth+2) + "┤")

		// 打印表格内容
		for i, result := range displayResults {
			leftLines := formatFileInfo(result.LeftInfo)
			rightLines := formatFileInfo(result.RightInfo)
			status := getStatusDisplay(result.Status)

			// 如果有差异详情，添加到状态列
			statusLines := []string{status}
			if len(result.Differences) > 0 {
				for _, diff := range result.Differences {
					// 格式化差异信息
					diffLines := formatDiffDetails(diff, result)
					statusLines = append(statusLines, diffLines...)
				}
			}

			// 计算需要多少行
			maxLines := maxInt(len(leftLines), len(rightLines))
			maxLines = maxInt(maxLines, len(statusLines))

			// 打印每一行
			for lineIdx := 0; lineIdx < maxLines; lineIdx++ {
				var leftText, rightText, statusText string

				if lineIdx < len(leftLines) {
					leftText = leftLines[lineIdx]
				}
				if lineIdx < len(rightLines) {
					rightText = rightLines[lineIdx]
				}
				if lineIdx < len(statusLines) {
					statusText = statusLines[lineIdx]
				}

				// 处理长文本换行
				leftWrapped := wrapText(leftText, leftColWidth)
				rightWrapped := wrapText(rightText, rightColWidth)
				statusWrapped := wrapText(statusText, statusColWidth)

				// 计算需要多少行来显示（考虑换行）
				wrappedMaxLines := maxInt(len(leftWrapped), len(rightWrapped))
				wrappedMaxLines = maxInt(wrappedMaxLines, len(statusWrapped))

				// 打印换行后的内容
				for wrapIdx := 0; wrapIdx < wrappedMaxLines; wrapIdx++ {
					var leftWrap, rightWrap, statusWrap string
					if wrapIdx < len(leftWrapped) {
						leftWrap = leftWrapped[wrapIdx]
					}
					if wrapIdx < len(rightWrapped) {
						rightWrap = rightWrapped[wrapIdx]
					}
					if wrapIdx < len(statusWrapped) {
						statusWrap = statusWrapped[wrapIdx]
					}

					// 确保文本不超过列宽，并正确对齐
					leftDisplay := truncateString(leftWrap, leftColWidth)
					rightDisplay := truncateString(rightWrap, rightColWidth)
					statusDisplay := truncateString(statusWrap, statusColWidth)

					fmt.Printf("│ %-*s │ %-*s │ %-*s │\n",
						leftColWidth, leftDisplay,
						rightColWidth, rightDisplay,
						statusColWidth, statusDisplay)
				}
			}

			// 添加分隔线（最后一个不添加）
			if i < len(displayResults)-1 {
				fmt.Println("├" + strings.Repeat("─", leftColWidth+2) + "┼" + strings.Repeat("─", rightColWidth+2) + "┼" + strings.Repeat("─", statusColWidth+2) + "┤")
			}
		}

		fmt.Println("└" + strings.Repeat("─", leftColWidth+2) + "┴" + strings.Repeat("─", rightColWidth+2) + "┴" + strings.Repeat("─", statusColWidth+2) + "┘")
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
			lines = append(lines, fmt.Sprintf("大小: %s→%s", leftSize, rightSize))
		}
	} else if strings.Contains(diff, "修改时间不同") {
		if result.LeftInfo != nil && result.RightInfo != nil {
			leftTime := result.LeftInfo.ModTime.Format("2006-01-02 15:04:05")
			rightTime := result.RightInfo.ModTime.Format("2006-01-02 15:04:05")
			lines = append(lines, fmt.Sprintf("时间: %s→%s", leftTime, rightTime))
		}
	} else if strings.Contains(diff, "权限不同") {
		if result.LeftInfo != nil && result.RightInfo != nil {
			leftPerm := result.LeftInfo.Mode.Perm().String()
			rightPerm := result.RightInfo.Mode.Perm().String()
			lines = append(lines, fmt.Sprintf("权限: %s→%s", leftPerm, rightPerm))
		}
	} else if strings.Contains(diff, "仅存在于") {
		if strings.Contains(diff, "左侧") {
			lines = append(lines, "仅左侧存在")
		} else {
			lines = append(lines, "仅右侧存在")
		}
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
