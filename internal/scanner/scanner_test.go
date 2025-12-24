package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileScanner(t *testing.T) {
	// 创建临时目录用于测试
	tmpDir, err := os.MkdirTemp("", "file_syn_test_*")
	if err != nil {
		t.Fatalf("无法创建临时目录: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建测试文件
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("无法创建测试文件: %v", err)
	}

	// 创建测试子目录
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("无法创建子目录: %v", err)
	}

	// 创建子目录中的文件
	subFile := filepath.Join(subDir, "subfile.txt")
	if err := os.WriteFile(subFile, []byte("sub content"), 0644); err != nil {
		t.Fatalf("无法创建子文件: %v", err)
	}

	// 测试扫描
	scanner := NewFileScanner(tmpDir)
	if err := scanner.Scan(); err != nil {
		t.Fatalf("扫描失败: %v", err)
	}

	files := scanner.GetFiles()

	// 验证文件数量（应该包含 test.txt, subdir/, subdir/subfile.txt）
	expectedCount := 3
	if len(files) != expectedCount {
		t.Errorf("期望 %d 个文件，实际得到 %d 个", expectedCount, len(files))
	}

	// 验证特定文件存在
	if _, exists := files["test.txt"]; !exists {
		t.Error("未找到 test.txt")
	}

	if _, exists := files["subdir"]; !exists {
		t.Error("未找到 subdir")
	}

	if _, exists := files["subdir/subfile.txt"]; !exists {
		t.Error("未找到 subdir/subfile.txt")
	}
}
