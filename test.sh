#!/bin/bash

# 文件同步监测工具测试脚本

set -e

echo "=========================================="
echo "文件同步监测工具 - 测试脚本"
echo "=========================================="
echo ""

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# 清理之前的测试目录
echo "1. 清理之前的测试目录..."
rm -rf test_dirs
mkdir -p test_dirs/left test_dirs/right

# 创建测试文件
echo "2. 创建测试文件..."

# 相同内容的文件
echo "相同内容" > test_dirs/left/same.txt
echo "相同内容" > test_dirs/right/same.txt

# 左侧独有文件
echo "左侧独有文件" > test_dirs/left/left_only.txt

# 右侧独有文件
echo "右侧独有文件" > test_dirs/right/right_only.txt

# 内容不同的文件（大小不同）
echo "左侧内容 - 修改前" > test_dirs/left/changed.txt
echo "右侧内容 - 修改后，内容更长" > test_dirs/right/changed.txt

# 子目录测试
mkdir -p test_dirs/left/subdir test_dirs/right/subdir
echo "子目录文件" > test_dirs/left/subdir/subfile.txt
echo "子目录文件" > test_dirs/right/subdir/subfile.txt

# 左侧子目录独有
echo "左侧子目录独有" > test_dirs/left/subdir/left_sub.txt

# 更新配置文件
echo "3. 更新配置文件..."
cat > config/config.json << EOF
{
  "left_dir": "$SCRIPT_DIR/test_dirs/left",
  "right_dir": "$SCRIPT_DIR/test_dirs/right",
  "show_unchanged": false
}
EOF

echo "配置文件已更新"
echo ""

# 运行测试（不显示未变更文件）
echo "4. 运行测试（不显示未变更文件）..."
echo "----------------------------------------"
./bin/file_syn
echo ""

# 运行测试（显示未变更文件）
echo "5. 运行测试（显示未变更文件）..."
echo "----------------------------------------"
cat > config/config.json << EOF
{
  "left_dir": "$SCRIPT_DIR/test_dirs/left",
  "right_dir": "$SCRIPT_DIR/test_dirs/right",
  "show_unchanged": true
}
EOF
./bin/file_syn
echo ""

# 验证文件列表
echo "6. 验证文件列表..."
echo "左侧目录文件:"
find test_dirs/left -type f | sort
echo ""
echo "右侧目录文件:"
find test_dirs/right -type f | sort
echo ""

echo "=========================================="
echo "测试完成！"
echo "=========================================="

