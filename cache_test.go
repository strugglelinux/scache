package cache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	cache := NewCache()
	result := cache.Set("a", "t", time.Duration(2)*time.Second)
	if !result {
		t.Errorf("key-value a=>t 未设置成功")
	}
	value, ok := cache.Get("a")
	if !ok {
		t.Errorf("键a对应值未找到 %s", value)
	}
	len := cache.Length()
	if len != 1 {
		t.Errorf("缓存元素个数不为 1 当前元素个数为:%d", len)
	}
	result = cache.Set("a", "", time.Duration(2)*time.Second)
	if !result {
		t.Errorf("key-value a=>'' 未设置成功 result:%v", result)
	}
	value, ok = cache.Get("a")
	if ok {
		t.Errorf(" 未删除成功 value:%v ok:%v", value, ok)
	}
}

func TestCacheDelete(t *testing.T) {
	cache := NewCache()
	result := cache.Set("a", "t", time.Duration(2)*time.Second)
	if !result {
		t.Errorf("key-value a=>t 未设置成功")
	}
	result = cache.Set("b", "q", time.Duration(2)*time.Second)
	if !result {
		t.Errorf("key-value a=>q 未设置成功")
	}
	len := cache.Length()
	if len != 2 {
		t.Errorf("缓存元素个数不为 2 当前元素个数为:%d", len)
	}
	result = cache.Delete("a")
	if !result {
		t.Errorf("删除 a 失败 ！ result:%v", result)
	}
	value, ok := cache.Get("a")
	if ok {
		t.Errorf(" 未删除成功 value:%v ok:%v", value, ok)
	}
}

func TestCacheRefresh(t *testing.T) {
	cache := NewCache()
	result := cache.Set("a", "t", time.Duration(2)*time.Second)
	if !result {
		t.Errorf("key-value a=>t 未设置成功")
	}
	result = cache.Set("b", "q", time.Duration(2)*time.Second)
	if !result {
		t.Errorf("key-value a=>q 未设置成功")
	}
	len := cache.Length()
	if len != 2 {
		t.Errorf("缓存元素个数不为 2 当前元素个数为:%d", len)
	}
	cache.Refresh()
	len = cache.Length()
	if len != 0 {
		t.Errorf("清空缓存失败 当前元素个数为:%d", len)
	}
}
