# Probe 安全与操作指南

`cmd/bpi-probe` 是仓库中的确定性契约审计器，也是带显式开关的实时只读工具。离线命令绝不会加载账户文件，也不会执行网络请求。

## 离线审计

```powershell
go run ./cmd/bpi-probe audit
```

审计内容包括：

- 锁定的 `bpi-rs` 提交以及 27 个领域、206 条契约的清单；
- 115 条公开读取、32 条认证读取、52 条私有读取和 7 条登录会话契约的精确风险统计；
- 每条已纳入契约都有且只有一个 Go 映射，并且不存在多余映射；
- 契约元数据、账户配置用例、fixture 路径、JSON 有效性和预期 API 状态码；
- `parity/contracts.lock.json` 中全部 642 个快照文件的哈希；
- 已提交响应 fixture 中高可信度的凭据或个人数据泄露。

其他离线命令：

```powershell
go run ./cmd/bpi-probe sanitize-audit
go run ./cmd/bpi-probe api-doc --check
go run ./cmd/bpi-probe lock
```

`lock` 是维护操作。只有在完成有意的契约快照更新评审后才能运行，否则它会把未知漂移写入基准，而不是发现漂移。

## 实时只读 Probe

网络请求受两道独立开关控制，必须同时满足：

```powershell
$env:BPI_PROBE = "1"
go run ./cmd/bpi-probe batch-run --read-only --profiles anonymous
```

运行器接受 `anonymous`、`normal` 和 `vip` 配置，默认只运行 `anonymous`。认证配置通过临时环境变量显式提供完整 Cookie 请求头：

```powershell
$env:BPI_PROBE = "1"
$env:BPI_COOKIE_NORMAL = "DedeUserID=...; SESSDATA=...; bili_jct=...; buvid3=..."
$env:BPI_COOKIE_VIP = "DedeUserID=...; SESSDATA=...; bili_jct=...; buvid3=..."
go run ./cmd/bpi-probe batch-run `
  --read-only `
  --profiles normal,vip `
  --output probe-output/read-summary.json
Remove-Item Env:BPI_COOKIE_NORMAL, Env:BPI_COOKIE_VIP
```

`BPI_COOKIE_NORMAL` 和 `BPI_COOKIE_VIP` 只由 Probe 命令读取；`bpi` SDK 不读取环境变量。该命令只运行 `public-read`、`authenticated-read` 和 `private-read` 契约，并始终跳过 `login-session`；写入和消费类契约不受支持。删除 `--read-only` 不会扩大权限，而会让命令在发出请求前直接失败。

## 输出隐私

Probe 输出是有意不包含响应体的摘要。每条结果仅包含：

- 契约名称、账户配置和风险级别；
- HTTP 状态码和 API 状态码；
- 响应字节长度和 SHA-256 摘要；
- `passed`、`api_mismatch` 或 `transport_error` 等有限结果标签。

它绝不会写入响应体、请求 URL 或参数、请求头、Cookie、CSRF 值、账户标识或传输错误文本。`probe-output/` 已由 Git 忽略，作为第二道保护。不要把外部 HTTP 工具的详细输出重定向到仓库中。

## Cookie 输入

Probe 不再定义或解析账户文件格式。认证配置必须通过对应的 `BPI_COOKIE_NORMAL` 或 `BPI_COOKIE_VIP` 提供原始 HTTP Cookie 请求头值；运行结束后应立即清除环境变量。一旦凭据发生暴露，应立即轮换。

## 源码同步

本地存在锁定版本的同级 Rust 仓库时，可以验证生成清单和按字节比对的契约快照：

```powershell
go run ./cmd/bpi-sourcegen `
  -rust ../bpi-rs `
  -expect-commit 94bcf43e46d6e11b55ec4d36d4848692cecd4213 `
  -expect-domains 27 `
  -expect-contracts 206 `
  -check
```

源码同步器会将 JSON 文本的 CRLF 规范化为 LF，再进行逐字节比较和写入；非 JSON 文件保持原始字节。该规则使契约快照及其 SHA-256 锁在 Windows 与 Linux 的全新检出中保持一致。

对于已经明确评审的上游更新，去掉 `-check` 并添加 `-sync-contracts`。同步命令会写入已变化或缺失的文件，但拒绝删除过期文件；删除操作必须经过人工明确评审。完成后重新生成契约锁和 API 索引。
