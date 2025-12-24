package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// displayWidth 计算字符串的显示宽度
func displayWidth(s string) int {
	width := 0
	for _, r := range s {
		if r >= 0x1100 && (r <= 0x115F || r >= 0x2E80 && r <= 0xFFEF) {
			width += 2
		} else if r != utf8.RuneError {
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

func main() {
	const leftColWidth = 50
	const rightColWidth = 50
	const statusColWidth = 20

	leftSeparatorLen := leftColWidth + 2
	rightSeparatorLen := rightColWidth + 2
	statusSeparatorLen := statusColWidth + 2

	// 测试行
	leftText := padString("📄 changed.txt", leftColWidth, true)
	rightText := padString("📄 changed.txt", rightColWidth, true)
	statusText := padString("🔄 修改", statusColWidth, true)

	// 打印分隔线
	headerLine := "┌" + strings.Repeat("─", leftSeparatorLen) + "┬" + strings.Repeat("─", rightSeparatorLen) + "┬" + strings.Repeat("─", statusSeparatorLen) + "┐"
	fmt.Println(headerLine)

	// 打印内容行
	contentLine := fmt.Sprintf("│ %s │ %s │ %s │", leftText, rightText, statusText)
	fmt.Println(contentLine)

	// 打印分隔线
	separatorLine := "├" + strings.Repeat("─", leftSeparatorLen) + "┼" + strings.Repeat("─", rightSeparatorLen) + "┼" + strings.Repeat("─", statusSeparatorLen) + "┤"
	fmt.Println(separatorLine)

	// 验证长度
	fmt.Printf("Header line length: %d\n", len(headerLine))
	fmt.Printf("Content line length: %d\n", len(contentLine))
	fmt.Printf("Separator line length: %d\n", len(separatorLine))
}

