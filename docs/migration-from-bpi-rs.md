# 从 bpi-rs 迁移

`bpi-go` 复现了 `bpi-rs` 0.2.4 提交 `36cb1104befee33b4281c59a4365d42dfdde45a4` 中已纳入对齐范围的行为。它是语义迁移，不是逐行源码翻译，也不是早期 Go 骨架的兼容层。

## 客户端构造

Rust:

```rust
let client = BpiClient::new()?;
let view = client.video().view(params).await?;
```

Go:

```go
client, err := bpi.NewClient()
if err != nil {
	return err
}
view, err := client.Video().View(ctx, params)
```

Go 领域客户端是由同一个隔离根客户端支撑的轻量值。应复用根客户端，而不是为每次请求重新构造。

根 `bpi` 包有意设计为兼容门面。具体实现与参数、模型一起放在 `activity`、`video` 等领域包中，通用传输和会话行为位于 `client`。通过 `bpi.NewClient()` 调用仍是首选的稳定入口；需要显式组合时也可以使用底层包。

## 参数与标识符

Rust 的 builder 参数对应 Go 的“构造函数加选项”值。构造函数负责设置必需标识符和默认值；`With...` 方法返回更新后的值，并在需要时报告无效的可选输入。

```go
bvid, err := ids.ParseBVID("BV1xx411c7mD")
if err != nil {
	return err
}
params := video.ViewByBVID(bvid)
```

数值标识符拒绝零值。AID/BVID 等互斥标识符通过不同构造函数表达，从源头避免无效组合。

## 结果与错误

普通领域方法直接返回业务数据，不要求每个调用点重复处理传输层包络。

可以结合公开错误类型、语义辅助函数和 `errors.As` 处理错误：

```go
value, err := client.Live().MyMedals(ctx, live.NewMyMedalsParams())
switch {
case err == nil:
	_ = value
case bpi.RequiresLogin(err):
	// 显式提供账户。
case bpi.IsRiskControl(err):
	// 主动退避，不要激进重试。
default:
	return err
}
```

相关类型包括 `ParameterError`、`TransportError`、`HTTPError`、`APIError`、`ResponseDecodeError` 和 `ResponseTooLargeError`。辅助函数覆盖登录、VIP、权限和风控结果。

## 凭据与会话

项目不存在全局账户，也不会隐式查找配置文件。可以选择：

- 使用 `bpi.WithCookie(cookieHeader)` 提供原始 HTTP Cookie 请求头值；
- 使用 `bpi.WithAccount(account)` 提供经过校验的结构化凭据；

Go 版本有意不移植 Rust 的账户文件加载器，也不约定 TOML、JSON 或环境变量格式。调用方从自己的配置层取得凭据后，再传给 SDK；因此模块只依赖 Go 标准库。

HTTP Cookie 请求头的格式是统一的 `name=value; name2=value2`，但不同端点所需的 Cookie 集合并不完全相同。`SESSDATA` 表示登录会话，`bili_jct` 用于 CSRF，`DedeUserID` 和 `buvid3` 常用于身份、设备识别或降低风控；SDK 还会原样保留其他合法 Cookie 对。

每个客户端独立持有 Cookie、WBI 缓存、日志器、时钟和 HTTP 配置。受信任哔哩哔哩主机返回的 `Set-Cookie` 只更新当前客户端。跨主机自定义请求不会携带哔哩哔哩凭据。

## 异步与取消

Rust future 对应接收 `context.Context` 的同步 Go 方法。Go 调用会阻塞到完成、取消或超时；需要并发时由调用方选择 goroutine。

## 可选值与原始响应

成功数据可能为 `null` 的端点会返回指针，或返回文档中明确说明的零值、nil 集合。二进制和 XML 接口使用 `[]byte` 或格式专用模型，不会强行经过 JSON。

模型发生漂移时，只能通过 `ResponseDecodeError.Body()` 恢复原始字节。即使端点被归类为公开接口，也应将这些字节视为敏感数据。

## 命名对应关系

Rust 的 snake_case 领域名称转换为 Go 访问器，例如：

| Rust | Go |
| --- | --- |
| `client.video()` | `client.Video()` |
| `client.video_ranking()` | `client.VideoRanking()` |
| `client.web_widget()` | `client.WebWidget()` |
| `client.creativecenter()` | `client.CreativeCenter()` |
| `client.historytoview()` | `client.HistoryToView()` |

自动生成的 [API 对齐索引](api-index.md) 是逐方法对应关系的权威记录。

## 范围边界

首个 Go 版本包含全部已纳入对齐范围的契约，包括认证读取、私有读取、登录二维码与会话流程、WBI 读取、压缩 XML 和 protobuf 字节响应。尚未纳入的 Rust 自由函数、写入、购买操作和不完整的流式占位接口不会对外提供；它们需要单独的风险审查和新的契约证据。
