package tcp

import (
	"context"
	"go-redis/interface/tcp"
	"go-redis/lib/logger"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/2/28
 * Time: 15:14
 * Description: No Description
 */

type Config struct {
	Address string
}

// ListenAndServeWithSignal 接收系统信号
func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	closeChan := make(chan struct{})
	signChan := make(chan os.Signal)
	signal.Notify(signChan, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	go func() {
		sig := <-signChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()

	listen, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}
	logger.Info("start listen")
	ListenAndServe(listen, handler, closeChan)
	return nil
}

func ListenAndServe(listen net.Listener, handler tcp.Handler, closeChan <-chan struct{}) {

	// 服务不是正常退出，而是被杀掉时，走这个方法关闭资源
	go func() {
		<-closeChan
		logger.Info("shutting done...")
		_ = listen.Close()
		_ = handler.Close()
	}()

	// 服务正常退出
	defer func() {
		_ = listen.Close()
		_ = handler.Close()
	}()

	ctx := context.Background()
	waitDone := sync.WaitGroup{}
	for {
		conn, err := listen.Accept()
		if err != nil {
			break
		}
		logger.Info("accepted link")
		waitDone.Add(1)
		go func() {
			defer func() {
				waitDone.Done()
			}()
			handler.Handle(ctx, conn)
		}()
	}
	waitDone.Wait()
}
