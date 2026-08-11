# bpi-go

`bpi-go` 是 [`bpi-rs`](https://github.com/Yuelioi/bpi-rs) 的 Go 版本，用于访问哔哩哔哩 HTTP API。首个版本锁定到 `bpi-rs` 0.2.4 的提交 `36cb1104befee33b4281c59a4365d42dfdde45a4`，完整实现了 27 个领域客户端和全部 206 条已纳入对齐范围的契约。

本项目遵循常见的 Go 工程实践：

- 客户端实例彼此隔离，可安全并发使用；
- 所有网络操作都接收 `context.Context`；
- 账户或 Cookie 必须显式配置，不会隐式读取文件；
- 凭据只会发送到受信任的哔哩哔哩主机；
- 响应缓冲区有明确上限，`log/slog` 默认静默并会清理敏感信息；
- 提供类型化参数、响应模型、错误、WBI 签名和 Cookie 刷新；
- 仅使用 Go 标准库，不引入第三方运行依赖；
- 契约测试完全离线且可复现，覆盖 JSON、XML、压缩数据和二进制响应。

## 安装

```bash
go get github.com/Yuelioi/bpi-go
```

当前模块面向 Go 1.25。

## 快速开始

```go
package main

import (
	"context"
	"log"

	"github.com/Yuelioi/bpi-go"
	"github.com/Yuelioi/bpi-go/ids"
	"github.com/Yuelioi/bpi-go/video"
)

func main() {
	client, err := bpi.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	bvid, err := ids.ParseBVID("BV1xx411c7mD")
	if err != nil {
		log.Fatal(err)
	}
	view, err := client.Video().View(context.Background(), video.ViewByBVID(bvid))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s (%s)", view.Title, view.BVID)
}
```

## 身份认证

凭据始终由调用方显式提供。最通用的入口是浏览器或其他凭据来源提供的原始 HTTP `Cookie` 请求头值：

```go
cookieHeader := os.Getenv("BPI_COOKIE") // 仅为示例，来源由调用方决定
client, err := bpi.NewClient(bpi.WithCookie(cookieHeader))
```

如果自己的配置层已经把常用字段解析为结构化数据，也可以使用 `Account`：

```go
client, err := bpi.NewClient(bpi.WithAccount(bpi.Account{
	DedeUserID: "...",
	SESSDATA:   "...",
	BiliJCT:    "...",
	Buvid3:     "...",
}))
```

SDK 不提供 `LoadCookie` 或账户配置文件加载器，也不会读取文件、环境变量或秘密管理系统。调用方负责取得凭据，再把 Cookie 请求头或 `Account` 传入客户端。请勿提交 Cookie、账户文件或未经处理的私有 API 响应。

## 自定义请求与模型恢复

领域客户端是稳定的主要接口。对于尚未建模的端点，调用方仍可复用同一套有界传输和响应包络处理：

```go
request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
if err != nil {
	return err
}
payload, err := bpi.SendPayload[MyPayload](ctx, client, request, "my.endpoint")
```

当响应不再匹配现有模型时，可以通过 `errors.As` 从 `*bpi.ResponseDecodeError` 中取得响应体的私有副本。该内容可能包含账户或隐私数据，只应在本地检查，不得记录到日志或提交到仓库。普通错误文本和结构化日志会有意省略它。

## 领域与接口对齐

目前支持 `activity`、`article`、`audio`、`bangumi`、`cheese`、`clientinfo`、`comment`、`creativecenter`、`danmaku`、`dynamic`、`electric`、`fav`、`historytoview`、`live`、`login`、`manga`、`message`、`misc`、`note`、`opus`、`search`、`user`、`video`、`video_ranking`、`vip`、`wallet` 和 `web_widget` 等领域。

自动生成的 [API 对齐索引](docs/api-index.md) 记录每条 Rust 契约对应的 Go 访问器、方法、参数类型、响应模型和风险级别。[迁移说明](docs/migration-from-bpi-rs.md) 介绍 Rust 与 Go 之间的语义对应关系。

## 仓库结构

根 `bpi` 包是一个轻量兼容门面，负责客户端构造、公开别名和 27 个领域访问器。通用的 HTTP、会话、签名、响应和选项逻辑位于 `client/`。每个领域目录集中存放该领域的客户端实现、参数、模型和契约测试：

```text
bpi.go, domains.go     根兼容门面
client/                通用客户端基础设施
activity/, video/, ... 领域客户端、模型、参数和测试
ids/                   经过校验的哔哩哔哩标识符
internal/contracttest/ 共享的离线契约测试支持
```

大多数调用方应继续使用 `bpi.NewClient()`，再通过返回的根客户端访问各领域。底层 `client` 和领域构造函数保持公开，供需要显式组合的集成场景使用。

本项目使用作者维护的 [Flightdeck](https://github.com/Yuelioi/flightdeck) 管理长期开发计划和跨会话工作记录；对应的普通 Markdown 工作台保存在仓库的 `flightdeck/` 目录中。Flightdeck 仅用于开发协作，不是 `bpi-go` 的运行依赖。

## 安全验证与 Probe

以下命令均为离线操作：

```powershell
go run ./cmd/bpi-probe audit
go run ./cmd/bpi-probe api-doc --check
go vet ./...
go test -count=1 ./...
go test -race -count=1 ./...
```

实时 Probe 默认禁用，只支持读取类契约，并且必须同时通过环境变量和命令行两道开关：

```powershell
$env:BPI_PROBE = "1"
go run ./cmd/bpi-probe batch-run --read-only --profiles anonymous
```

运行普通账户或 VIP 账户的 Probe 时，需要分别显式设置 `BPI_COOKIE_NORMAL` 或 `BPI_COOKIE_VIP`。这些变量只由开发工具读取，SDK 本身不会读取环境变量。Probe 摘要绝不包含响应体、Cookie、CSRF 值、请求参数或账户标识。详情参见 [Probe 安全与操作指南](docs/probe.md)。

## 项目文档

- [贡献指南](CONTRIBUTING.md)
- [安全策略](SECURITY.md)
- [发布流程](docs/releasing.md)
- [更新日志](CHANGELOG.md)
- [源码对齐与同步](parity/README.md)

本项目采用 [MIT 许可证](LICENSE)。
