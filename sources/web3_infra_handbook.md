# Web3 Infra 转型规划与行动指南（针对 10 年 Java 后端工程师）

[TOC]

---

## 个人背景与目标分析

- 年龄：35+
- 英文阅读流畅，口语一般（⚠️ 远程岗位日常需要英语 standup/design review/incident response，口语是需要重点突破的瓶颈）
- Java 技术：SpringCloud + MQ
- K8S：掌握基本命令
- Netty：未掌握
- 无金融支付和风控经验
- 可接受远程和夜间沟通
- 收入期望：高且尽量稳定
- 可转方向：基础设施、钱包支付，全栈可考虑
- 当前技能：Java、Go + Gin

**核心目标：**  
成为市场认可的 Web3 基础设施工程师，优先聚焦钱包/支付系统和链监听系统。

---

## 最佳定位与市场选择

### 推荐岗位定位
- Senior Backend / Infra Engineer with Web3 experience
- Web3 Infrastructure Backend Engineer

> 后端基础设施岗位与 Solidity/DApp 岗位竞争同样激烈，但竞争维度不同：后端更看重分布式系统、安全、链底层理解；Solidity/DApp 更看重 EVM 深度、Gas 优化、前端生态。选择应基于自身技术积累优势，而非单纯回避竞争。

### 目标公司类型
1. 钱包/托管公司（Wallet, Custody, MPC）
2. 支付公司（Stablecoin, Cross-border）
3. Infra 公司（RPC, Node Service, Indexer, Chain Analytics）

### 谨慎或不推荐
- 交易所：高压、高卷，金融背景缺失
- 发币小团队：风险高，token 激励为主，不稳定

---

## 核心技术路线

### 1. Go 深化
- 并发：goroutine、channel、worker pool、context cancel、sync.Mutex/sync.RWMutex、atomic
  - > 注意：sync.Map 仅适用于 key 写一次读多次的特定场景，大多数情况用 map + sync.RWMutex
- 网络：websocket、grpc、protobuf
- 性能优化：pprof、GC、memory leak、object pool

