package scanner

import (
	"fmt"
	"os"
	"path/filepath"

	"file_syn/pkg/models"
)

// FileScanner 文件扫描器
type FileScanner struct {
	rootPath string
	files    map[string]*models.FileInfo
}

// NewFileScanner 创建新的文件扫描器
func NewFileScanner(rootPath string) *FileScanner {
	return &FileScanner{
		rootPath: rootPath,
		files:    make(map[string]*models.FileInfo),
	}
}

// Scan 扫描目录
func (fs *FileScanner) Scan() error {
	return filepath.Walk(fs.rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// 如果无法访问某个文件，记录错误但继续扫描
			fmt.Fprintf(os.Stderr, "警告: 无法访问 %s: %v\n", path, err)
			return nil
		}

		// 计算相对路径
		relPath, err := filepath.Rel(fs.rootPath, path)
		if err != nil {
			relPath = path
		}

		// 统一路径分隔符为 /
		relPath = filepath.ToSlash(relPath)

		// 跳过根目录本身
		if relPath == "." {
			return nil
		}

		fileInfo := &models.FileInfo{
			Path:    relPath,
			Size:    info.Size(),
			ModTime: info.ModTime(),
			IsDir:   info.IsDir(),
			Mode:    info.Mode(),
			AbsPath: path,
		}

		fs.files[relPath] = fileInfo
		return nil
	})
}

// GetFiles 获取所有文件信息
func (fs *FileScanner) GetFiles() map[string]*models.FileInfo {
	return fs.files
}
