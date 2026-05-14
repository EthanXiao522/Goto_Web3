# Web3 Infra 3 个月冲刺行动计划

> **前置文档**：[Web3 Infra 转型规划与行动指南](./web3_infra_handbook.md)  
> **学习强度**：每天 6 小时 × 每周 6 天（周一至周六），周日为轻量复习/弹性日  
> **总时长**：12 周 ≈ 432 小时  
> **核心目标**：3 个月后具备可演示的 GitHub 项目 + 面试能力

[TOC]

---

## 总览

### 三阶段划分

| 阶段 | 周次 | 主题 | 核心交付 |
|------|------|------|----------|
| **Phase 1** | 第 1-4 周 | Go 深化 + 链基础 + 简易 Listener | 可运行的 HTTP/WS 链监听器 |
| **Phase 2** | 第 5-8 周 | Kafka + 钱包系统（充值/提现/签名） | 完整钱包后端服务 |
| **Phase 3** | 第 9-12 周 | K8S 部署 + 监控 + 面试准备 | 生产级 GitHub 项目 + 面试就绪 |

### 每日作息建议

| 时段 | 时长 | 内容 |
|------|------|------|
| 09:00 - 11:00 | 2h | 阅读/视频学习（核心概念、文档、源码阅读） |
| 11:00 - 12:00 | 1h | 笔记整理 + 代码片段练习 |
| 14:00 - 17:00 | 3h | 项目编码（当天核心开发任务） |
| 晚上 | 弹性 | 如果当天任务提前完成，可自由探索或复习 |

### 周日安排（弹性日）

- 上午（1-2h）：回顾本周笔记，补完未完成的任务
- 下午：自由探索（阅读行业博客、逛 GitHub Trending、看 Devcon 演讲）

---

## Phase 1：基础夯实（第 1-4 周）

> **阶段目标：** Go 并发编程熟练化，理解 EVM 核心数据模型（账户/交易/区块/事件），能独立写出一个可运行的链监听器。
>
> **阶段产出：**
>
> - `chain-listener/` 项目仓库（WebSocket 订阅 + HTTP fallback + Reorg 处理）
> - Kafka producer/consumer 初步集成
> - 去重与幂等写入机制
> - README 架构图 + 快速开始指南
> - tag v0.2.0

### 第 1 周：Go 并发 + EVM 基础

> **本周目标：** 掌握 goroutine/channel/context 并发模式，理解以太坊账户/交易/区块基础模型，能通过 go-ethereum 连接 RPC 节点并查询链上数据。
>
> **本周产出：** Go 并发代码练习集（producer-consumer / worker pool / pipeline），RPC 调用 demo（查余额 / 查交易 / 查区块 / 解码 Event），概念笔记 1 篇。

