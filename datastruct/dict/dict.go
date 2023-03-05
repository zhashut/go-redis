package dict

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/3/5
 * Time: 14:38
 * Description: No Description
 */

type Consumer func(key string, val interface{}) bool

type Dict interface {
	Get(key string) (val interface{}, exists bool)
	Len() int
	Put(key string, val interface{}) (result int)
	PutIsAbsent(key string, val interface{}) (result int) // 如果没有要插入的键的话，就插入
	PutIsExists(key string, val interface{}) (result int) // 如果存在，就插入
	Remove(key string) (result int)
	ForEach(consumer Consumer) // 遍历键
	Keys() []string
	RandomKeys(limit int) []string
	RandomDistinctKeys(limit int) []string // 返回不重复的键
	Clear()
}
