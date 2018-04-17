package simplecache

import (
	"fmt"
	"strings"
	"sync"
)

//From redigo
type Error string

func (err Error) Error() string { return string(err) }

type SimpleCache struct {
	sync.Map
}

const HashMapToken = "$%HASHMAP%$"

func (s *SimpleCache) Set(key string, value interface{}) {
	v, loaded := s.LoadOrStore(key, value)

	if loaded {
		if t, ok := v.(string); ok {
			if t == HashMapToken {
				panic("simplecache: cannot call Set() on HASHMAP item")
			}
		} else {
			s.Store(key, value)
		}
	}
	return
}

func (s *SimpleCache) Get(key string) (interface{}, error) {
	if value, ok := s.Load(key); ok {
		if t, ok := value.(string); ok {
			if t == HashMapToken {
				panic("simplecache: cannot call Get() on HASHMAP item")
			}
		}

		return value, nil
	}
	return nil, ErrNil
}

func (s *SimpleCache) HSet(key string, key2 string, value interface{}) {
	v, loaded := s.LoadOrStore(key, HashMapToken)

	if loaded {
		if t, ok := v.(string); ok {
			if t != HashMapToken {
				panic("simplecache: cannot call Set() on not HASHMAP item")
			}
		}
	}

	k := fmt.Sprintf("%s/%s", key, key2)
	s.Store(k, value)
}

func (s *SimpleCache) HGet(key, key2 string) (interface{}, error) {
	k := fmt.Sprintf("%s/%s", key, key2)
	if v, ok := s.Load(k); ok {
		return v, nil
	}

	return nil, ErrNil
}

func (s *SimpleCache) HMSet(key string, args ...interface{}) {
	if len(args)%2 != 0 {
		panic("HMSet param count wrong")
	}

	v, loaded := s.LoadOrStore(key, HashMapToken)

	if loaded {
		if t, ok := v.(string); ok {
			if t != HashMapToken {
				panic("simplecache: cannot call Set() on not HASHMAP item")
			}
		}
	}

	for i := 0; i < len(args); i += 2 {
		k := fmt.Sprintf("%s/%s", key, args[i].(string))
		s.Store(k, args[i+1])
	}
}

func (s *SimpleCache) HMGet(key string, args ...interface{}) (interface{}, error) {
	values := make([]interface{}, 0, len(args))
	hasValue := false
	for i := 0; i < len(args); i++ {
		k := fmt.Sprintf("%s/%s", key, args[i].(string))
		if v, ok := s.Load(k); ok {
			values = append(values, v)
			hasValue = true
		} else {
			values = append(values, nil)
		}
	}
	if !hasValue {
		return nil, ErrNil
	}
	return values, nil
}

func (s *SimpleCache) Del(key string) error {
	v, ok := s.Load(key)

	if !ok {
		return ErrNil
	}

	if _v, ok := v.(string); ok {
		if _v == HashMapToken {
			delList := make([]string, 0, 10)
			s.Range(func(k, v interface{}) bool {
				keys := strings.SplitN(k.(string), "/", 2)
				if keys[0] == key {
					delList = append(delList, k.(string))
				}
				return true
			})

			for _, k := range delList {
				s.Delete(k)
			}

			return nil
		}
	}

	s.Delete(key)
	return nil
}
