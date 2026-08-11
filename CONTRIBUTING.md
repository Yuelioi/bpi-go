# 贡献指南

感谢你改进 `bpi-go`。本项目将锁定的 Rust 源码和经过脱敏、已纳入对齐范围的契约作为行为依据；方法能够编译并不代表端点已经实现完整。

## 修改 API 前

1. 在 `parity/source.json` 和 [API 对齐索引](docs/api-index.md) 中找到对应条目。
2. 阅读 `testdata/contracts/` 中复制的契约及所有相关响应配置。
3. 保持风险级别和认证边界不变。
4. 稳定字段应优先使用明确类型；只有真正不稳定的子树才使用 `json.RawMessage` 隔离。
5. 保持根 `bpi` 门面精简。领域客户端实现、参数、模型和契约测试应放在同一个领域包中；传输、会话、签名和响应机制应放在 `client/`。

尚未纳入对齐范围的新端点，以及写入或消费类端点，不属于首个版本的对齐范围。实现前必须进行明确的风险审查，提供经过脱敏的契约证据，并设置安全开关。

## 必需的垂直切片

一次端点变更通常必须同时包含：

- 带校验的参数构造；
- 接收 `context.Context` 的公开领域方法；
- 对请求方法、URL、查询参数、表单或请求体、请求头、CSRF 和 WBI 的断言；
- 对所有已纳入配置和错误结果的响应解码；
- 对齐映射以及复制并脱敏后的 fixture；
- 对不直观行为的公开说明。

测试使用 `http.RoundTripper` 和仓库内 fixture。默认测试必须保持离线，并且不得读取凭据。

## 验证

```powershell
gofmt -w (rg --files -g '*.go')
go vet ./...
go test -count=1 ./...
go test -race -count=1 ./...
go run ./cmd/bpi-probe audit
go run ./cmd/bpi-probe api-doc --check
```

依赖发生变化后运行 `go mod tidy`。优先使用标准库；新增任何依赖时，都应在变更说明中解释原因。

## 凭据与 fixture

严禁提交账户文件、Cookie 字符串、CSRF 值、令牌、原始私有响应或实时 Probe 响应体。请使用既有脱敏约定（`<redacted>`、确定性的 fixture 标识符和空私有列表）。提交 fixture 变更前运行 `bpi-probe sanitize-audit`。

## 提交与评审说明

说明受影响的契约名称、涉及的风险级别以及已执行的验证命令。明确标注有意保留的原始 JSON 边界，以及 Go API 为安全性或类型有效性而收紧的任何上游行为。