**学习资源**  
- 官方文档：https://go.dev/doc/  
- Go 并发模式：https://go.dev/blog/pipelines  
- 代码风格参考：[Uber Go Style Guide](https://github.com/uber-go/guide)  
- 视频：[JustForFunc](https://youtube.com/@JustForFunc)（Francesc Campoy，前 Go 团队成员，并发模型讲解极佳）  
- 视频：[Ardan Labs](https://youtube.com/@ardanlabs)（Bill Kennedy，生产级 Go 训练）

### 2. 链基础设施理解（EVM 链为起点）

> 本章及后续 Listener/钱包设计均以 EVM 兼容链（Ethereum/L2）为默认上下文。Solana、Sui 等非 EVM 链的数据模型差异巨大（无 nonce、无 eth_getLogs、不同账户模型），不能直接套用，需要在后续学习中单独深入。

- RPC 调用：eth_call、eth_getLogs、pending tx、websocket subscription
- 交易生命周期：nonce、gas、pending、reorg、retry
- Event 系统：topics、indexed、ABI decoding

**学习资源**  
- 官方文档：[Ethereum JSON-RPC](https://ethereum.org/en/developers/docs/apis/json-rpc/)  
- 开源项目：[go-ethereum](https://github.com/ethereum/go-ethereum)  
- 视频：[Cyfrin Updraft](https://updraft.cyfrin.io/)（Patrick Collins，全球最全面的免费区块链课程，含 Solidity/Foundry/DeFi 安全）  
- 视频：[ETHGlobal](https://youtube.com/@ETHGlobal)（官方黑客松 Workshop 和演讲）  
- 视频：[Finematics](https://youtube.com/@Finematics)（DeFi 概念可视化讲解）

### 3. 钱包系统核心
- 钱包模型：EOA、Smart Contract Wallet、HD Wallet、MPC Wallet
- 签名机制：secp256k1、transaction signing、message signing（EIP-191 / EIP-712）
- 安全：热/冷钱包隔离、HSM/KMS 密钥管理、提现审核、防 replay/nonce attack

**学习资源**  
- 官方文档：[Bitcoin Developer Guide](https://developer.bitcoin.org/devguide/)  
- Mastering Ethereum（免费在线阅读）：https://github.com/ethereumbook/ethereumbook  
- HD 钱包基础：https://iancoleman.io/bip39/  
- 视频：[Devcon 演讲存档](https://archive.devcon.org/)（以太坊基金会年度开发者大会，MPC/ZKP/安全等前沿话题）

### 4. Infra & 运维
- K8S：StatefulSet、PVC、Helm、Prometheus、Grafana
  - > 如果仅掌握基本命令，建议先补 Deployment/Service/ConfigMap/Secrets 基础，再深入 StatefulSet 和 Helm
- Linux：fd、tcp、io、memory、network tuning
- 监控与可观测性：block lag、pending tx、rpc latency

**学习资源**  
- 官方文档：[Kubernetes](https://kubernetes.io/docs/)  
- Helm 官方文档：https://helm.sh/docs/  
- Prometheus 官方：https://prometheus.io/docs/  
- 视频：[TechWorld with Nana](https://www.youtube.com/@TechWorldwithNana)（最受欢迎的 K8s 教程，适合入门到 CKA 备考）  
- 视频：[CNCF](https://www.youtube.com/@cncf)（KubeCon 官方演讲，了解云原生生态前沿）

---

## 链监听系统（Listener）设计

### 架构
```text
Blockchain Node
       ↓
Websocket Listener
       ↓
Block Parser
       ↓
Kafka
       ↓
Consumer Group
       ↓
Database
```

### 难点与解决方案
1. **reorg（区块回滚）**  
   - 区分 "已确认" 与 "未确认" 数据：未确认数据以 block_hash 为 key 存储，确认后才按 block_number 索引
   - 检测 reorg：当某高度 block_hash 与已存储的不同时触发回滚
   - 处理流程：标记旧链数据（或设置 `is_canonical=false`）→ 重新解析新区块 → 幂等写入新数据
   - 建议：充提相关操作仅在 >= N 个确认后（如 ETH 建议 12+ 确认）标记为最终
2. **Websocket 中断**  
   - 自动重连 + 指数退避
   - checkpoint last_processed_block，重启时从断点扫描
   - gap scan：用 HTTP RPC 拉取缺失高度区间，与 WebSocket 互补
3. **重复事件**  
   - EVM 链使用 tx_hash + log_index + block_hash 复合唯一约束（仅 tx_hash+log_index 无法应对 reorg 场景）
   - 或采用：未确认数据先入临时表，确认后迁移到主表再做去重
4. **高吞吐量**  
   - MQ 解耦（Kafka/Redis Streams）  
   - Consumer Group 并行处理（按 address/tx_hash 分片避免乱序）

**学习资源**  
- Kafka 官方文档：https://kafka.apache.org/documentation/  
- 参考实现：[go-ethereum](https://github.com/ethereum/go-ethereum)（阅读 ethclient 和 event subscription 部分）  
- 视频：[Confluent 官方频道](https://www.youtube.com/@Confluent)（Kafka 创始团队出品，含 Kafka 101/Kafka Streams 等完整系列）

---

## 钱包系统设计（生产级）

### 架构图
```text
                    +----------------+
                    | API Gateway    |
                    +----------------+
                             |
        -----------------------------------------
        |                  |                    |
+----------------+ +----------------+ +----------------+
| Wallet Service | | Deposit Service| | WithdrawSvc   |
+----------------+ +----------------+ +----------------+
        |                  |                    |
        |          +----------------+           |
        |          | Block Listener |           |
        |          +----------------+           |
        |                  |                    |
        |             +---------+               |
        |             | Kafka   |               |
        |             +---------+               |
        |                             +----------------+
        |                             | Nonce Manager  |
        |                             +----------------+
        |                                      |
        |                             +----------------+
        |                             | Signing Service|
        |                             +----------------+
        |                                      |
        ----------------------------------------
                             |
                     +----------------+
                     | Blockchain RPC |
                     +----------------+
```

> **架构补充说明：**
> - **Signing Service 必须与 HSM（硬件安全模块）或云 KMS 集成**。生产级钱包系统绝不能将私钥以明文存储在应用服务器内存中。
> - **Nonce Manager 核心挑战**：同一地址并发发交易时需 nonce 预留/锁定；nonce gap（如 nonce 5 卡住导致 6/N 全部阻塞）需要检测与重发机制；同一 nonce 的 replacement/cancellation 策略。
> - 区块链 RPC 节点建议至少 3 个不同提供商的 fallback（如 Alchemy + Infura + 自建节点），避免单点故障。

**学习资源**  
- 参考：[go-ethereum](https://github.com/ethereum/go-ethereum)（accounts/keystore 部分）  
- AWS KMS 文档：https://docs.aws.amazon.com/kms/  
- GCP Cloud HSM：https://cloud.google.com/kms/docs/hsm

---

## 高级面试题汇总

1. 多链充值系统设计  
2. 如何处理 reorg（深度追问：block_hash vs block_number、确认数选择依据）  
3. nonce 冲突处理（追问：nonce gap、replacement、并发控制）  
4. websocket 中断恢复（追问：断点续扫、与 HTTP polling 的互补）  
5. 如何避免重复到账（追问：reorg 场景下的幂等性）  
6. signing service 设计（追问：HSM 集成、密钥生命周期管理）  
7. 提现系统可靠性保障（追问：状态机设计、异常回滚）
8. 常见安全攻击面（frontrunning、sandwich attack、reentrancy、签名重放）

---

## GitHub 项目结构建议

```text
wallet-infra/
 ├── cmd/
 ├── internal/
 ├── api/
 ├── listener/
 ├── deposit/
 ├── withdraw/
 ├── nonce/
 ├── signer/
 ├── kafka/
 ├── storage/
 ├── docker/
 ├── deploy/
 ├── helm/
 ├── docs/
 └── README.md
```

---

## 行动计划（12 个月）

> 以下为建议时间线。实际速度取决于每天投入时间和是否有全职工作。6 个月仅为极高强度脱产学习的乐观估计；正常兼职学习建议按 12-18 个月规划。

| 月份 | 目标 |
|------|------|
| 1-2 | 链基础学习 + Go 深化（websocket, RPC, tx lifecycle） |
| 3-4 | 完成链监听系统（listener + Kafka + DB），重点攻克 reorg 处理 |
| 5-6 | 完成钱包系统（deposit + withdraw + nonce manager + signing） |
| 7-8 | K8S 部署链节点与服务，Prometheus + Grafana 监控 |
| 9-10 | 工程化（Docker + Helm + CI/CD）+ 安全审计意识培养 |
| 11-12 | 投递简历，mock interview，系统设计复习，英语口语强化 |

---

## 下一阶段推荐深入内容

1. MPC 钱包原理与架构  
2. 提现系统完整状态机  
3. 高可用 RPC 集群设计  
4. 多链系统架构（ETH/Solana/Sui）—— Solana/Sui 需学习 Rust 和完全不同的账户/交易模型  
5. Kafka 在 Web3 的真实玩法  
6. 链数据平台设计  
7. Web3 风控设计  
8. Web3 Infra DevOps 体系  
9. 海外 remote 面试完整套路（英语口语表达 + 系统设计白板）  
10. 足够拿 offer 的 GitHub 项目设计  
11. Rust 入门（Solana/Sui/Polkadot 生态必备，与 Go 互补）  
12. Web3 安全审计基础（常见漏洞：reentrancy、frontrunning、flash loan、签名重放）

---

## 学习资源（汇总）

1. **Go 语言**  
   - 官网：https://go.dev/doc/  
   - Go 并发模式：https://go.dev/blog/pipelines  
   - 代码风格参考：[Uber Go Style Guide](https://github.com/uber-go/guide)  
   - 视频：[JustForFunc](https://youtube.com/@JustForFunc)（Francesc Campoy） | [Ardan Labs](https://youtube.com/@ardanlabs)（Bill Kennedy）

2. **Ethereum / RPC**  
   - 官网：https://ethereum.org/en/developers/docs/apis/json-rpc/  
   - 开源：[go-ethereum](https://github.com/ethereum/go-ethereum)  
   - 视频：[Cyfrin Updraft](https://updraft.cyfrin.io/)（Patrick Collins，免费全套课程） | [ETHGlobal](https://youtube.com/@ETHGlobal) | [Finematics](https://youtube.com/@Finematics)（DeFi 可视化）

3. **钱包系统 / HD 钱包**  
   - Bitcoin Developer Guide：https://developer.bitcoin.org/devguide/  
   - Mastering Ethereum：https://github.com/ethereumbook/ethereumbook  
   - HD 钱包工具：https://iancoleman.io/bip39/  
   - 视频：[Devcon 演讲存档](https://archive.devcon.org/)（MPC/ZKP/安全前沿话题）

4. **Kubernetes / Infra**  
   - 官网：https://kubernetes.io/docs/  
   - Helm：https://helm.sh/docs/  
   - Prometheus：https://prometheus.io/docs/  
   - 视频：[TechWorld with Nana](https://www.youtube.com/@TechWorldwithNana) | [CNCF](https://www.youtube.com/@cncf)

5. **Kafka / MQ**  
   - 官网：https://kafka.apache.org/documentation/  
   - 视频：[Confluent 官方频道](https://www.youtube.com/@Confluent)（Kafka 创始团队出品）

6. **安全**  
   - SWC Registry（智能合约漏洞分类）：https://github.com/SmartContractSecurity/SWC-registry  
   - OpenZeppelin 文档：https://docs.openzeppelin.com/
