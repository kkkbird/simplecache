package simplecache

import (
	"sync"
)

const HashMapToken = "$%HASHMAP%$"

type Error string

func (err Error) Error() string { return string(err) }

type SimpleCache struct {
	m sync.Map
}

func (s *SimpleCache) Set(key string, value interface{}) {
	if _, ok := s.m.Load(HashMapToken + key); ok {
		panic("simplecache: cannot call Set() on HASHMAP item")
	}

	s.m.Store(key, value)
}

func (s *SimpleCache) Get(key string) (interface{}, error) {
	if _, ok := s.m.Load(HashMapToken + key); ok {
		panic("simplecache: cannot call Get() on HASHMAP item")
	}

	if value, ok := s.m.Load(key); ok {
		return value, nil
	}
	return nil, ErrNil
}

func (s *SimpleCache) HSet(key string, key2 string, value interface{}) {
	s.HMSet(key, key2, value)
}

func (s *SimpleCache) HGet(key, key2 string) (interface{}, error) {
	r, err := s.HMGet(key, key2)

	if err != nil {
		return nil, err
	}
	return r.([]interface{})[0], err
}

func (s *SimpleCache) HMSet(key string, args ...interface{}) {
	if len(args)%2 != 0 {
		panic("HMSet param count wrong")
	}

	if _, ok := s.m.Load(key); ok {
		panic("simplecache: cannot call HMSet() on NOT-HASHMAP item")
	}

	key = HashMapToken + key

	//create new map
	hm := make(map[string]interface{})

	for i := 0; i < len(args); i += 2 {
		hm[args[i].(string)] = args[i+1]
	}

	_existHMap, loaded := s.m.LoadOrStore(key, hm)

	//if key already exists
	if loaded {
		existHMap, ok := _existHMap.(map[string]interface{})
		// value must be a hashmap type
		if !ok {
			panic("simplecache: cannot call HMSet() on NOT-HASHMAP item")
		}
		for k, v := range hm {
			existHMap[k] = v
		}
		s.m.Store(key, existHMap)
	}
}

func (s *SimpleCache) HMGet(key string, args ...interface{}) (interface{}, error) {
	if _, ok := s.m.Load(key); ok {
		panic("simplecache: cannot call HMSet() on NOT-HASHMAP item")
	}
	key = HashMapToken + key
	if _hm, ok := s.m.Load(key); ok {
		hm, ok := _hm.(map[string]interface{})
		if !ok {
			panic("simplecache: cannot call HMGet() on NOT-HASHMAP item")
		}
		values := make([]interface{}, 0, len(args))
		hasValue := false
		for i := 0; i < len(args); i++ {
			if v, ok := hm[args[i].(string)]; ok {
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
	return nil, ErrNil
}

func (s *SimpleCache) Del(key string) error {
	s.m.Delete(key)
	s.m.Delete(HashMapToken + key)
	return nil
}

// func (s *SimpleCache) Set(key string, value interface{}) {
// 	v, loaded := s.m.LoadOrStore(key, value)

// 	if loaded {
// 		if t, ok := v.(string); ok {
// 			if t == HashMapToken {
// 				panic("simplecache: cannot call Set() on HASHMAP item")
// 			}
// 		}
// 		s.m.Store(key, value)
// 	}
// 	return
// }

// func (s *SimpleCache) Get(key string) (interface{}, error) {
// 	if value, ok := s.m.Load(key); ok {
// 		if t, ok := value.(string); ok {
// 			if t == HashMapToken {
// 				panic("simplecache: cannot call Get() on HASHMAP item")
// 			}
// 		}

// 		return value, nil
// 	}
// 	return nil, ErrNil
// }

// func (s *SimpleCache) HSet(key string, key2 string, value interface{}) {
// 	v, loaded := s.m.LoadOrStore(key, HashMapToken)

// 	if loaded {
// 		t, ok := v.(string)
// 		if !ok || t != HashMapToken {
// 			panic("simplecache: cannot call Set() on not HASHMAP item")
// 		}
// 	}

// 	k := strings.Join([]string{key, key2}, "/")

// 	s.m.Store(k, value)
// }

// func (s *SimpleCache) HGet(key, key2 string) (interface{}, error) {
// 	k := strings.Join([]string{key, key2}, "/")
// 	if v, ok := s.m.Load(k); ok {
// 		return v, nil
// 	}

// 	return nil, ErrNil
// }

// func (s *SimpleCache) HMSet(key string, args ...interface{}) {
// 	if len(args)%2 != 0 {
// 		panic("HMSet param count wrong")
// 	}

// 	v, loaded := s.m.LoadOrStore(key, HashMapToken)

// 	if loaded {
// 		if t, ok := v.(string); ok {
// 			if t != HashMapToken {
// 				panic("simplecache: cannot call Set() on not HASHMAP item")
// 			}
// 		}
// 	}

// 	for i := 0; i < len(args); i += 2 {
// 		k := strings.Join([]string{key, args[i].(string)}, "/")
// 		s.m.Store(k, args[i+1])
// 	}
// }

// func (s *SimpleCache) HMGet(key string, args ...interface{}) (interface{}, error) {
// 	l := len(args)
// 	values := make([]interface{}, 0, l)
// 	hasValue := false
// 	for i := 0; i < l; i++ {
// 		k := strings.Join([]string{key, args[i].(string)}, "/")
// 		if v, ok := s.m.Load(k); ok {
// 			values = append(values, v)
// 			hasValue = true
// 		} else {
// 			values = append(values, nil)
// 		}
// 	}
// 	if !hasValue {
// 		return nil, ErrNil
// 	}
// 	return values, nil
// }

// func (s *SimpleCache) Del(key string) error {
// 	v, ok := s.m.Load(key)

// 	if !ok {
// 		return ErrNil
// 	}

// 	if _v, ok := v.(string); ok {
// 		if _v == HashMapToken {
// 			//TODO: low effeciency to enumerate all keys
// 			delList := make([]string, 0)
// 			keydir := key + "/"
// 			lKeydir := len(keydir)
// 			s.m.Range(func(k, v interface{}) bool {
// 				_k := k.(string)
// 				if len(_k) > lKeydir && _k[:lKeydir] == keydir {
// 					delList = append(delList, _k)
// 				}
// 				return true
// 			})

// 			fmt.Println(delList)

// 			for _, k := range delList {
// 				s.m.Delete(k)
// 			}
// 			s.m.Delete(key)

// 			return nil
// 		}
// 	}

// 	s.m.Delete(key)
// 	return nil
// }
