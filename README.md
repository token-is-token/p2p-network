# P2P Network

P2P网络 | Peer-to-Peer Network Layer

## 项目介绍

P2P Network 是 LLM Share Network 的去中心化通信层，负责节点发现、消息路由、请求分发和网络维护。

## 架构

```
p2p-network/
├── cmd/node/          # 应用程序入口
├── pkg/
│   ├── node/         # P2P节点核心
│   ├── dht/          # 分布式哈希表
│   ├── pubsub/       # 发布/订阅
│   ├── discovery/    # 节点发现
│   ├── relay/       # 中继服务
│   ├── protocol/    # 自定义协议
│   └── utils/       # 工具函数
└── test/            # 测试文件
```

## 技术栈

- **语言**: Go 1.21+
- **P2P**: libp2p-go
- **DHT**: go-libp2p-kad-dht
- **PubSub**: go-libp2p-pubsub

## 安装

```bash
# 克隆仓库
git clone https://github.com/your-org/p2p-network.git
cd p2p-network

# 下载依赖
go mod download

# 构建
make build

# 运行测试
make test
```

## 快速开始

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/your-org/p2p-network/pkg/node"
)

func main() {
    // 创建配置
    cfg := node.DefaultConfig()
    
    // 创建节点
    n, err := node.NewNode(cfg)
    if err != nil {
        panic(err)
    }
    
    // 启动节点
    if err := n.Start(context.Background()); err != nil {
        panic(err)
    }
    
    fmt.Printf("Node started with ID: %s\n", n.ID())
    
    // 阻塞直到关闭
    <-n.Context().Done()
}
```

## 模块说明

### Node (节点核心)
P2P 节点核心模块，负责管理 libp2p 主机、DHT 和 PubSub。

### DHT (分布式哈希表)
提供分布式键值存储，用于节点发现和内容路由。

### PubSub (发布/订阅)
支持主题订阅和消息发布，用于节点间通信。

### Discovery (发现服务)
提供多种发现方式：mDNS、引导节点等。

### Relay (中继服务)
支持 NAT 穿透和连接中继。

### Protocol (自定义协议)
实现应用层自定义协议。

## 开发

```bash
# 运行单元测试
make test

# 运行集成测试
make test-integration

# 代码检查
make lint

# 性能测试
make benchmark
```

## 许可证

MIT License - see LICENSE file for details
