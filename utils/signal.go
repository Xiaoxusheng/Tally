package utils

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func Listen() {

	signalCh := make(chan os.Signal, 1)
	// 监听指定的系统信号
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// 启动一个 goroutine 来处理接收到的系统信号
	// 接收信号
	for {
		select {
		case sig := <-signalCh:
			fmt.Printf("Received signal: %s\n", sig)
			os.Exit(0)
		}

	}
	// 打印接收到的信号
	// 进行一些清理工作或其他操作
	// ...

	// 退出程序
	//os.Exit(0)
}
