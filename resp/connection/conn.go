package connection

import (
	"go-redis/lib/sync/wait"
	"net"
	"sync"
	"time"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/3/2
 * Time: 13:13
 * Description: No Description
 */

type Connection struct {
	conn         net.Conn
	waitingReply wait.Wait
	mu           sync.Mutex
	selectDB     int
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		conn: conn,
	}
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Connection) Close() error {
	c.waitingReply.WaitWithTimeout(10 * time.Second)
	_ = c.conn.Close()
	return nil
}

func (c *Connection) Write(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}
	c.mu.Lock()
	c.waitingReply.Add(1)
	defer func() {
		c.mu.Unlock()
		c.waitingReply.Done()
	}()
	_, err := c.conn.Write(bytes)

	return err
}

func (c *Connection) GetDBIndex() int {
	return c.selectDB
}

func (c *Connection) SelectDB(dbNum int) {
	c.selectDB = dbNum
}
