package cluster

import (
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/3/7
 * Time: 20:45
 * Description: No Description
 */

// Rename k1 k2
func rename(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	if len(cmdArgs) != 3 {
		return reply.MakeErrReply("ERR Wrong number args")
	}
	src := string(cmdArgs[1])
	dest := string(cmdArgs[2])

	srcPeer := cluster.peerPicker.PickNode(src)
	destPeer := cluster.peerPicker.PickNode(dest)

	// 判断是否在同一节点上
	if srcPeer != destPeer {
		return reply.MakeErrReply("ERR rename must within on peer")
	}

	return cluster.relay(srcPeer, c, cmdArgs)
}
