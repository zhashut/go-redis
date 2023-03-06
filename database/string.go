package database

import (
	"go-redis/interface/database"
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/3/6
 * Time: 14:36
 * Description: No Description
 */

// Get k1
func execGet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.MakeNullBulkReply()
	}
	return reply.MakeBulkReply(entity.Data.([]byte))
}

// SET K1 V
func execSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	value := args[1]
	entity := &database.DataEntity{
		Data: value,
	}
	db.PutEntity(key, entity)

	return reply.MakeOkReply()
}

// SETNX k1 v, 如果存在则不插入，否则插入
func execSetNX(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	value := args[1]
	entity := &database.DataEntity{
		Data: value,
	}
	result := db.PutIfAbsent(key, entity)

	return reply.MakeIntReply(int64(result))
}

// GETSET k1 v1 v2 (这个命令是获取key对应的value值并修改成新的value，但是返回的还是旧的value)
func execGetSet(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	value := args[1]
	entity, exists := db.GetEntity(key)
	db.PutEntity(key, &database.DataEntity{
		Data: value,
	})
	if !exists {
		return reply.MakeNullBulkReply()
	}
	return reply.MakeBulkReply(entity.Data.([]byte))
}

// STRLEN k -> value
func execStrLen(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.MakeNullBulkReply()
	}
	bytes := entity.Data.([]byte)

	return reply.MakeIntReply(int64(len(bytes)))
}

func init() {
	RegisterCommand("Get", execGet, 2)  // get k1
	RegisterCommand("Set", execSet, -3) //set k1 v1
	RegisterCommand("SetNx", execSetNX, 3)
	RegisterCommand("GetSet", execGetSet, 3)
	RegisterCommand("StrLen", execStrLen, 2)
}
