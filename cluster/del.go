package cluster

import (
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/3/7
 * Time: 21:04
 * Description: No Description
 */

func del(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	replies := cluster.broadcast(c, cmdArgs)
	var errReply reply.ErrorReply
	var deleted int64
	for _, r := range replies {
		if reply.IsErrorReply(r) {
			errReply = r.(reply.ErrorReply)
			break
		}
		intReply, ok := r.(*reply.IntReply)
		if !ok {
			errReply = reply.MakeErrReply("error intReply")
			continue
		}
		deleted += intReply.Code
	}
	if errReply != nil {
		return reply.MakeErrReply("error: " + errReply.Error())
	}

	return reply.MakeIntReply(deleted)
}
