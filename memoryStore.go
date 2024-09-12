package extension

import "time"

type MemoryStoreItem struct {
	Value   any
	Expired time.Time
}

type MemoryStore struct {
	data map[string]MemoryStoreItem
}

func (s *MemoryStore) Init() {
	if s.data == nil {
		s.data = map[string]MemoryStoreItem{}
	}
}

func (s *MemoryStore) Get(key string) (any, bool) {
	item, isExist := s.data[key]
	if !isExist {
		return nil, isExist
	}
	if !item.Expired.IsZero() && time.Now().After(item.Expired) {
		s.Remove(key)
		return nil, false
	}
	return item.Value, isExist
}

func (s *MemoryStore) Set(key string, value any, ttl ...int) {
	item := MemoryStoreItem{
		Value: value,
	}
	if len(ttl) > 0 {
		t := ttl[0]
		item.Expired = time.Now().Add(time.Duration(t) * time.Second)
	}

	s.data[key] = item
}

func (s *MemoryStore) Remove(key string) {
	delete(s.data, key)
}
