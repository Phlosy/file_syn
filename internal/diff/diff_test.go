package diff

import (
	"os"
	"path/filepath"
	"testing"
)

func TestComparer(t *testing.T) {
	// 创建两个临时目录用于测试
	leftDir, err := os.MkdirTemp("", "file_syn_left_*")
	if err != nil {
		t.Fatalf("无法创建左侧临时目录: %v", err)
	}
	defer os.RemoveAll(leftDir)

	rightDir, err := os.MkdirTemp("", "file_syn_right_*")
	if err != nil {
		t.Fatalf("无法创建右侧临时目录: %v", err)
	}
	defer os.RemoveAll(rightDir)

	// 在左侧目录创建文件
	leftFile := filepath.Join(leftDir, "common.txt")
	if err := os.WriteFile(leftFile, []byte("common content"), 0644); err != nil {
		t.Fatalf("无法创建左侧文件: %v", err)
	}

	// 在右侧目录创建同名文件（内容不同）
	rightFile := filepath.Join(rightDir, "common.txt")
	if err := os.WriteFile(rightFile, []byte("different content"), 0644); err != nil {
		t.Fatalf("无法创建右侧文件: %v", err)
	}

	// 只在左侧创建文件
	leftOnlyFile := filepath.Join(leftDir, "left_only.txt")
	if err := os.WriteFile(leftOnlyFile, []byte("left only"), 0644); err != nil {
		t.Fatalf("无法创建左侧独有文件: %v", err)
	}

	// 只在右侧创建文件
	rightOnlyFile := filepath.Join(rightDir, "right_only.txt")
	if err := os.WriteFile(rightOnlyFile, []byte("right only"), 0644); err != nil {
		t.Fatalf("无法创建右侧独有文件: %v", err)
	}

	// 执行对比
	comparer := NewComparer()
	results, err := comparer.Compare(leftDir, rightDir)
	if err != nil {
		t.Fatalf("对比失败: %v", err)
	}

	// 验证结果
	if len(results) != 3 {
		t.Errorf("期望 3 个差异结果，实际得到 %d 个", len(results))
	}

	// 验证每个结果的状态
	foundAdded := false
	foundDeleted := false
	foundModified := false

	for _, result := range results {
		switch result.Path {
		case "right_only.txt":
			if result.Status != "added" {
				t.Errorf("right_only.txt 应该是 added 状态，实际是 %s", result.Status)
			}
			foundAdded = true
		case "left_only.txt":
			if result.Status != "deleted" {
				t.Errorf("left_only.txt 应该是 deleted 状态，实际是 %s", result.Status)
			}
			foundDeleted = true
		case "common.txt":
			if result.Status != "modified" {
				t.Errorf("common.txt 应该是 modified 状态，实际是 %s", result.Status)
			}
			foundModified = true
		}
	}

	if !foundAdded {
		t.Error("未找到 added 状态的文件")
	}
	if !foundDeleted {
		t.Error("未找到 deleted 状态的文件")
	}
	if !foundModified {
		t.Error("未找到 modified 状态的文件")
	}
}
