# 更新日志

本文件记录项目中所有值得关注的变更。格式遵循 Keep a Changelog，版本号遵循语义化版本规范。

## [Unreleased]（未发布）

## [0.3.0] - 2026-08-27

### 变更

- 将对齐基准升级到 `bpi-rs` 0.3.0 标签提交 `94bcf43e46d6e11b55ec4d36d4848692cecd4213`；27 个领域、206 条契约、风险分类和 642 个契约快照文件均保持不变。
- **破坏性变更**：`video.DASHFLAC.Audio` 从 `[]video.DASHStream` 改为 `*video.DASHStream`，与上游 FLAC 单对象或 `null` 形态一致，并新增可选 `Display` 字段。
- **破坏性变更**：`bangumi.Detail.Total` 从 `uint32` 改为 `int32`，保留未知集数的 `-1` 哨兵。

### 修复

- 非零 API code 携带不兼容的部分 `data` 时优先返回 `APIError`，不再被成功模型的 `ResponseDecodeError` 遮蔽；兼容的非零 payload 仍可通过显式 `IntoData` 读取。
- 兼容历史动态中的数字 ID，以及已删除转发原动态的 `null` ID 占位对象。
- 修正番剧和课程播放流把当前画质 `durl` 误读为画质分组 `durls` 的问题。
- 兼容 FLAC 单音轨/空音轨、DASH 与 MP4 的 `null` 备用地址、番剧精简统计的 `coin`/`play` 字段，以及未开播和历史番剧的缺省字段。

## [0.1.0] - 2026-08-11

### 新增

- 符合 Go 习惯且相互隔离的客户端基础设施，支持显式凭据、有界 HTTP 执行、上下文传播、敏感日志清理、类型化错误和可恢复的解码失败。
- 类型化标识符、Cookie 会话与刷新、带客户端级 single-flight 缓存的 WBI 签名、Bili-ticket 签名和通用请求策略。
- 完整实现 `bpi-rs` 0.2.4 提交 `36cb1104befee33b4281c59a4365d42dfdde45a4` 中的 27 个领域客户端和 206 条已纳入对齐范围的契约。
- 对 115 条公开读取、32 条认证读取、52 条私有读取和 7 条登录会话契约进行离线请求与响应验证。
- 支持弹幕相关契约所需的 XML、raw-deflate 和 protobuf 字节响应。
- 可复现的源码清单、按字节锁定的契约快照、完整的 Go 对齐清单和自动生成的 API 索引。
- 不保存响应体、带双重开关且只读的 `bpi-probe`，支持匿名、普通和 VIP 配置，并提供隐私与契约审计。
- CI、迁移说明、贡献与安全策略，以及可重复执行的发布检查清单。

### 变更

- 将根 `bpi` 包收敛为兼容门面，通用客户端能力归入 `client/`，27 个领域实现和契约测试分别归入对应领域目录。
- 将 README、更新日志及维护说明统一为中文优先，并由生成器输出中文 API 对齐索引。
- 移除 `LoadAccountProfile`、`AccountProfile` 和 TOML 解析依赖；配置文件、环境变量及秘密管理由调用方负责，SDK 只接收显式 `Account` 或原始 Cookie 请求头。
- 实时 Probe 改为通过 `BPI_COOKIE_NORMAL` 和 `BPI_COOKIE_VIP` 接收可选测试凭据，不再依赖账户文件格式。

### 修复

- 修复契约快照锁在全新检出时受 CRLF/LF 差异影响的问题，并确保 CI 使用 Go 1.25 系列的最新补丁版本。

### 安全

- 哔哩哔哩 Cookie 不会附加到无关主机。
- 包含秘密的查询字段和原始响应体不会进入日志或普通错误格式。
- 实时 Probe 无法执行登录会话、写入或消费类契约。

[Unreleased]: https://github.com/Yuelioi/bpi-go/compare/v0.3.0...HEAD
[0.3.0]: https://github.com/Yuelioi/bpi-go/compare/v0.1.0...v0.3.0
[0.1.0]: https://github.com/Yuelioi/bpi-go/releases/tag/v0.1.0
