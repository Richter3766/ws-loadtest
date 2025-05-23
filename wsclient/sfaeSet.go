package wsclient

import "sync"

type SafeSet struct {
	mu  sync.Mutex
	set map[int64]struct{}
}

func NewSafeSet() *SafeSet {
	return &SafeSet{set: make(map[int64]struct{})}
}

// 추가 (true면 중복, false면 새로 추가)
func (s *SafeSet) Add(id int64) (already bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.set[id]; exists {
		return true // 중복
	}
	s.set[id] = struct{}{}
	return false
}
