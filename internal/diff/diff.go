package diff

import (
	"fmt"
	"sort"
	"time"

	"file_syn/internal/scanner"
	"file_syn/pkg/models"
)

// Comparer 目录对比器
type Comparer struct{}

// NewComparer 创建新的对比器
func NewComparer() *Comparer {
	return &Comparer{}
}

// Compare 对比两个目录
func (c *Comparer) Compare(leftDir, rightDir string) ([]*models.DiffResult, error) {
	// 扫描左侧目录
	leftScanner := scanner.NewFileScanner(leftDir)
	if err := leftScanner.Scan(); err != nil {
		return nil, fmt.Errorf("扫描左侧目录失败: %v", err)
	}

	// 扫描右侧目录
	rightScanner := scanner.NewFileScanner(rightDir)
	if err := rightScanner.Scan(); err != nil {
		return nil, fmt.Errorf("扫描右侧目录失败: %v", err)
	}

	leftFiles := leftScanner.GetFiles()
	rightFiles := rightScanner.GetFiles()

	// 收集所有文件路径
	allPaths := make(map[string]bool)
	for path := range leftFiles {
		allPaths[path] = true
	}
	for path := range rightFiles {
		allPaths[path] = true
	}

	// 转换为排序的切片
	var sortedPaths []string
	for path := range allPaths {
		sortedPaths = append(sortedPaths, path)
	}
	sort.Strings(sortedPaths)

	var results []*models.DiffResult

	// 对比每个文件
	for _, path := range sortedPaths {
		leftFile := leftFiles[path]
		rightFile := rightFiles[path]

		result := &models.DiffResult{
			Path:        path,
			LeftInfo:    leftFile,
			RightInfo:   rightFile,
			Differences: []string{},
		}

		if leftFile == nil {
			// 文件只在右侧存在
			result.Status = models.StatusAdded
			result.Differences = []string{"文件仅存在于右侧目录"}
		} else if rightFile == nil {
			// 文件只在左侧存在
			result.Status = models.StatusDeleted
			result.Differences = []string{"文件仅存在于左侧目录"}
		} else {
			// 文件在两侧都存在，检查差异
			diffs := compareFileInfo(leftFile, rightFile)
			if len(diffs) > 0 {
				result.Status = models.StatusModified
				result.Differences = diffs
			} else {
				result.Status = models.StatusUnchanged
			}
		}

		results = append(results, result)
	}

	return results, nil
}

// compareFileInfo 对比两个文件信息
func compareFileInfo(left, right *models.FileInfo) []string {
	var differences []string

	// 检查是否为目录
	if left.IsDir != right.IsDir {
		if left.IsDir {
			differences = append(differences, "左侧是目录，右侧不是")
		} else {
			differences = append(differences, "右侧是目录，左侧不是")
		}
		return differences
	}

	// 如果是目录，只检查类型差异（已在上面检查）
	if left.IsDir {
		return differences
	}

	// 对比文件大小
	if left.Size != right.Size {
		differences = append(differences, fmt.Sprintf("大小不同: 左侧=%d 字节, 右侧=%d 字节", left.Size, right.Size))
	}

	// 对比修改时间（允许1秒的误差，因为不同文件系统的时间精度可能不同）
	timeDiff := left.ModTime.Sub(right.ModTime)
	if timeDiff < 0 {
		timeDiff = -timeDiff
	}
	if timeDiff > time.Second {
		differences = append(differences, fmt.Sprintf("修改时间不同: 左侧=%s, 右侧=%s",
			left.ModTime.Format("2006-01-02 15:04:05"),
			right.ModTime.Format("2006-01-02 15:04:05")))
	}

	// 对比文件权限（只对比基本权限位，忽略特殊位）
	leftPerm := left.Mode.Perm()
	rightPerm := right.Mode.Perm()
	if leftPerm != rightPerm {
		differences = append(differences, fmt.Sprintf("权限不同: 左侧=%s, 右侧=%s",
			leftPerm.String(), rightPerm.String()))
	}

	return differences
}
