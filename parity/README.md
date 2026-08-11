# 源码与实现对齐

`parity` 目录将首个 Go 公开接口锁定到 `bpi-rs` 0.2.4 的提交 `36cb1104befee33b4281c59a4365d42dfdde45a4`。

- `source.json` 是 27 个 Rust 领域客户端和 206 条已纳入契约的确定性清单。
- `source.lock.json` 记录 Rust 仓库、提交版本和 Cargo 版本。
- `implemented.json` 将每条已纳入契约映射到一个 Go 领域访问器、方法、参数类型和响应模型。
- `contracts.lock.json` 记录 `testdata/contracts` 下每个文件的大小和 SHA-256 摘要。

当前快照已完成 206/206 条映射，其中包括 115 条公开读取、32 条认证读取、52 条私有读取和 7 条登录会话契约。

## 离线验证

主要检查不需要同级源码仓库，也不需要网络：

```powershell
go run ./cmd/bpi-probe audit
```

该检查会拒绝缺失、重复或多余的映射，身份与风险漂移，缺失或格式错误的 fixture，API 状态码不匹配，快照哈希变化，以及高可信度的秘密残留。

自动生成的文档有单独的确定性检查：

```powershell
go run ./cmd/bpi-probe api-doc --check
```

## 与 bpi-rs 对照验证

本地存在锁定版本的同级 Rust 仓库时：

```powershell
go run ./cmd/bpi-sourcegen `
  -rust ../bpi-rs `
  -expect-commit 36cb1104befee33b4281c59a4365d42dfdde45a4 `
  -expect-domains 27 `
  -expect-contracts 206 `
  -check
```

该命令会按字节比较 `source.json`、`source.lock.json` 和全部 642 个已复制的契约及 fixture 文件，同时拒绝使用存在未提交相关变更的 Rust 源码。

## 有意更新上游基准

评审新的 Rust 提交及其风险变化后，重新生成清单，并复制已变化或缺失的快照文件：

```powershell
go run ./cmd/bpi-sourcegen `
  -rust ../bpi-rs `
  -expect-commit <reviewed-commit> `
  -expect-domains <reviewed-count> `
  -expect-contracts <reviewed-count> `
  -sync-contracts
```

同步命令绝不会删除目标目录中的过期文件。请明确评审并删除这些文件，并且只在同一次已评审迁移中更新硬编码源码锁。最后执行：

```powershell
go run ./cmd/bpi-probe lock
go run ./cmd/bpi-probe api-doc
go run ./cmd/bpi-probe audit
```

不要为了掩盖无法解释的漂移而重新生成锁文件。
