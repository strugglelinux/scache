package cache

import (
	"sync"
	"sync/atomic"
	"time"
)

const (
	//DefaultExpiration 缓存默认过期时间
	DefaultExpiration time.Duration = 0
)

//Item 元素
type Item struct {
	Object     interface{}
	ExpireTime int64
}

//IsExpired  判断是否过期
func (item Item) IsExpired() bool {
	if item.ExpireTime == 0 { //为零时不过期
		return false
	}
	return time.Now().UnixNano() > item.ExpireTime
}

//Cache 结构体
type Cache struct {
	mu    sync.RWMutex
	items map[string]Item //缓存值
	num   int64
}

//Get 获取value
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.num == 0 || key == "" { //
		return nil, false
	}
	item, found := c.items[key]
	if found {
		if isExprired := item.IsExpired(); isExprired {
			c.delete(key)
			return nil, false
		}
		return item.Object, true
	}
	return nil, false
}

func (c *Cache) get(key string) (interface{}, bool) {
	item, found := c.items[key]
	if found {
		if isExprired := item.IsExpired(); isExprired {
			c.delete(key)
			return nil, false
		}
		return item.Object, true
	}
	return nil, false

}

// Set 设置
func (c *Cache) Set(key string, value interface{}, expireTime time.Duration) bool {
	var t int64
	if expireTime > DefaultExpiration {
		t = time.Now().Add(expireTime).UnixNano()
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if value == "" || value == nil {
		c.delete(key)
		return true
	}
	c.items[key] = Item{value, t}
	atomic.AddInt64(&c.num, 1)
	return true
}

//Delete 删除指定key缓存值
func (c *Cache) Delete(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	result := c.delete(key)
	return result
}

func (c *Cache) delete(key string) bool {
	if _, ok := c.items[key]; ok {
		delete(c.items, key)
		atomic.AddInt64(&c.num, -1)
		return true
	}
	return false
}

//Refresh 刷新缓存 清空
func (c *Cache) Refresh() {
	c.mu.Lock()
	c.items = make(map[string]Item)
	c.num = 0
	c.mu.Unlock()
}

//Length 获取元素个数
func (c *Cache) Length() int64 {
	return c.num
}

//NewCache 创建
func NewCache() *Cache {
	c := new(Cache)
	c.items = make(map[string]Item)
	c.num = 0
	return c
}
