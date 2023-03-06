package database

import (
	"go-redis/datastruct/dict"
	"go-redis/interface/database"
	"go-redis/interface/resp"
	"go-redis/resp/reply"
	"strings"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/3/5
 * Time: 16:15
 * Description: No Description
 */

type DB struct {
	index  int
	data   dict.Dict
	addAof func(CmdLine)
}

type ExecFunc func(db *DB, args [][]byte) resp.Reply

type CmdLine = [][]byte

func makeDB() *DB {
	db := &DB{
		data:   dict.MakeSyncDict(),
		addAof: func(line CmdLine) {},
	}
	return db
}

func (db *DB) Exec(c resp.Connection, cmdLine CmdLine) resp.Reply {
	cmdName := strings.ToLower(string(cmdLine[0]))
	cmd, ok := cmdTable[cmdName]
	if !ok {
		return reply.MakeErrReply("ERR unknown command " + cmdName)
	}
	if !validateArity(cmd.arity, cmdLine) {
		return reply.MakeArgNumErrReply(cmdName)
	}
	execFunc := cmd.exector

	//SET K V -> K V
	return execFunc(db, cmdLine[1:])
}

func validateArity(arity int, cmdArgs [][]byte) bool {
	argsNum := len(cmdArgs)
	if arity > 0 {
		return argsNum == arity
	}
	return argsNum >= -arity
}

func (db *DB) GetEntity(key string) (*database.DataEntity, bool) {
	raw, exists := db.data.Get(key)
	if !exists {
		return nil, false
	}
	entity := raw.(*database.DataEntity)

	return entity, true
}

func (db *DB) PutEntity(key string, val *database.DataEntity) int {
	return db.data.Put(key, val)
}

func (db *DB) PutIfExists(key string, val *database.DataEntity) int {
	return db.data.PutIsExists(key, val)
}

func (db *DB) PutIfAbsent(key string, val *database.DataEntity) int {
	return db.data.PutIsAbsent(key, val)
}

func (db *DB) Remove(key string) {
	db.data.Remove(key)
}

func (db *DB) Removes(keys ...string) (deleted int) {
	for _, key := range keys {
		_, exists := db.data.Get(key)
		if exists {
			db.data.Remove(key)
			deleted++
		}
	}
	return
}

func (db *DB) Flush() {
	db.data.Clear()
}
