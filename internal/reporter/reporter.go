package reporter

import (
	"fmt"

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

// PrintResults 打印对比结果
func (r *Reporter) PrintResults(results []*models.DiffResult) {
	fmt.Println()
	fmt.Println("=== 文件同步监测结果 ===")
	fmt.Println()

	addedCount := 0
	deletedCount := 0
	modifiedCount := 0
	unchangedCount := 0

	for _, result := range results {
		switch result.Status {
		case models.StatusAdded:
			addedCount++
			fmt.Printf("[新增] %s\n", result.Path)
			for _, diff := range result.Differences {
				fmt.Printf("  - %s\n", diff)
			}
			fmt.Println()
		case models.StatusDeleted:
			deletedCount++
			fmt.Printf("[删除] %s\n", result.Path)
			for _, diff := range result.Differences {
				fmt.Printf("  - %s\n", diff)
			}
			fmt.Println()
		case models.StatusModified:
			modifiedCount++
			fmt.Printf("[修改] %s\n", result.Path)
			for _, diff := range result.Differences {
				fmt.Printf("  - %s\n", diff)
			}
			fmt.Println()
		case models.StatusUnchanged:
			unchangedCount++
			if r.showUnchanged {
				fmt.Printf("[未变更] %s\n", result.Path)
			}
		}
	}

	// 打印统计信息
	fmt.Println("=== 统计信息 ===")
	fmt.Printf("新增文件: %d\n", addedCount)
	fmt.Printf("删除文件: %d\n", deletedCount)
	fmt.Printf("修改文件: %d\n", modifiedCount)
	if r.showUnchanged {
		fmt.Printf("未变更文件: %d\n", unchangedCount)
	}
	fmt.Printf("总计: %d\n", len(results))
}
