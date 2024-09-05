package extension

import "time"

type MemoryStorageItem struct {
	Value   any
	Expired time.Time
}

type MemoryStorage struct {
	data map[string]MemoryStorageItem
}

func (s *MemoryStorage) Init() {
	if s.data == nil {
		s.data = map[string]MemoryStorageItem{}
	}
}

func (s *MemoryStorage) Get(key string) (any, bool) {
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

func (s *MemoryStorage) Set(key string, value any, ttl ...int) {
	item := MemoryStorageItem{
		Value: value,
	}
	for _, t := range ttl {
		if t > 0 {
			item.Expired = time.Now().Add(time.Duration(t) * time.Second)
		}
	}
	s.data[key] = item
}

func (s *MemoryStorage) Remove(key string) {
	delete(s.data, key)
}
