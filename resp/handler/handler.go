package handler

import (
	"context"
	"go-redis/database"
	dataface "go-redis/interface/database"
	"go-redis/lib/logger"
	"go-redis/lib/sync/atomic"
	"go-redis/resp/connection"
	"go-redis/resp/parser"
	"go-redis/resp/reply"
	"io"
	"net"
	"strings"
	"sync"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/3/4
 * Time: 16:42
 * Description: No Description
 */

var (
	unknownErrReply = []byte("-ERR unknown\r\n")
)

type RespHandler struct {
	activeConn sync.Map
	db         dataface.Database
	closing    atomic.Boolean
}

func MakeHandler() *RespHandler {
	// TODO: 实现Database接口
	var db dataface.Database
	db = database.NewEchoDatabase()
	return &RespHandler{
		db: db,
	}
}

// Handle 处理Tcp连接
func (r *RespHandler) Handle(ctx context.Context, conn net.Conn) {
	if r.closing.Get() {
		_ = conn.Close()
	}
	client := connection.NewConnection(conn)
	r.activeConn.Store(client, struct{}{})
	ch := parser.ParseStream(conn)
	for payload := range ch {
		// Error
		if payload.Err != nil {
			if payload.Err == io.EOF ||
				payload.Err == io.ErrUnexpectedEOF ||
				strings.Contains(payload.Err.Error(), "use of closed network connection") {
				r.closeClient(client)
				logger.Info("connection closed: " + client.RemoteAddr().String())
				return
			}
			// protocol error
			errReply := reply.MakeErrReply(payload.Err.Error())
			if err := client.Write(errReply.ToBytes()); err != nil {
				r.closeClient(client)
				logger.Info("connection closed: " + client.RemoteAddr().String())
				return
			}
			continue
		}
		// Exec
		if payload.Data == nil {
			logger.Info("payload.Data is nil")
			continue
		}

		reply, ok := payload.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Info("require multi bulk reply")
			continue
		}

		result := r.db.Exec(client, reply.Args)
		if result == nil {
			_ = client.Write(unknownErrReply)
			continue
		}
		_ = client.Write(result.ToBytes())
	}
}

// 关闭单个client
func (r *RespHandler) closeClient(client *connection.Connection) {
	_ = client.Close()
	r.db.AfterClientClose(client)
	r.activeConn.Delete(client)
}

// Close 关闭协议
func (r *RespHandler) Close() error {
	logger.Info("handler shutting down.....")
	r.closing.Set(true)
	r.activeConn.Range(
		func(key, value interface{}) bool {
			client := key.(*connection.Connection)
			_ = client.Close()
			return true
		},
	)
	r.db.Close()

	return nil
}
