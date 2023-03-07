package cluster

import (
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/3/7
 * Time: 20:58
 * Description: No Description
 */

func flushdb(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	replies := cluster.broadcast(c, cmdArgs)
	var errReply reply.ErrorReply
	for _, r := range replies {
		if reply.IsErrorReply(r) {
			errReply = r.(reply.ErrorReply)
			break
		}
	}
	if errReply != nil {
		return reply.MakeErrReply("error: " + errReply.Error())
	}

	return reply.MakeOkReply()
}
