package tcp

import (
	"context"
	"net"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/2/28
 * Time: 15:07
 * Description: No Description
 */

type Handler interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}
