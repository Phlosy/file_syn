package reporter

import (
	"fmt"
	"strings"
	"unicode/utf8"

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

// displayWidth 计算字符串的显示宽度（中文字符占2个宽度，emoji通常占2个宽度）
func displayWidth(s string) int {
	width := 0
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		r := runes[i]

		// 检查是否是 emoji 或特殊符号（通常占2个宽度）
		if r >= 0x1F300 && r <= 0x1F9FF { // Emoji range
			width += 2
			continue
		}
		if r >= 0x2600 && r <= 0x26FF { // Miscellaneous Symbols
			width += 2
			continue
		}
		if r >= 0x2700 && r <= 0x27BF { // Dingbats
			width += 2
			continue
		}
		if r >= 0xFE00 && r <= 0xFE0F { // Variation Selectors (emoji modifiers)
			width += 0 // 这些是修饰符，不占宽度
			continue
		}
		if r >= 0x200D { // Zero Width Joiner (emoji sequences)
			// 检查是否是 emoji 序列的一部分
			if i+1 < len(runes) && (runes[i+1] >= 0x1F300 && runes[i+1] <= 0x1F9FF) {
				width += 0
				continue
			}
		}

		// 判断是否为全角字符（中文、日文、韩文等）
		if r >= 0x1100 && (r <= 0x115F || // Hangul Jamo
			r >= 0x2E80 && r <= 0x2EFF || // CJK Radicals Supplement
			r >= 0x2F00 && r <= 0x2FDF || // Kangxi Radicals
			r >= 0x3000 && r <= 0x303F || // CJK Symbols and Punctuation
			r >= 0x3040 && r <= 0x309F || // Hiragana
			r >= 0x30A0 && r <= 0x30FF || // Katakana
			r >= 0x3100 && r <= 0x312F || // Bopomofo
			r >= 0x3130 && r <= 0x318F || // Hangul Compatibility Jamo
			r >= 0x3200 && r <= 0x32FF || // Enclosed CJK Letters and Months
			r >= 0x3300 && r <= 0x33FF || // CJK Compatibility
			r >= 0x3400 && r <= 0x4DBF || // CJK Unified Ideographs Extension A
			r >= 0x4E00 && r <= 0x9FFF || // CJK Unified Ideographs
			r >= 0xA000 && r <= 0xA48F || // Yi Syllables
			r >= 0xA490 && r <= 0xA4CF || // Yi Radicals
			r >= 0xAC00 && r <= 0xD7AF || // Hangul Syllables
			r >= 0xF900 && r <= 0xFAFF || // CJK Compatibility Ideographs
			r >= 0xFE30 && r <= 0xFE4F || // CJK Compatibility Forms
			r >= 0xFF00 && r <= 0xFFEF) { // Halfwidth and Fullwidth Forms
			width += 2
		} else if r == utf8.RuneError {
			width += 1
		} else {
			width += 1
		}
	}
	return width
}

// padString 填充字符串到指定显示宽度
func padString(s string, width int, alignLeft bool) string {
	currentWidth := displayWidth(s)
	if currentWidth >= width {
		return s
	}

	padding := width - currentWidth
	if alignLeft {
		return s + strings.Repeat(" ", padding)
	}
	return strings.Repeat(" ", padding) + s
}

// truncateStringByWidth 按显示宽度截断字符串
func truncateStringByWidth(s string, maxWidth int) string {
	if displayWidth(s) <= maxWidth {
		return s
	}

	width := 0
	var result strings.Builder
	for _, r := range s {
		charWidth := 2
		if r < 0x1100 || (r > 0x115F && r < 0x2E80) || (r > 0xFFEF) {
			if r != utf8.RuneError {
				charWidth = 1
			}
		}
		if width+charWidth > maxWidth-3 {
			result.WriteString("...")
			break
		}
		result.WriteRune(r)
		width += charWidth
	}
	return result.String()
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

// wrapTextByWidth 按显示宽度换行文本
func wrapTextByWidth(text string, width int) []string {
	if displayWidth(text) <= width {
		return []string{text}
	}

	var lines []string
	currentLine := ""
	currentWidth := 0

	for _, r := range text {
		charWidth := 2
		if r < 0x1100 || (r > 0x115F && r < 0x2E80) || (r > 0xFFEF) {
			if r != utf8.RuneError {
				charWidth = 1
			}
		}

		if currentWidth+charWidth > width {
			if currentLine != "" {
				lines = append(lines, currentLine)
				currentLine = string(r)
				currentWidth = charWidth
			} else {
				// 单个字符就超过宽度，强制换行
				lines = append(lines, string(r))
				currentWidth = 0
			}
		} else {
			currentLine += string(r)
			currentWidth += charWidth
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
		// 列宽度定义（显示宽度）
		const leftColWidth = 50
		const rightColWidth = 50
		const statusColWidth = 20

		// 计算分隔线的实际字符数
		// 每列格式：│ + 空格(1) + 内容(width) + 空格(1) = width + 2
		leftSeparatorLen := leftColWidth + 2
		rightSeparatorLen := rightColWidth + 2
		statusSeparatorLen := statusColWidth + 2

		// 创建分隔线的辅助函数
		createSeparator := func(left, middle, right string) string {
			return left + strings.Repeat("─", leftSeparatorLen) + middle + strings.Repeat("─", rightSeparatorLen) + middle + strings.Repeat("─", statusSeparatorLen) + right
		}

		// 打印表头
		headerLine := createSeparator("┌", "┬", "┐")
		fmt.Println(headerLine)

		leftHeader := padString("左侧目录", leftColWidth, true)
		rightHeader := padString("右侧目录", rightColWidth, true)
		statusHeader := padString("状态", statusColWidth, true)
		fmt.Printf("│ %s │ %s │ %s │\n", leftHeader, rightHeader, statusHeader)

		separatorLine := createSeparator("├", "┼", "┤")
		fmt.Println(separatorLine)

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

				// 按显示宽度换行
				leftWrapped := wrapTextByWidth(leftText, leftColWidth)
				rightWrapped := wrapTextByWidth(rightText, rightColWidth)
				statusWrapped := wrapTextByWidth(statusText, statusColWidth)

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

					// 按显示宽度截断并填充
					leftDisplay := truncateStringByWidth(leftWrap, leftColWidth)
					rightDisplay := truncateStringByWidth(rightWrap, rightColWidth)
					statusDisplay := truncateStringByWidth(statusWrap, statusColWidth)

					// 填充到指定宽度
					leftPadded := padString(leftDisplay, leftColWidth, true)
					rightPadded := padString(rightDisplay, rightColWidth, true)
					statusPadded := padString(statusDisplay, statusColWidth, true)

					fmt.Printf("│ %s │ %s │ %s │\n", leftPadded, rightPadded, statusPadded)
				}
			}

			// 添加分隔线（最后一个不添加）
			if i < len(displayResults)-1 {
				fmt.Println(separatorLine)
			}
		}

		footerLine := createSeparator("└", "┴", "┘")
		fmt.Println(footerLine)
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
