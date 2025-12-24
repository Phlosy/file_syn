package models

import (
	"os"
	"time"
)

// FileInfo 存储文件的元数据信息
type FileInfo struct {
	Path    string      // 相对路径
	Size    int64       // 文件大小（字节）
	ModTime time.Time   // 修改时间
	IsDir   bool        // 是否为目录
	Mode    os.FileMode // 文件权限
	AbsPath string      // 绝对路径（用于区分来源）
}

// DiffResult 存储差异结果
type DiffResult struct {
	Path        string    // 文件相对路径
	Status      string    // 差异状态：added, deleted, modified, unchanged
	LeftInfo    *FileInfo // 左侧目录的文件信息（如果存在）
	RightInfo   *FileInfo // 右侧目录的文件信息（如果存在）
	Differences []string  // 差异的属性列表
}

// Status constants
const (
	StatusAdded     = "added"
	StatusDeleted   = "deleted"
	StatusModified  = "modified"
	StatusUnchanged = "unchanged"
)
