package dict

import "sync"

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/3/5
 * Time: 15:24
 * Description: No Description
 */

type SyncDict struct {
	m sync.Map
}

func MakeSyncDict() *SyncDict {
	return &SyncDict{}
}

func (dict *SyncDict) Get(key string) (val interface{}, exists bool) {
	value, ok := dict.m.Load(key)
	return value, ok
}

func (dict *SyncDict) Len() int {
	length := 0
	dict.m.Range(func(key, value any) bool {
		length++
		return true
	})
	return length
}

func (dict *SyncDict) Put(key string, val interface{}) (result int) {
	_, exists := dict.m.Load(key)
	dict.m.Store(key, val)
	if exists {
		return 0
	}
	return 1
}

func (dict *SyncDict) PutIsAbsent(key string, val interface{}) (result int) {
	_, exists := dict.m.Load(key)
	if exists {
		return 0
	}
	dict.m.Store(key, val)
	return 1
}

func (dict *SyncDict) PutIsExists(key string, val interface{}) (result int) {
	_, exists := dict.m.Load(key)
	if exists {
		dict.m.Store(key, val)
		return 1
	}
	return 0
}

func (dict *SyncDict) Remove(key string) (result int) {
	_, exists := dict.m.Load(key)
	dict.m.Delete(key)
	if exists {
		return 1
	}
	return 0
}

func (dict *SyncDict) ForEach(consumer Consumer) {
	dict.m.Range(func(key, value any) bool {
		consumer(key.(string), value)
		return true
	})
}

func (dict *SyncDict) Keys() []string {
	result := make([]string, dict.Len())
	i := 0
	dict.m.Range(func(key, value any) bool {
		result[i] = key.(string)
		i++
		return true
	})
	return result
}

// RandomKeys 随机可重复插入的键
func (dict *SyncDict) RandomKeys(limit int) []string {
	result := make([]string, dict.Len())
	for i := 0; i < limit; i++ {
		dict.m.Range(func(key, value any) bool {
			result[i] = key.(string)
			return false
		})
	}
	return result
}

// RandomDistinctKeys 随机插入键，但不可重复
func (dict *SyncDict) RandomDistinctKeys(limit int) []string {
	result := make([]string, dict.Len())
	i := 0
	dict.m.Range(func(key, value any) bool {
		result[i] = key.(string)
		i++
		if i == limit {
			return false
		}
		return true
	})
	return result
}

func (dict *SyncDict) Clear() {
	*dict = *MakeSyncDict()
}
