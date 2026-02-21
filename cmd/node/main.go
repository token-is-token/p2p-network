package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/your-org/p2p-network/pkg/node"
	"github.com/your-org/p2p-network/pkg/utils"
)

func main() {
	// 创建上下文，用于优雅关闭
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建日志
	logger, err := utils.NewLogger("p2p-node", utils.LogLevelInfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	logger.Info("Starting P2P Node...")

	// 创建默认配置
	cfg := node.DefaultConfig()

	// 可选的命令行参数覆盖
	if len(os.Args) > 1 {
		for i := 1; i < len(os.Args); i++ {
			switch os.Args[i] {
			case "--port":
				if i+1 < len(os.Args) {
					cfg.ListenPort = os.Args[i+1]
					i++
				}
			case "--bootnodes":
				if i+1 < len(os.Args) {
					cfg.BootstrapPeers = parsePeers(os.Args[i+1])
					i++
				}
			case "--enable-relay":
				cfg.EnableRelay = true
			case "--verbose", "-v":
				logger.SetLevel(utils.LogLevelDebug)
			}
		}
	}

	// 创建节点
	n, err := node.NewNode(cfg, node.WithLogger(logger))
	if err != nil {
		logger.Error("Failed to create node", "error", err)
		os.Exit(1)
	}

	// 启动节点
	if err := n.Start(ctx); err != nil {
		logger.Error("Failed to start node", "error", err)
		os.Exit(1)
	}

	logger.Info("Node started successfully",
		"peerID", n.ID(),
		"listenAddrs", n.Addrs(),
	)

	// 创建指标收集器
	metrics, err := utils.NewMetrics("p2p-node")
	if err != nil {
		logger.Warn("Failed to create metrics", "error", err)
	} else {
		logger.Info("Metrics enabled", "port", "9090")
		go metrics.Start(ctx, 9090)
	}

	// 等待信号以优雅关闭
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		logger.Info("Received signal, shutting down", "signal", sig)
	case <-ctx.Done():
	}

	// 停止节点
	if err := n.Stop(ctx); err != nil {
		logger.Error("Error stopping node", "error", err)
	}

	logger.Info("Node stopped")
}

// parsePeers 解析引导节点地址
func parsePeers(peersStr string) []string {
	if peersStr == "" {
		return nil
	}
	// 简单解析逗号分隔的地址
	// 在实际使用中，应该使用 multiaddr 解析
	return nil
}
