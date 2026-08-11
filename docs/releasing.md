# 发布流程

首个对齐版本基于 `bpi-rs` 0.2.4 的提交 `36cb1104befee33b4281c59a4365d42dfdde45a4`。任何发布都不得在未明确说明的情况下移动该基准。

## 前置条件

- 工作树中只包含已经评审的发布变更。
- `go.mod` 和 `go.sum` 已整理完毕。
- `CHANGELOG.md` 已记录本次发布。
- `README.md`、包文档、迁移说明和自动生成的 API 索引与公开接口一致。
- 暂存区中不存在本地凭据文件、Probe 响应体、Cookie、令牌或私有 fixture。

## 验证

在仓库根目录执行：

```powershell
gofmt -w (rg --files -g '*.go')
go mod tidy
go vet ./...
go test -count=1 ./...
go test -race -count=1 ./...
go run ./cmd/bpi-probe audit
go run ./cmd/bpi-probe api-doc --check
govulncheck ./...
```

如果本地存在锁定版本的同级 Rust 仓库，还应执行 [Probe 安全与操作指南](probe.md) 中记录的 `bpi-sourcegen -check` 命令。

## 版本管理

版本标签遵循语义化版本规范。在 `v1.0.0` 之前，次版本可以调整公开模型和构造器，但凭据隔离、上下文行为、错误脱敏和安全开关属于兼容性承诺。Rust 源码提交、已纳入契约集合、风险级别或 fixture 快照发生任何变化，都必须在更新日志中明确说明。

## 创建标签

1. 确认发布提交上的 CI 全部通过。
2. 最后检查一次暂存文件列表和生成产物。
3. 创建带注释的标签，例如 `v0.1.0`。
4. 推送提交和标签。
5. 根据更新日志中的对应章节编写发布说明。

该模块不需要生成二进制发布产物；使用者直接引用带标签的 Go 模块。Probe 和源码生成器可以通过 `go install` 从同一标签安装。
