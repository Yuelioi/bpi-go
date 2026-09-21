# bpi-go

## 相关项目

- [bpi-rs](https://github.com/Yuelioi/bpi-rs)：Rust 版本。
- [bpi-py](https://github.com/Yuelioi/bpi-py)：Python 版本。

一个面向 Go 的 Bilibili API SDK。

如果你想在 Go 里获取视频信息、搜索、用户资料、排行榜、直播、动态、评论、收藏夹、音频、番剧，或者调用登录态与创作中心接口，`bpi-go` 提供了一套统一、类型化的调用方式。

当前 v0.3.0 覆盖 **27 个领域**，面向 Go 1.25。

```go
package main

import (
	"context"
	"fmt"
	"log"

	bpi "github.com/Yuelioi/bpi-go"
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

	fmt.Printf("%s by %s\n", view.Title, view.Owner.Name)
}
```

## 这个项目从哪里来？

`bpi-go` 是 [`bpi-rs`](https://github.com/Yuelioi/bpi-rs) **v0.3.0** 的 Go 版本。

Go 版本延续了 `bpi-rs` 已经整理过的领域划分、接口语义、响应模型、错误语义和离线契约验证方式，同时按照 Go 的习惯重新组织为 `context.Context`、显式参数类型和可并发使用的客户端。

## 为什么用 bpi-go？

从使用者角度，主要是不用反复处理 Bilibili HTTP API 的签名、Cookie、响应包络和各种特殊格式。

| 你需要做的事 | `bpi-go` 帮你处理 |
| --- | --- |
| 调很多不同类型的 B 站接口 | 27 个领域都挂在同一个根客户端上 |
| 手工拼 URL 和查询参数 | 使用领域方法和类型化参数 |
| 自己解析 JSON | 常用响应直接返回 Go struct |
| 自己处理 WBI | SDK 统一处理签名 |
| 自己管理 Cookie / CSRF | Cookie / Account 显式传入，登录态集中管理 |
| 控制超时和取消请求 | 所有网络调用都接收 `context.Context` |
| 做服务端并发调用 | Client 实例彼此隔离，可安全并发使用 |
| 控制依赖 | 运行时只使用 Go 标准库 |
| 处理 XML、压缩和二进制响应 | SDK 已覆盖对应特殊协议路径 |
| 判断登录、VIP、权限和风控错误 | 提供稳定的错误类型与语义判断函数 |

调用方式统一按领域组织：

```go
client.Video().View(...)
client.Search().Videos(...)
client.User().Card(...)
client.Live().RoomInfo(...)
client.Dynamic().All(...)
client.Comment().List(...)
client.CreativeCenter().SeasonList(...)
```

这样不需要记一整套扁平函数名。

## 安装

```bash
go get github.com/Yuelioi/bpi-go@v0.3.0
```

或者直接使用最新兼容版本：

```bash
go get github.com/Yuelioi/bpi-go
```

## 示例

### 1. 获取视频信息

不需要登录。

```go
ctx := context.Background()

bvid, err := ids.ParseBVID("BV1xx411c7mD")
if err != nil {
	return
}

view, err := client.Video().View(ctx, video.ViewByBVID(bvid))
if err != nil {
	return
}

fmt.Println(view.Title)
fmt.Println(view.Owner.Name)
fmt.Println(view.Stat.View)
```

### 2. 搜索视频

```go
params, err := search.NewVideoParams("Python")
if err != nil {
	return
}

result, err := client.Search().Videos(ctx, params)
if err != nil {
	return
}

if result.Result != nil {
	for _, item := range *result.Result {
		fmt.Println(item.Title)
	}
}
```

搜索模块还提供文章、番剧、影视、用户、直播间等分类搜索。

### 3. 使用登录态

需要登录的接口显式传入 Cookie：

```go
client, err := bpi.NewClient(
	bpi.WithCookie("SESSDATA=...; bili_jct=...; DedeUserID=..."),
)
if err != nil {
	return
}

nav, err := client.Login().Nav(ctx)
if err != nil {
	return
}

fmt.Println(nav.IsLogin)
```

也可以使用结构化 `Account`：

```go
client, err := bpi.NewClient(bpi.WithAccount(bpi.Account{
	DedeUserID: "...",
	SESSDATA:   "...",
	BiliJCT:    "...",
	Buvid3:     "...",
}))
```

SDK 不会自动读取浏览器 Cookie、本地账号文件、环境变量或秘密管理系统。

### 4. 二维码登录

```go
generated, err := client.Login().GenerateQR(ctx)
if err != nil {
	return
}

fmt.Println(generated.URL)
fmt.Println(generated.Key)

params, err := login.NewQRPollParams(generated.Key)
if err != nil {
	return
}

status, err := client.Login().PollQR(ctx, params)
if err != nil {
	return
}

fmt.Println(status.Code)
```

SDK 负责请求和登录状态解析，二维码如何展示、多久轮询一次由你的应用决定。

## 覆盖了哪些模块？

当前提供以下 27 个领域客户端：

| 类型 | 模块 |
| --- | --- |
| 视频与内容 | `video`、`videoranking`、`bangumi`、`cheese`、`audio`、`article`、`note`、`opus`、`manga` |
| 用户与互动 | `user`、`comment`、`dynamic`、`message`、`fav`、`historytoview` |
| 直播 | `live` |
| 搜索 | `search` |
| 账号 | `login`、`vip`、`wallet`、`electric` |
| 创作者 | `creativecenter` |
| 其他 | `activity`、`clientinfo`、`danmaku`、`misc`、`webwidget` |

自动生成的 [API 对齐索引](docs/api-index.md) 记录 Rust 契约对应的 Go 访问器、方法、参数类型、响应模型和风险级别。

[迁移说明](docs/migration-from-bpi-rs.md) 介绍 Rust 与 Go 之间的语义对应关系。

## 类型化参数与响应

常用 ID 使用专门类型进行校验：

```go
bvid, err := ids.ParseBVID("BV1xx411c7mD")
```

接口参数也使用明确的参数结构：

```go
params := video.ViewByBVID(bvid)
view, err := client.Video().View(ctx, params)
```

这样很多无效参数可以在真正发出 HTTP 请求之前被发现。

响应直接返回业务模型：

```go
fmt.Println(view.Title)
fmt.Println(view.Owner.Name)
fmt.Println(view.Stat.View)
```

对于上游本身没有稳定结构的字段，会保留为 `json.RawMessage` 或开放结构，而不是猜测协议。

## context 与并发

所有网络操作都接收 `context.Context`：

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

view, err := client.Video().View(ctx, params)
```

客户端实例没有隐式全局账号状态，可以安全地为不同账号创建不同 Client，也可以在服务端并发复用同一个 Client。

## 错误处理

常见错误类型包括：

- `ParameterError`
- `TransportError`
- `HTTPError`
- `APIError`
- `ResponseDecodeError`
- `ResponseTooLargeError`
- `ErrMissingData`
- `ErrAuthenticationRequired`

同时提供稳定的语义判断：

```go
if bpi.RequiresLogin(err) {
	fmt.Println("需要登录")
}

if bpi.RequiresVIP(err) {
	fmt.Println("需要大会员")
}

if bpi.IsPermissionError(err) {
	fmt.Println("权限不足")
}

if bpi.IsRiskControl(err) {
	fmt.Println("触发风控")
}
```

## 自定义请求与模型恢复

领域客户端是主要入口。如果某个端点还没有独立领域方法，也可以复用同一套有界传输和响应包络：

```go
request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
if err != nil {
	return
}

payload, err := bpi.SendPayload[MyPayload](ctx, client, request, "my.endpoint")
```

当响应不再匹配现有模型时，可以通过 `errors.As` 从 `*bpi.ResponseDecodeError` 中显式取得响应体副本。响应可能包含账户或隐私数据，不应直接写入日志或提交到仓库。

## 仓库结构

```text
bpi.go, domains.go     根客户端和 27 个领域入口
client/                HTTP、会话、签名、响应和通用选项
activity/, video/, ... 各领域客户端、模型、参数和契约测试
ids/                   类型化 Bilibili ID
internal/contracttest/ 离线契约测试支持
docs/                  API 对齐、迁移和发布文档
```

大多数调用方只需要从 `bpi.NewClient()` 开始，再通过 `client.Video()`、`client.Search()` 等领域入口访问 API。

## 验证与开发

项目包含离线契约测试，普通测试不依赖真实 Bilibili 账号：

```bash
go vet ./...
go test -count=1 ./...
go test -race -count=1 ./...
```

实时 Probe 默认禁用，只面向开发维护场景，SDK 自身不会自动读取 Probe 的账号环境变量。

详细说明见：

- [API 对齐索引](docs/api-index.md)
- [迁移说明](docs/migration-from-bpi-rs.md)
- [Probe 安全与操作指南](docs/probe.md)
- [贡献指南](CONTRIBUTING.md)
- [安全策略](SECURITY.md)
- [发布流程](docs/releasing.md)
- [更新日志](CHANGELOG.md)

## 接口稳定性

Bilibili Web API 并不是官方稳定公开 API，上游接口、字段和错误码可能随时变化。

`bpi-go` 使用离线契约、脱敏响应样例和类型模型来降低升级成本，但这不代表第三方 Web API 本身具有长期稳定性。

## License

MIT License，见 [LICENSE](LICENSE)。
