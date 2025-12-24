# File Sync Monitor

一个用 Go 语言实现的文件同步监测工具，用于对比两个目录的文件差异，类似于 `git status` 的功能。

## 功能特性

- 🔍 **递归目录扫描**：自动递归扫描两个目录内的所有文件和子目录
- 📊 **文件元数据收集**：在不打开文件内容的前提下，收集以下信息：
  - 文件/目录名称
  - 文件大小（字节）
  - 修改时间
  - 文件权限
  - 是否为目录
- 🔄 **智能对比**：对比两个目录中相同相对路径的文件，检测差异
- 📝 **详细差异报告**：输出详细的差异信息，包括：
  - **新增文件**：仅在右侧目录存在的文件
  - **删除文件**：仅在左侧目录存在的文件
  - **修改文件**：两侧都存在但属性不同的文件（大小、修改时间、权限）
  - **未变更文件**：两侧完全一致的文件（可选显示）

## 项目结构

```
file_syn/
├── cmd/
│   └── file_syn/          # 主程序入口
│       └── main.go
├── config/                # 配置文件目录
│   └── config.json        # 配置文件示例
├── internal/              # 内部包（不对外暴露）
│   ├── config/           # 配置加载模块
│   │   ├── config.go
│   │   └── config_test.go
│   ├── scanner/          # 文件扫描模块
│   │   ├── scanner.go
│   │   └── scanner_test.go
│   ├── diff/             # 文件对比模块
│   │   ├── diff.go
│   │   └── diff_test.go
│   └── reporter/         # 结果输出模块
│       └── reporter.go
├── pkg/                   # 公共包
│   └── models/           # 数据模型
│       └── models.go
├── bin/                   # 编译输出目录
├── Makefile              # 构建脚本
├── go.mod                # Go模块文件
└── .gitignore           # Git忽略文件
```

## 安装

### 从源码编译

```bash
# 克隆项目
git clone <repository-url>
cd file_syn

# 配置目录路径（编辑 config/config.json）
# 复制示例配置文件
cp config/config.json.example config/config.json
# 编辑 config/config.json，设置 left_dir 和 right_dir

# 编译当前平台
make build

# 或使用 go 命令直接编译
go build -o bin/file_syn ./cmd/file_syn
```

### 多平台编译

项目支持编译 Linux、Windows 和 macOS 多个平台：

```bash
# 编译所有平台
make build-all

# 或单独编译特定平台
make build-linux    # Linux (amd64, arm64)
make build-windows  # Windows (amd64)
make build-darwin   # macOS (amd64, arm64)
```

编译后的可执行文件位于 `bin/` 目录：
- `file_syn-linux-amd64` - Linux x86_64
- `file_syn-linux-arm64` - Linux ARM64
- `file_syn-windows-amd64.exe` - Windows x86_64
- `file_syn-darwin-amd64` - macOS Intel
- `file_syn-darwin-arm64` - macOS Apple Silicon

## 配置

程序使用 JSON 配置文件来指定要对比的目录。首次使用前需要配置 `config/config.json`。

### 配置文件设置

1. 复制示例配置文件：
```bash
cp config/config.json.example config/config.json
```

2. 编辑 `config/config.json` 文件，设置要对比的目录路径：

```json
{
  "left_dir": "/path/to/left/directory",
  "right_dir": "/path/to/right/directory",
  "show_unchanged": false
}
```

配置项说明：
- `left_dir`: 左侧目录的路径（必填）
- `right_dir`: 右侧目录的路径（必填）
- `show_unchanged`: 是否显示未变更的文件（可选，默认为 false）

### 配置文件查找顺序

如果未指定配置文件路径，程序将按以下顺序查找：
1. `config/config.json`
2. `./config/config.json`
3. `config.json`

## 使用方法

### 基本用法

使用默认配置文件 `config/config.json`：

```bash
./bin/file_syn
```

### 指定配置文件

```bash
./bin/file_syn /path/to/custom-config.json
```

### 示例

```bash
# 使用默认配置文件
./bin/file_syn

# 使用自定义配置文件
./bin/file_syn /path/to/my-config.json
```

### 输出示例

```
配置文件: /home/node7/xpk/file_syn/config/config.json
左侧目录: /home/user/dir1
右侧目录: /home/user/dir2
正在扫描和对比...

╔════════════════════════════════════════════════════════════════════════════╗
║                        文件同步监测结果                                    ║
╚════════════════════════════════════════════════════════════════════════════╝

┌──────────┬──────────────────────────────────────────────────────────────┐
│ 状态       │ 文件路径                                                         │
├──────────┼──────────────────────────────────────────────────────────────┤
│ ➕ 新增     │ new_file.txt                                                 │
│          │   文件仅存在于右侧目录                                                 │
├──────────┼──────────────────────────────────────────────────────────────┤
│ ➖ 删除     │ old_file.txt                                                 │
│          │   文件仅存在于左侧目录                                                 │
├──────────┼──────────────────────────────────────────────────────────────┤
│ 🔄 修改     │ changed_file.txt                                            │
│          │   大小: 1.0 KB → 2.0 KB                                        │
│          │   修改时间: 2024-01-01 10:00:00 → 2024-01-01 11:00:00          │
└──────────┴──────────────────────────────────────────────────────────────┘

╔════════════════════════════════════════════════════════════════════════════╗
║                              统计信息                                       ║
╚════════════════════════════════════════════════════════════════════════════╝

┌──────────────────┬────────┐
│ 新增文件             │      1 │
├──────────────────┼────────┤
│ 删除文件             │      1 │
├──────────────────┼────────┤
│ 修改文件             │      1 │
├──────────────────┼────────┤
│ 总计               │      3 │
└──────────────────┴────────┘
```

## Makefile 命令

项目提供了丰富的 Makefile 命令来简化开发流程：

### 常用命令

- `make help` - 显示所有可用命令
- `make build` - 编译当前平台的可执行文件
- `make run` - 编译并运行程序（使用默认配置文件）
- `make run CONFIG=/path/to/config.json` - 编译并运行程序（使用指定配置文件）
- `make test` - 运行测试
- `make test-coverage` - 运行测试并生成覆盖率报告
- `make clean` - 清理编译产物

### 多平台编译

- `make build-all` - 编译所有平台（Linux/Windows/macOS，amd64/arm64）
- `make build-linux` - 仅编译 Linux 版本
- `make build-windows` - 仅编译 Windows 版本
- `make build-darwin` - 仅编译 macOS 版本

### 代码质量

- `make fmt` - 格式化代码
- `make vet` - 运行 go vet 代码静态检查
- `make lint` - 运行代码检查（fmt + vet）

### 其他

- `make deps` - 下载依赖
- `make mod-verify` - 验证依赖
- `make install` - 安装到系统路径（需要 sudo 权限）
- `make all` - 执行完整流程（清理、检查、测试、编译所有平台）

## 开发

### 运行测试

```bash
make test
```

### 生成测试覆盖率报告

```bash
make test-coverage
```

生成的覆盖率报告文件：`coverage.html`

### 代码格式化

```bash
make fmt
```

### 代码检查

```bash
make lint
```

## 技术细节

- **不读取文件内容**：只比较文件的元数据（大小、修改时间、权限等），性能高效
- **自动处理路径差异**：使用相对路径进行对比，不关心目录路径本身
- **错误处理**：遇到无法访问的文件会记录警告但继续扫描
- **统计信息**：输出包含详细的统计信息，方便快速了解差异情况

## 依赖

- Go 1.21 或更高版本

## 许可证

[根据项目实际情况添加许可证信息]

## 贡献

欢迎提交 Issue 和 Pull Request！

