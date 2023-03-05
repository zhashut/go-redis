package database

import "strings"

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/3/5
 * Time: 16:18
 * Description: No Description
 */

var cmdTable = make(map[string]*Command)

type Command struct {
	exector ExecFunc
	arity   int
}

func RegisterCommand(name string, exector ExecFunc, arity int) {
	name = strings.ToLower(name)
	cmdTable[name] = &Command{
		exector: exector,
		arity:   arity,
	}
}