#### Day 1（周一）— Go 并发模型
- [ ] 阅读 [Go Tour - Concurrency](https://go.dev/tour/concurrency/1)（2h）
- [ ] 浏览 [Go 官方文档](https://go.dev/doc/) 了解标准库索引（碎片时间）
- [ ] 动手写代码（3h）：
  - goroutine + channel 实现生产者-消费者模型
  - 使用 `context.WithCancel` 实现优雅退出
  - 使用 `sync.WaitGroup` 等待所有 worker 完成
- [ ] 笔记：画出 goroutine 调度模型 GMP 的草图（1h）

#### Day 2（周二）— Go 并发进阶
- [ ] 阅读 [Go Concurrency Patterns](https://go.dev/blog/pipelines)（1.5h）
- [ ] 动手写代码（3.5h）：
  - 实现 worker pool（固定数量 goroutine 处理任务队列）
  - 实现 fan-in / fan-out 模式
  - 使用 `sync.Mutex` / `sync.RWMutex` 保护共享状态
  - 使用 `sync.Map` 并理解其适用场景（key 写一次读多次）
- [ ] 笔记：对比 Java 线程池与 Go goroutine 的异同（1h）
- [ ] 浏览参考资源（碎片时间）：
  - [Uber Go Style Guide](https://github.com/uber-go/guide) — 了解 Go 代码规范
  - [JustForFunc](https://youtube.com/@JustForFunc) — 选看 1 集 goroutine 相关的实战视频
  - [Ardan Labs](https://youtube.com/@ardanlabs) — 收藏频道，后续关注 Bill Kennedy 的生产级 Go 训练

#### Day 3（周三）— Ethereum 账户与交易模型
- [ ] 阅读 [Ethereum Accounts](https://ethereum.org/en/developers/docs/accounts/)（1h）
- [ ] 阅读 [Ethereum Transactions](https://ethereum.org/en/developers/docs/transactions/)（1h）
- [ ] 动手（3h）：
  - 用 go-ethereum 的 `ethclient` 连接 Infura/Alchemy
  - 查询账户余额 `BalanceAt()`
  - 查询交易详情 `TransactionByHash()`
  - 字段理解：nonce、gasPrice、gasLimit、value、data、to
- [ ] 笔记：画出交易从签名到上链的完整生命周期（1h）

#### Day 4（周四）— Gas 机制 + 区块结构
- [ ] 阅读 [Ethereum Gas](https://ethereum.org/en/developers/docs/gas/)（1h）
- [ ] 阅读 [Ethereum Blocks](https://ethereum.org/en/developers/docs/blocks/)（1h）
- [ ] 动手（3h）：
  - 用 `BlockByNumber()` 拉取最新区块
  - 解析 block header（parentHash、stateRoot、txRoot、receiptRoot）
  - 遍历 block 中所有 transactions
  - 计算区块中所有交易的 gas 消耗总和
- [ ] 笔记：EIP-1559 的 baseFee + priorityFee 机制（1h）

#### Day 5（周五）— JSON-RPC 深入
- [ ] 阅读 [Ethereum JSON-RPC 规范](https://ethereum.org/en/developers/docs/apis/json-rpc/)（1h）
- [ ] 动手（4h）：
  - 直接构造 JSON-RPC 请求调用 `eth_getBlockByNumber`
  - 调用 `eth_getTransactionReceipt` 获取交易回执
  - 调用 `eth_getLogs` 按地址/topics 过滤事件日志
  - 调用 `eth_blockNumber` 和 `eth_syncing` 了解节点状态
- [ ] 笔记：整理常用 RPC 方法速查表（1h）

#### Day 6（周六）— Event 与 ABI
- [ ] 阅读 [Ethereum Events & Logs](https://docs.soliditylang.org/en/latest/contracts.html#events)（1h）
- [ ] 阅读 ABI 规范（1h）
- [ ] 动手（3h）：
  - 手写一个 ERC-20 Transfer 事件的 topic hash 计算（keccak256）
  - 解析 log.data 和 log.topics
  - 用 `abi.JSON.Unpack()` 解码 event log
  - 解析一个真实的 USDT Transfer 交易
- [ ] 笔记：topics 索引规则（indexed 参数 → topic，非 indexed → data）（1h）

#### Day 7（周日）— 弹性
- [ ] 补完本周未完成的任务
- [ ] 回顾所有笔记，建立一个 Go + EVM 概念脑图
- [ ] 推荐浏览：
  - [Finematics](https://youtube.com/@Finematics) — 选 1-2 个 DeFi 概念视频
  - [Cyfrin Updraft](https://updraft.cyfrin.io/) — 注册并浏览课程目录，了解免费区块链课程体系
  - [ETHGlobal](https://youtube.com/@ETHGlobal) — 选看 1 个 Workshop，感受开发氛围

---

### 第 2 周：Go 网络编程 + 区块扫描器（HTTP 轮询版）

> **本周目标：** 掌握 WebSocket/gRPC 网络编程，实现一个 HTTP 轮询版的区块扫描器（拉取区块 → 解析交易/事件 → 写入 PostgreSQL）。
>
> **本周产出：** `chain-listener/` 项目骨架，HTTP 轮询核心逻辑，PostgreSQL schema（blocks / transactions / event_logs），支持 event log 解析（ERC-20 Transfer）。

#### Day 1（周一）— Go WebSocket 客户端
- [ ] 阅读 [gorilla/websocket](https://pkg.go.dev/github.com/gorilla/websocket) 文档（1h）
- [ ] 动手（4h）：
  - 实现 WebSocket 连接到 Infura/Alchemy wss endpoint
  - 发送 `eth_subscribe` 订阅 newHeads
  - 接收并解析 newHeads 消息
  - 实现自动重连 + 指数退避（1s → 2s → 4s → 8s，max 60s）
- [ ] 笔记：WebSocket vs HTTP 轮询的优劣对比（1h）

#### Day 2（周二）— gRPC + Protobuf
- [ ] 阅读 [gRPC Go Quickstart](https://grpc.io/docs/languages/go/quickstart/)（1h）
- [ ] 动手（4h）：
  - 定义 .proto 文件（Block、Transaction 等消息）
  - 生成 Go 代码
  - 实现一个简单的 gRPC server + client
  - 理解为何链基础设施中 gRPC 比 REST 更常用（流式、二进制、强类型）
- [ ] 笔记：gRPC 四种通信模式（Unary / Server Stream / Client Stream / Bidirectional）（1h）

#### Day 3（周三）— 区块扫描器：项目骨架
- [ ] 初始化项目（1h）：
  - `go mod init github.com/<yourname>/chain-listener`
  - 搭建目录结构（cmd/、internal/、pkg/）
  - 配置管理（viper 或 envconfig）
- [ ] 编码（4h）：
  - 实现配置结构体（RPC URL、起始区块、确认数、数据库 DSN）
  - 实现 PostgreSQL 连接（database/sql + pgx）
  - 创建 blocks 表和 transactions 表（schema migration）
- [ ] 笔记：项目架构决策记录（1h）

#### Day 4（周四）— 区块扫描器：HTTP 轮询核心
- [ ] 编码（5h）：
  - 实现 HTTP 轮询循环：`eth_blockNumber` → 逐块拉取
  - 每个区块：`eth_getBlockByNumber(blockNum, true)` 获取含交易列表的完整区块
  - 解析后写入 blocks 表和 transactions 表
  - 实现 last_processed_block 的 checkpoint（写入 DB 或文件）
  - 重启后从 checkpoint + 1 继续
- [ ] 笔记：记录遇到的坑（1h）

#### Day 5（周五）— 区块扫描器：补全 + 日志
- [ ] 编码（4h）：
  - 实现 simple logging（使用 `log/slog` 结构化日志）
  - 添加 prometheus metrics（`processed_blocks_total`、`rpc_call_duration_seconds`）
  - 处理异常：RPC 超时重试、区块缺失时跳过并记录
- [ ] 跑起来，观察输出（2h）：
  - 从最新区块往前扫 100 个区块
  - 观察吞吐量：每秒能处理多少区块
  - 确认数据库写入正确

#### Day 6（周六）— Event Log 解析
- [ ] 编码（5h）：
  - 在 listener 中增加 event log 解析
  - 创建 event_logs 表（block_number, tx_hash, log_index, address, topic0, data）
  - 使用 `eth_getLogs` 按区块范围拉取
  - 初步尝试解析 ERC-20 Transfer（topic0 = `keccak256("Transfer(address,address,uint256)")`）
- [ ] 笔记：整理数据处理流水线的数据流图（1h）

#### Day 7（周日）— 弹性
- [ ] 运行自己的扫描器 2 小时，检查是否有数据遗漏
- [ ] 浏览 [go-ethereum 源码](https://github.com/ethereum/go-ethereum) ethclient 部分

---

### 第 3 周：Listener 升级（WebSocket + Reorg 处理）

> **本周目标：** 将 listener 切换为 WebSocket 实时订阅，实现 reorg 检测与回滚处理逻辑，实现去重与幂等写入。
>
> **本周产出：** WebSocket 订阅 + 自动重连（指数退避），reorg 检测/标记 is_canonical/重处理，gap scan 断点续扫，`tx_hash + log_index + block_hash` 复合唯一约束，集成测试。

#### Day 1（周一）— 切换为 WebSocket 订阅
- [ ] 改造 listener（5h）：
  - 将 HTTP 轮询核心替换为 WebSocket `eth_subscribe("newHeads")`
  - 订阅到新区块头后，通过 HTTP RPC 补拉完整区块（含交易列表）
  - 保留 HTTP 轮询作为 fallback（WS 断开时启用）
- [ ] 运行测试（1h）：对比 HTTP 版和 WS 版的数据一致性

#### Day 2（周二）— Reorg 处理：理论 + 设计
- [ ] 阅读 [Ethereum Reorgs](https://ethereum.org/en/developers/docs/blocks/#block-reorgs)（1h）
- [ ] 设计 reorg 处理方案（1h）：
  - 每写入一个区块，记录 `block_hash`
  - 再次遇到同一高度时比较 block_hash
  - 如果不同 → 触发 reorg 流程：标记旧数据 `is_canonical=false`
- [ ] 编码（3h）：
  - 在 blocks 表增加 `is_canonical` 字段（默认 true）
  - 实现 reorg 检测逻辑：`SELECT block_hash FROM blocks WHERE block_number = $1`

#### Day 3（周三）— Reorg 处理：编码实现
- [ ] 编码（5h）：
  - 实现 reorg 处理流程：
    1. 检测 block_hash 不一致
    2. 将旧链数据 `is_canonical` 标记为 false
    3. 重新解析新区块及其所有交易/日志
    4. 幂等写入新数据
  - 实现确认机制：区块 ≥ 12 个确认后，将 `is_canonical` 锁定
- [ ] 单元测试（1h）：

#### Day 4（周四）— WebSocket 断线恢复 + Gap Scan
- [ ] 编码（4h）：
  - 实现 checkpoint 持久化（last_processed_block + last_processed_hash）
  - WS 重连后从 checkpoint 恢复
  - 实现 gap scan：用 HTTP RPC 拉取 `[checkpoint, latest-12]` 区间，检查是否有遗漏
- [ ] 测试（2h）：
  - 手动断网测试重连
  - 故意 kill 进程后重启，检查从 checkpoint 恢复

#### Day 5（周五）— 去重与幂等性
- [ ] 编码（3h）：
  - transaction 表：`tx_hash` UNIQUE 约束，INSERT ... ON CONFLICT DO NOTHING
  - event_log 表：`(tx_hash, log_index, block_hash)` 复合 UNIQUE 约束
  - 确保 reorg 后重新处理不会产生重复数据
- [ ] 编写集成测试（3h）：
  - 模拟 reorg：手动往 DB 写入旧数据 → 运行 listener → 验证数据被标记/替换
  - 模拟重复事件：多次插入同一 tx → 验证幂等

#### Day 6（周六）— 代码整理 + 文档
- [ ] 代码整理（3h）：
  - 统一错误处理风格
  - 添加必要的注释（重点在 reorg 处理逻辑）
  - 函数拆分（单一职责原则）
- [ ] README 初稿（2h）：
  - 项目简介
  - 架构图（用 Mermaid 或图片）
  - 快速开始指南
- [ ] 跑 1 小时数据，统计处理速度和稳定性（1h）

#### Day 7（周日）— 弹性
- [ ] Phase 1 回顾：能否用自己的话讲清楚 "从区块到数据库" 的完整流程？
- [ ] 录制一段 5 分钟英文自述，尝试讲解自己的项目

---

### 第 4 周：Listener 打磨 + Phase 1 检查点

> **本周目标：** 性能优化（pprof + batch insert），完成 Kafka 集成（producer 发消息 + consumer 写库），通过 Phase 1 全部检查点。
>
> **本周产出：** pprof 性能分析报告，Kafka producer/consumer 代码，完善 README（架构图 + 运行指南），Phase 1 自检清单全部勾选，tag v0.2.0。

#### Day 1（周一）— 性能优化
- [ ] 使用 pprof 分析 listener 的性能瓶颈（2h）
- [ ] 优化（3h）：
  - 批量 INSERT（使用 PostgreSQL COPY 或 batch INSERT）
  - 减少不必要的 RPC 调用
  - 添加对象池（sync.Pool）减少 GC 压力
- [ ] 跑 benchmark，对比优化前后（1h）

#### Day 2（周二）— 错误处理与边界情况
- [ ] 编码（4h）：
  - RPC 返回 null 区块的处理
  - 节点落后（syncing）时的行为
  - pending block 的处理（只处理已出块）
  - 超大区块（如 1000+ tx）的处理
- [ ] 测试各种异常场景（2h）

#### Day 3（周三）— 为 Phase 2 做技术准备
- [ ] 阅读 Kafka 核心概念（2h）：
  - Topic、Partition、Consumer Group、Offset
  - [Kafka 官方文档](https://kafka.apache.org/documentation/#gettingStarted)
- [ ] 本地搭建 Kafka（Docker Compose）（1h）
- [ ] Go 连接 Kafka（3h）：
  - 使用 [segmentio/kafka-go](https://github.com/segmentio/kafka-go)
  - 实现 Producer 发送测试消息
  - 实现 Consumer 消费消息

#### Day 4（周四）— Listener 集成 Kafka Producer
- [ ] 编码（5h）：
  - 在现有 listener 中添加 Kafka Producer
  - 每解析完一个区块 → 发送 BlockMessage 到 Kafka
  - 消息格式：JSON 或 Protobuf
  - 确保消息不丢失（acks=all / 至少 ack=1）
- [ ] 测试（1h）：监听 → Kafka → 用命令行 consumer 验证消息

#### Day 5（周五）— Consumer + 简单存储
- [ ] 编码（5h）：
  - 实现独立的 consumer 服务：消费 Kafka BlockMessage → 写入 PostgreSQL
  - 实现 Consumer Group：启动 2 个 consumer 实例，验证自动负载均衡
  - 去重逻辑：与 listener 直写 DB 时的去重逻辑一致
- [ ] 端到端测试（1h）：Listener → Kafka → Consumer → DB

#### Day 6（周六）— Phase 1 收尾
- [ ] 补充单元测试和集成测试（3h）
- [ ] 完善 README（2h）：
  - 架构图（ASCII art 或 Mermaid）
  - 运行命令（make run-listener / make run-consumer）
  - 配置说明
- [ ] Phase 1 检查点自评（1h）

#### Day 7（周日）— Phase 1 检查点

**自检清单：**

- [ ] 能否用 go-ethereum ethclient 查询区块、交易、日志？
- [ ] 能否解释 goroutine + channel 的通信模式？
- [ ] Listener 是否支持 WebSocket 订阅 + 自动重连？
- [ ] Listener 是否处理 reorg（检测 + 回滚 + 重新处理）？
- [ ] 是否实现了 checkpoint 断点续扫？
- [ ] 是否具备幂等写入（重复数据不会导致错误）？
- [ ] Kafka producer/consumer 是否能正常工作？
- [ ] GitHub 项目 README 是否包含架构图和快速开始？

---

## Phase 2：核心系统（第 5-8 周）

> **阶段目标：** 在 Listener 基础上构建钱包后端核心系统（充值检测/提现/nonce管理/签名服务），完成两项目串联。
>
> **阶段产出：**
>
> - `wallet-backend/` 项目仓库（HD 钱包 / 充值 / 提现 / Nonce / 签名）
> - Nonce Manager（并发安全分配 + gap 检测 + replacement）
> - 充提全流程集成测试通过
> - docker-compose 一键启动全栈服务
> - 系统设计文档 2 篇
> - tag v0.5.0

### 第 5 周：Kafka 深化 + 生产级 Listener 重构

> **本周目标：** 深入 Kafka（partition / consumer group / 消息顺序），将 listener 消息格式升级为 Protobuf，实现数据一致性 reconciliation 工具。
>
> **本周产出：** Protobuf 消息定义（BlockMessage / TxMessage / LogMessage），多 partition + consumer group 并行消费，reconciliation 工具，钱包项目骨架创建。

#### Day 1（周一）— Kafka 深入
- [ ] 阅读（2h）：
  - Kafka 分区策略与消息顺序保证
  - Consumer Group rebalance 机制
  - [Confluent Kafka 101 视频](https://www.youtube.com/@Confluent)
- [ ] 动手（3h）：
  - 创建多 partition topic（如 4 partition）
  - 测试同一 key 的消息路由到同一 partition
  - 测试 consumer group rebalance
- [ ] 笔记：Kafka 消息顺序性 vs 并行度的权衡（1h）

#### Day 2（周二）— Listener Kafka 重构
- [ ] 编码（5h）：
  - 将 listener 的 Kafka 消息从 JSON 改为 Protobuf
  - 定义 proto：BlockMessage、TxMessage、LogMessage
  - 按区块高度 hash 做 partition key（保证同一区块的消息顺序）
  - Consumer 端按 tx_hash 分片，避免同一交易被多个 consumer 重复处理
- [ ] 笔记：为什么 Web3 场景中 Kafka 比 Redis Streams 更适合（1h）

#### Day 3（周三）— Consumer 业务逻辑
- [ ] 编码（5h）：
  - Consumer 端增加业务处理：识别 ERC-20 Transfer 事件
  - 建立地址索引表 address_tx_history
  - 去重 + 幂等处理的封装
- [ ] 测试（1h）：写入 100 个区块数据，验证 consumer 正确解析

#### Day 4（周四）— 数据一致性验证
- [ ] 编码（3h）：
  - 实现 reconciliation 工具：
    - 对比 blocks 表与 RPC 节点的区块总数
    - 对比 transactions 表与区块中的交易数
    - 发现不一致时告警
- [ ] 编写简单的前端或 CLI dashboard（3h）：
  - 显示最新处理的区块号
  - 显示 reorg 次数
  - 显示消息处理延迟

#### Day 5（周五）— 跑通完整数据链 + 性能测试
- [ ] 运行 listener + Kafka + consumer 全链路 2 小时（3h）
- [ ] 监控指标（2h）：
  - 区块处理延迟（从链上出块到写入 DB 的时间）
  - Kafka 消息积压情况
  - consumer lag
- [ ] 记录并优化瓶颈（1h）

#### Day 6（周六）— 总结 + 开始钱包系统设计
- [ ] 整理 listener 项目的架构文档（2h）
- [ ] 钱包系统设计（3h）：
  - 在纸上画出完整架构：API → Wallet Service → Deposit/Withdraw → Nonce Manager → Signing Service → RPC
  - 确定各组件间的数据流和接口
- [ ] 创建钱包项目骨架（1h）：`github.com/<yourname>/wallet-backend`

#### Day 7（周日）— 弹性
- [ ] 回顾 Kafka 核心概念，确认能讲清楚 "exactly-once vs at-least-once"
- [ ] 浏览 [Fireblocks 架构博客](https://www.fireblocks.com/blog/) 了解托管钱包设计思路

---

### 第 6 周：钱包核心（HD 钱包 + 签名 + 密钥管理）

> **本周目标：** 理解并实现 BIP-32/39/44 HD 钱包推导，掌握 EIP-1559/EIP-712 交易和消息签名，实现 keystore 加密存储方案。
>
> **本周产出：** HD 钱包生成/导入/派生地址代码，EIP-1559 + EIP-712 签名实现，keystore（AES-256-GCM 加密）存储，钱包服务 HTTP API（创建/查询/签名），安全加固（鉴权/审计日志/脱敏）。

#### Day 1（周一）— HD 钱包原理
- [ ] 阅读（2h）：
  - [BIP-32](https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki)（HD 钱包）
  - [BIP-39](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki)（助记词）
  - [BIP-44](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki)（多币种层级）
- [ ] 动手（3h）：
  - 用 Go 库（如 `go-ethereum/accounts` 或 `tyler-smith/go-bip39`）生成助记词
  - 从助记词 → seed → 主私钥 → 子私钥
  - 按 BIP-44 路径 `m/44'/60'/0'/0/0` 推导以太坊地址
- [ ] 笔记：画出 HD 钱包的密钥派生树（1h）
- [ ] 浏览参考资源（碎片时间）：
  - [Mastering Ethereum](https://github.com/ethereumbook/ethereumbook) — 第 4-5 章（密码学与钱包）
  - [iancoleman.io/bip39](https://iancoleman.io/bip39/) — 手动试验 HD 钱包推导
  - [Bitcoin Developer Guide](https://developer.bitcoin.org/devguide/) — 理解钱包设计思想的起源
  - [Devcon 演讲存档](https://archive.devcon.org/) — 搜索 "MPC wallet" 观看 1 个演讲，了解前沿方向

#### Day 2（周二）— 交易签名
- [ ] 阅读 [EIP-1559](https://eips.ethereum.org/EIPS/eip-1559)（1h）
- [ ] 动手（4h）：
  - 用私钥签名 Legacy 交易（EIP-155）
  - 用私钥签名 EIP-1559 交易（DynamicFeeTx）
  - 理解 chainID 在签名中的作用（防跨链重放）
  - 用 `eth_sendRawTransaction` 在测试网发送交易
- [ ] 笔记：对比 EIP-155 和 EIP-1559 的字段差异（1h）

#### Day 3（周三）— 消息签名
- [ ] 阅读 [EIP-191](https://eips.ethereum.org/EIPS/eip-191) 和 [EIP-712](https://eips.ethereum.org/EIPS/eip-712)（1.5h）
- [ ] 动手（3.5h）：
  - 实现 `personal_sign` 签名验证
  - 实现 EIP-712 类型化数据签名
  - 从签名恢复地址（`ecrecover`）
- [ ] 笔记：`personal_sign` vs `eth_sign` vs EIP-712 的安全差异（1h）

#### Day 4（周四）— 密钥存储设计
- [ ] 阅读（1h）：[go-ethereum keystore](https://goethereumbook.org/keystore/) 源码
- [ ] 设计（1h）：自己的密钥存储方案
  - 本地开发：加密 keystore 文件（AES-128-CTR + scrypt）
  - 生产环境：集成 AWS KMS / GCP Cloud HSM
- [ ] 编码（3h）：
  - 实现 keystore 创建（生成随机私钥 + 密码加密存储）
  - 实现 keystore 解密（密码 → 私钥）
  - 单元测试：错误的密码无法解密
- [ ] 笔记：HSM/KMS 与软件 keystore 的安全等级对比（1h）
- [ ] 浏览参考资源（碎片时间）：
  - [AWS KMS 文档](https://docs.aws.amazon.com/kms/) — 了解云 HSM 基本概念
  - [GCP Cloud HSM](https://cloud.google.com/kms/docs/hsm) — 对比不同云厂商方案

#### Day 5（周五）— 钱包服务骨架
- [ ] 编码（5h）：
  - 定义钱包服务 API（gRPC 或 REST）：
    - `POST /wallets` → 创建钱包（生成助记词 + keystore）
    - `GET /wallets/{id}` → 查询钱包信息（地址、余额）
    - `POST /wallets/{id}/sign` → 签名交易（返回签名后的 rawTx）
  - 实现 wallet service 核心逻辑
- [ ] 测试（1h）：用 curl/grpcurl 测试 API

#### Day 6（周六）— 安全加固
- [ ] 编码（4h）：
  - API 鉴权（API Key / JWT）
  - 签名操作增加二次确认（OTP 或二次密码）
  - 密钥操作日志（audit log）
  - 敏感信息脱敏（日志中不打印私钥、助记词）
- [ ] 安全自测（2h）：
  - 尝试用错误密码解密
  - 尝试绕过鉴权
  - 检查日志是否包含私钥
- [ ] 延伸阅读（碎片时间）：
  - 了解 Smart Contract Wallet（如 Safe/Gnosis Safe）与 EOA 钱包的区别
  - 理解热/冷/温钱包的隔离架构和各自的安全假设

#### Day 7（周日）— 弹性
- [ ] 用测试网 ETH（Sepolia/Holesky faucet）测试完整签名+发送流程
- [ ] 回顾本周，确保理解：一个 HD 钱包从助记词到签名交易的全流程

---

### 第 7 周：充值系统

> **本周目标：** 实现完整的充值检测流水线（Kafka → Deposit Consumer → 确认 → 入账），支持 ETH 主币和 ERC-20 token，处理 reorg 场景下的充值安全。
>
> **本周产出：** Deposit Consumer（ETH + ERC-20 统一处理），12 确认机制，reorg 回滚逻辑，防重复到账（复合唯一约束 + 幂等检查 + 定时对账），充值 API（查询 / 统计 / Webhook 通知）。

#### Day 1（周一）— 充值检测设计
- [ ] 设计（2h）：
  - 充值检测流程：Listener → Kafka → Deposit Consumer → 确认 → 入账
  - 地址管理：如何为新用户生成充值地址（HD 钱包派生）
  - 多币种支持：ETH 主币 vs ERC-20 token
- [ ] 编码（3h）：
  - 创建 deposit 数据库表（id, address, tx_hash, amount, token, confirmations, status, created_at）
  - 实现 address 索引（快速查找某个地址的所有充值记录）
- [ ] 笔记：充值系统状态机 draft（1h）

#### Day 2（周二）— 充值检测实现
- [ ] 编码（5h）：
  - 实现 Deposit Consumer：
    - 从 Kafka 消费交易数据
    - 筛选目标地址的转入交易（ETH 直接转账 + ERC-20 Transfer 事件）
    - 写入 deposit 表（状态：pending）
  - 确认逻辑：初始 confirmations=0，每来一个新块 confirmations+1
  - 达到 12 个确认后，状态变为 confirmed，触发入账逻辑
- [ ] 测试（1h）：向测试地址转入测试币，验证充值检测

#### Day 3（周三）— Reorg 场景下的充值处理
- [ ] 编码（4h）：
  - 当 deposit 对应的区块被 reorg 时：
    - 未确认（pending）：删除或标记为无效
    - 已确认（confirmed）：回滚余额（记录 rollback 日志）
  - 实现 deposit 回滚逻辑
- [ ] 测试（2h）：
  - 模拟 reorg：创建 pending deposit → 手动触发 reorg → 验证 deposit 状态

#### Day 4（周四）— 防重复到账
- [ ] 编码（4h）：
  - 唯一约束：`(tx_hash, log_index, to_address, token_address)` 复合唯一
  - 幂等检查：处理前先查是否已存在
  - 定时对账任务：每 10 分钟对比链上余额与系统余额
- [ ] 测试（2h）：
  - 重放同一交易 → 验证不会重复入账

#### Day 5（周五）— 充值通知 + API
- [ ] 编码（4h）：
  - 充值到账通知（WebSocket 推送或 Webhook）
  - 充值 API：
    - `GET /deposits?address=0x...` → 查询充值记录
    - `GET /deposits/{id}` → 单条充值详情
  - 充值统计：`GET /deposits/stats` → 总充值金额、笔数
- [ ] 测试（2h）：集成测试充值全流程

#### Day 6（周六）— ETH 与 ERC-20 统一处理
- [ ] 编码（4h）：
  - ETH 转账：检查 `tx.to` 是否为目标地址，`tx.value` 为金额
  - ERC-20 Transfer：解析 event log，检查 `to` 参数是否为目标地址
  - 统一的 deposit 记录（币种标识区分）
- [ ] 测试（2h）：ETH 充值 + USDC 充值 + 同时多笔

#### Day 7（周日）— 弹性
- [ ] 整理充值系统代码，确保测试覆盖率 > 70%
- [ ] 用文字或流程图描述 "一笔 ETH 从链上到账户余额" 的完整路径

---

### 第 8 周：提现系统 + Nonce 管理

> **本周目标：** 实现 Nonce Manager（并发安全分配 + gap 检测 + replacement）和提现完整状态机（创建 → 签名 → 广播 → 确认）。
>
> **本周产出：** Nonce Manager 完整实现（并发测试通过），提现状态机 + 核心流程，Signing Service（接口化 + KMS 扩展点），充提全链路集成测试，提现安全措施（金额上限 / 频率限制 / 白名单），tag v0.5.0。

#### Day 1（周一）— Nonce Manager 设计
- [ ] 设计（2h）：
  - Nonce 管理核心问题：
    1. 获取链上最新 nonce（`PendingNonceAt` vs `NonceAt`）
    2. 多笔并发交易的 nonce 分配
    3. nonce gap 检测与处理
  - 设计方案：内存 nonce pool + 定期链上同步
- [ ] 编码（3h）：
  - 实现 `NonceManager` 结构体
  - 方法：`ReserveNonce(address) → nonce`
  - 方法：`ReleaseNonce(address, nonce)`（交易失败时释放）
  - 方法：`ConfirmNonce(address, nonce)`（交易上链后确认）
- [ ] 笔记：画出 nonce 状态机（reserved → submitted → confirmed / released）（1h）

#### Day 2（周二）— Nonce Manager 编码
- [ ] 编码（5h）：
  - Nonce 预留时加锁（`sync.Mutex`），防止同一地址并发冲突
  - Nonce gap 检测：如果 nonce 5 长时间 pending → 用更高 gas 的替换交易覆盖
  - Nonce 同步：定期从链上查询 `PendingNonceAt` 校准
- [ ] 单元测试（1h）：
  - 并发预留 nonce 的竞争测试
  - gap 检测逻辑测试

#### Day 3（周三）— 提现流程设计 + Signing Service
- [ ] 设计提现状态机（1h）：
  - `pending` → `signing` → `broadcasting` → `pending_confirmation` → `confirmed` / `failed`
- [ ] 编码 Signing Service（4h）：
  - 接口：`SignTransaction(walletID, txParams) → signedTx`
  - 从 keystore 加载私钥（需要密码验证）
  - 构建并签名 EIP-1559 交易
  - 返回签名后的 rawTx bytes
- [ ] 安全审计自检（1h）

#### Day 4（周四）— 提现核心逻辑
- [ ] 编码（5h）：
  - Withdraw Service 主流程：
    1. 校验余额（余额 ≥ 提现金额 + 预估 gas 费）
    2. 调用 NonceManager.ReserveNonce() 获取 nonce
    3. 调用 SigningService.SignTransaction() 签名
    4. 调用 `eth_sendRawTransaction` 广播
    5. 更新 withdraw 状态、记录 tx_hash
    6. 失败时调用 NonceManager.ReleaseNonce()
- [ ] 测试（1h）：在测试网发起小额提现

#### Day 5（周五）— 提现状态跟踪 + 重试
- [ ] 编码（4h）：
  - 提现状态跟踪：定时轮询 tx receipt，确认状态
  - 失败重试策略：
    - nonce 冲突 → 重新签名（新 nonce）
    - gas 不足 → 重新估算 gas 后重试
    - 网络错误 → 指数退避重发
  - 提现超时处理：超过 N 分钟未确认 → 告警
- [ ] 测试（2h）：模拟各种失败场景

#### Day 6（周六）— 集成测试 + 安全加固
- [ ] 端到端测试（3h）：
  - 创建钱包 → 充值（模拟） → 发起提现 → 链上确认 → 余额扣减
  - 并发提现：10 笔并发提现，验证 nonce 管理和余额一致性
- [ ] 安全加固（2h）：
  - 提现金额上限检查
  - 提现频率限制（rate limit）
  - 提现白名单地址
- [ ] git tag v0.2.0（1h）

#### Day 7（周日）— Phase 2 检查点

**自检清单：**

- [ ] 能否解释 HD 钱包的密钥派生路径（BIP-32/39/44）？
- [ ] 能否展示 EIP-1559 交易的字段并解释每个字段的含义？
- [ ] Deposit Consumer 是否正确处理 ETH 和 ERC-20 两种充值？
- [ ] 充值是否经过 12 确认才视为有效？
- [ ] 是否实现了 reorg 场景下的充值回滚？
- [ ] Nonce Manager 是否处理了并发、gap、替换场景？
- [ ] 提现系统是否有完整的状态机（pending → confirmed/failed）？
- [ ] Signing Service 是否保证了私钥不离开安全边界？
- [ ] 是否有集成测试覆盖充提全流程？

---

## Phase 3：工程化与面试（第 9-12 周）

> **阶段目标：** 将项目部署到 K8S，完善 CI/CD 和可观测性，系统设计复习与 Mock Interview，投递简历拿到面试。
>
> **阶段产出：**
>
> - K8S 部署文件 + Helm Chart
> - Prometheus + Grafana 监控面板（含告警规则）
> - GitHub Actions CI/CD pipeline
> - 测试覆盖率 ≥ 70%
> - 英文简历 + Cover Letter + 技术博客 1 篇
> - 投递 ≥ 10 个岗位
> - 完整 Portfolio（项目截图 / 架构图 / Demo 视频）
> - tag v1.0.0

### 第 9 周：Docker + K8S 部署

> **本周目标：** 将全部服务容器化，编写 K8S 部署文件，部署到 minikube/kind，配置 Prometheus + Grafana 监控和告警。
>
> **本周产出：** Dockerfile × 2（多阶段构建），docker-compose.yml（全栈一键启动），K8S manifests（Deployment / Service / ConfigMap / Secret / PVC），Helm Chart，Prometheus metrics + Grafana Dashboard + Alert rules。

#### Day 1（周一）— Docker 化
- [ ] 编码（5h）：
  - 为 listener 和 wallet-backend 分别编写 Dockerfile
  - 使用多阶段构建（编译阶段 + 运行阶段），减少镜像体积
  - 编写 docker-compose.yml：
    - PostgreSQL
    - Kafka (KRaft mode, 无需 ZooKeeper)
    - Listener
    - Wallet Backend
  - 本地 `docker-compose up` 跑通
- [ ] 笔记：Docker 最佳实践（非 root 用户、健康检查、信号处理）（1h）

#### Day 2（周二）— K8S 基础补充 + 部署
- [ ] 学习（2h）：如果 K8S 基础薄弱，先看 [TechWorld with Nana K8S Crash Course](https://www.youtube.com/@TechWorldwithNana)
  - Deployment、Service、ConfigMap、Secret、PVC
- [ ] 浏览资料：
  - [Kubernetes 官方文档](https://kubernetes.io/docs/) — 核心概念章节
  - [CNCF YouTube](https://www.youtube.com/@cncf) — 收藏频道，后续关注区块链节点运维相关演讲
- [ ] 安装 minikube 或 kind（1h）
- [ ] 编码（3h）：
  - 编写 listener-deployment.yaml
  - 编写 wallet-deployment.yaml
  - 编写 postgres-statefulset.yaml（使用 PVC 持久化）
  - 编写 kafka-deployment.yaml

#### Day 3（周三）— K8S 配置管理
- [ ] 编码（4h）：
  - 使用 ConfigMap 管理配置文件
  - 使用 Secret 管理私钥/密码/API Key
  - 编写 Service（ClusterIP + LoadBalancer）
  - 配置 Readiness Probe 和 Liveness Probe
- [ ] 部署到 minikube，验证所有 pod 正常运行（2h）

#### Day 4（周四）— Helm Chart
- [ ] 阅读 [Helm 官方文档](https://helm.sh/docs/) 快速入门（0.5h）
- [ ] 学习 Helm 基础（1h）
- [ ] 编码（4h）：
  - 为项目编写 Helm Chart
  - 模板化：values.yaml 定义可配置参数
  - 支持 `helm install` 一键部署
- [ ] 测试：`helm install` → 验证 → `helm uninstall` 清理（1h）

#### Day 5（周五）— 可观测性：Prometheus + 系统调优意识
- [ ] 阅读 [Prometheus 官方文档](https://prometheus.io/docs/) 核心概念（0.5h）
- [ ] 编码（4h）：
  - 在 listener 中暴露 `/metrics` endpoint
  - 实现关键指标：
    - `blocks_processed_total`
    - `reorgs_detected_total`
    - `rpc_call_duration_seconds`
    - `deposit_detected_total`
    - `withdraw_processed_total`
    - `nonce_gap_detected_total`
  - 部署 Prometheus + 配置 ServiceMonitor
- [ ] 测试 Prometheus 数据采集（2h）
- [ ] 碎片时间了解（后续深入）：Linux 系统调优基础（fd 限制 / tcp keepalive / io scheduler / 内存分配），运维区块链节点的必备知识

#### Day 6（周六）— 可观测性：Grafana
- [ ] 编码/配置（4h）：
  - 部署 Grafana
  - 创建 Dashboard：
    - Block Processing Overview（速率、延迟、reorg 次数）
    - RPC Health（延迟、错误率）
    - Wallet Operations（充值/提现速率）
  - 配置 Alert Rule：
    - 区块处理延迟 > 2 分钟 → 告警
    - reorg 检测到 → 告警
- [ ] 截图保存到 docs/ 目录（2h）

#### Day 7（周日）— 弹性
- [ ] 用 k9s 或 Lens 观察集群运行状态
- [ ] 手动 kill pod，观察 K8S 自动恢复 + Grafana 告警

---

### 第 10 周：CI/CD + 安全审计 + 文档

> **本周目标：** 搭建 CI/CD pipeline，补齐测试覆盖率，完成安全自查，撰写完整项目文档和英文技术博客。
>
> **本周产出：** GitHub Actions CI（lint/test/build）+ CD（deploy），测试覆盖率 ≥ 70%（核心模块 ≥ 80%），安全审计报告（SQL 注入 / 私钥泄露 / 竞态条件 / 依赖漏洞），完整文档（ARCHITECTURE.md / DEPLOYMENT.md / MONITORING.md），英文技术博客 1 篇。

#### Day 1（周一）— GitHub Actions CI
- [ ] 编码（5h）：
  - 编写 `.github/workflows/ci.yml`：
    - 代码格式检查（gofmt、goimports）
    - Lint（golangci-lint）
    - 单元测试 + 覆盖率报告
    - 构建 Docker 镜像
  - 编写 `.github/workflows/deploy.yml`
- [ ] 提交代码，观察 CI 流水线运行（1h）

#### Day 2（周二）— 测试覆盖率提升
- [ ] 补充单元测试（4h）：
  - NonceManager 测试（并发、gap、替换）
  - Deposit dedup 测试
  - Withdraw state machine 测试
- [ ] 目标是测试覆盖率 ≥ 70%（2h）

#### Day 3（周三）— 安全审计自查
- [ ] 阅读参考资料（0.5h）：
  - [SWC Registry](https://github.com/SmartContractSecurity/SWC-registry) — 浏览 SWC-114/115/134 等与后端相关的条目
  - [OpenZeppelin 文档](https://docs.openzeppelin.com/) — 了解智能合约安全最佳实践
- [ ] 代码安全审查（4h）：
  - SQL 注入检查（所有 DB 查询是否使用参数化）
  - 私钥泄露检查（日志/错误消息/api 响应）
  - 竞态条件检查（`go test -race`）
  - 输入校验（金额范围、地址格式）
  - 依赖漏洞扫描（`govulncheck`）
- [ ] 记录发现的问题并修复（1.5h）

#### Day 4（周四）— 文档完善
- [ ] README（3h）：
  - 项目简介 + 核心功能
  - 架构图（Mermaid 或 Excalidraw 导出）
  - 快速开始（docker-compose 一键启动）
  - API 文档
  - 配置说明
- [ ] docs/ 目录（3h）：
  - ARCHITECTURE.md（详细架构设计 + 技术选型理由）
  - DEPLOYMENT.md（K8S + Helm 部署指南）
  - MONITORING.md（Prometheus 指标说明 + Grafana 截图）

#### Day 5（周五）— 技术博客（英文）
- [ ] 撰写一篇英文技术博客（5h）：
  - 主题建议 3 选 1：
    1. "Building a Production-Grade Ethereum Block Listener in Go"
    2. "Handling Blockchain Reorgs: A Practical Guide"
    3. "Designing a Nonce Manager for High-Throughput Ethereum Wallets"
  - 篇幅：1500-2000 词
  - 包含代码片段和架构图
- [ ] 发布到 Medium / Dev.to / 个人博客（1h）

#### Day 6（周六）— 系统设计复习
- [ ] 复习系统设计题（5h）：
  - 多链充值系统设计（5 条链以上，如何处理不同链的差异）
  - 高可用 RPC 集群设计（负载均衡、failover、限流）
  - 钱包系统安全架构（HSM 集成、多重签名、审计日志）
- [ ] 每个题目画架构图 + 列出关键技术点（1h）

#### Day 7（周日）— 弹性
- [ ] 用 Excalidraw/Draw.io 画系统设计题的架构图
- [ ] 对着镜子用英文讲一遍

---

### 第 11 周：面试冲刺（一）— 系统设计 + 编码

> **本周目标：** 掌握 Web3 面试高频系统设计题（多链充值 / RPC 集群 / 钱包架构），Go 编码题熟练，英语技术表达流畅。
>
> **本周产出：** 4 道系统设计白板练习（含架构图 + 英文录音），Go 编码练习（LRU Cache / Rate Limiter / Nonce Tracker），英文自我介绍（30s + 2min），2 次 Mock Interview 复盘记录。

#### Day 1（周一）— 系统设计 Mock：多链充值系统
- [ ] 自问自答（全程英文，录音）（5h）：
  1. 需求澄清：支持多少条链？充值确认数要求？多币种？
  2. 架构设计：每条链独立 listener？还是统一 listener？
  3. 数据模型：如何统一不同链的充值记录？
  4. 安全：reorg 处理、双花攻击防护
  5. 扩展性：当链从 5 条扩展到 20 条，系统如何扩展？
- [ ] 回听录音，找出表达不清晰的地方（1h）

#### Day 2（周二）— 系统设计 Mock：高可用 RPC 集群
- [ ] 自问自答（全程英文，录音）（5h）：
  1. 需求：99.99% 可用性，延迟 < 500ms
  2. 架构：多 provider 负载均衡 + 自建节点 + 健康检查
  3. 缓存策略：哪些数据可以缓存？（block、receipt、balance）
  4. 限流：防止 API Key 超限
  5. 降级：某 provider 宕机时如何无感切换？
- [ ] 回听录音，找出表达不清晰的地方（1h）

#### Day 3（周三）— 编码 Mock：LeetCode + Go 专项
- [ ] 刷题（4h）：
  - Go 并发题：实现一个并发安全的 LRU Cache
  - Go 并发题：实现一个 rate limiter（token bucket + sliding window）
  - Go 并发题：实现一个超时控制的重试机制
- [ ] 代码审查自查：错误处理、资源释放、并发安全（2h）

#### Day 4（周四）— 编码 Mock：区块链专项
- [ ] 编程题（4h）：
  - 实现一个简易的 Ethereum 交易解码器（RLP 解码）
  - 实现一个简易的 Merkle Tree 验证
  - 实现一个 nonce tracker（支持并发 safe 的 nonce 分配）
- [ ] 代码审查自查（2h）

#### Day 5（周五）— 行为面试 + 英语
- [ ] 准备行为面试回答（3h）：
  - "Tell me about yourself"（30 秒 + 2 分钟两个版本）
  - "Why Web3?"（表达从传统后端到 Web3 的动机）
  - "Describe a challenging technical problem you solved"
  - "How do you handle disagreements in a team?"
- [ ] 英语技术术语发音纠正（2h）：
  - 重点：reorg（/riˈɔːrɡ/）、nonce（/nɒns/）、Merkle（/ˈmɜːrkəl/）、Ethereum（/ɪˈθɪəriəm/）
- [ ] 录音并回听（1h）

#### Day 6（周六）— 综合 Mock Interview
- [ ] 找一个朋友或在 Discord/Telegram 社区找一个 Web3 工程师（2h）
- [ ] 如果没有真人，用 ChatGPT Voice Mode 模拟面试（3h）：
  - 自我介绍
  - 系统设计：设计一个去中心化交易所的充值系统
  - 技术深度：解释 EVM 的 storage 和 memory 区别
- [ ] 复盘表现，记录待改进点（1h）

#### Day 7（周日）— 弹性
- [ ] 整理所有 Mock Interview 中的薄弱点
- [ ] 针对性补强（重看对应章节的笔记）

---

### 第 12 周：面试冲刺（二）— 简历 + 投递

> **本周目标：** 简历定稿、正式投递目标岗位、完成最终项目展示（Demo 视频）。
>
> **本周产出：** 英文简历 + Cover Letter 定稿，≥ 10 份岗位投递，Demo 视频录制，LinkedIn 优化 + 社区加入，全部学习笔记归档，tag v1.0.0。

#### Day 1（周一）— 简历重构
- [ ] 简历改写（4h）：
  - 突出 Web3 项目经验（listener + wallet）
  - 量化成果（"处理过 X 万笔交易"、"支持 Y 条链"）
  - 关键词优化（EVM、reorg、nonce management、HD wallet、Kafka、K8S）
  - 链接 GitHub 项目和技术博客
- [ ] 英文简历版本（2h）

#### Day 2（周二）— LinkedIn + 社区
- [ ] LinkedIn 优化（2h）：
  - 头像 + banner
  - About 部分（英文，突出 Web3 Infra 方向）
  - 发布一篇技术文章
- [ ] 加入社区（3h）：
  - Discord：ETHGlobal、Ethereum R&D
  - Telegram：各 Web3 招聘群
  - 关注 Twitter/X 上的 Web3 infra 项目和工程师
- [ ] 提交 2-3 个开源项目的 PR（typofix 或 doc fix 也可）（1h）

#### Day 3（周三）— 投递简历
- [ ] 整理目标公司清单（2h）：
  - Tier 1：钱包/托管（Fireblocks、Copper、Cobo、Safeheron）
  - Tier 2：Infra/RPC（Alchemy、QuickNode、Chainstack、Blockdaemon）
  - Tier 3：支付（Circle、MoonPay、Ramp）
- [ ] 投递简历（4h）：
  - LinkedIn EasyApply + 官网投递
  - Web3 招聘平台：cryptojobslist.com、web3.career
  - 主动 DM 目标公司的 HR/工程师

#### Day 4（周四）— 面试跟随练习
- [ ] 针对心仪公司的 JD 做定向准备（3h）
- [ ] 系统设计终极复习（3h）：
  - reorg 处理、nonce 管理、签名服务、多链架构
  - 每个题目能 5 分钟英文流利讲清楚

#### Day 5（周五）— 最终项目展示准备
- [ ] 准备 Demo 视频或 Live Demo（4h）：
  - localhost 启动全套服务
  - 演示：监听新区块 → Kafka → 消费 → 检测充值 → 发起提现 → 链上确认
  - 展示 Grafana Dashboard
- [ ] 录制 Demo 视频（2h）

#### Day 6（周六）— 全局回顾
- [ ] 回顾 12 周学习笔记（3h）
- [ ] 整理 "高频面试题" 回答（3h）
- [ ] git tag v1.0.0

#### Day 7（周日）— 最终检查点 + 下一步方向

> **3 个月后的持续学习方向：**
> - **Rust 入门**：Solana / Sui / Polkadot 生态的节点和 SDK 均为 Rust，是 Web3 Infra 的第二语言
> - **非 EVM 链深入**：Solana（账户模型 / 无 nonce / getProgramAccounts）、Sui（对象模型 / 无区块概念）、Bitcoin（UTXO / Taproot）
> - **MPC 钱包**：深入学习安全多方计算在钱包中的应用（ZenGo、Fireblocks 等产品的技术方案）
> - **链数据平台**：The Graph / Dune Analytics / 自建索引器的架构设计


**最终检查清单：**

- [ ] GitHub 仓库是否包含完整的 README + 架构图 + 部署文档？
- [ ] 是否有技术博客（Medium/Dev.to）？
- [ ] 是否能用英文在 5 分钟内讲清楚自己的项目？
- [ ] 是否能流利回答 reorg / nonce / 签名服务 / 多链架构 这 4 个核心系统设计题？
- [ ] 是否有 Grafana Dashboard 截图？
- [ ] 是否在 3 个以上渠道投递了简历？
- [ ] 是否加入了 Discord/Telegram 社区并参与了讨论？

---

## 附录 A：每周时间分配参考

| 学习模块 | 每周小时数 | 占比 |
|----------|-----------|------|
| 阅读/视频学习 | 10h | ~28% |
| 项目编码 | 20h | ~56% |
| 笔记整理 | 3h | ~8% |
| 测试 + 调试 | 3h | ~8% |
| **合计** | **36h** | **100%** |

## 附录 B：关键 Milestone

| 时间点 | Milestone | 标志 |
|--------|-----------|------|
| 第 4 周末 | Phase 1 检查点 | Listener 稳定运行 2 小时 |
| 第 8 周末 | Phase 2 检查点 | 充提全流程可演示 |
| 第 10 周末 | 技术博客发布 | 博客链接 |
| 第 12 周末 | 最终检查点 | 简历投递 ≥ 10 家公司 |

## 附录 C：工具链速查

| 用途 | 工具 | 优先级 |
|------|------|--------|
| IDE | VS Code + Go 插件 / GoLand | 必须 |
| 链 RPC | Infura / Alchemy 免费 tier | 必须 |
| 数据库 | PostgreSQL 15+ | 必须 |
| 消息队列 | Kafka (KRaft mode) | 必须 |
| 容器 | Docker + docker-compose | 必须 |
| K8S | minikube / kind | 必须 |
| 监控 | Prometheus + Grafana | 必须 |
| CI/CD | GitHub Actions | 推荐 |
| 测试网水龙头 | Sepolia / Holesky faucet | 必须 |
| API 测试 | curl / grpcurl / Postman | 推荐 |
| 版本管理 | Git + GitHub | 必须 |
