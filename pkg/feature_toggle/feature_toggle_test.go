package feature_toggle_test

import (
	"fmt"
	ft "github.com/khanh-nguyen-code/go_util/pkg/feature_toggle"
	"sync"
	"testing"
)

type storage struct {
	s *sync.Map
}

func (s *storage) Load(key interface{}) bool {
	if val, ok := s.s.Load(key); ok {
		if b, ok := val.(bool); ok && b {
			return true
		}
	}
	return false
}

func (s *storage) Set(key interface{}) {
	s.s.Store(key, true)
}
func (s *storage) Del(key interface{}) {
	s.s.Delete(key)
}

func TestFeatureToggle(t *testing.T) {
	s := &storage{s: &sync.Map{}}
	ft.Set(s)
	fmt.Println("set feature_1")
	s.Set("feature_1")
	ft.If("feature_1", func() {
		fmt.Println("feature_1 is enable")
	}, func() {
		fmt.Println("feature_1 is disable")
	}).If("feature_2", func() {
		fmt.Println("feature_2 is enable")
	}, func() {
		fmt.Println("feature_2 is disable")
	}).Exec()
	fmt.Println("del feature_1")
	s.Del("feature_1")
	ft.If("feature_1", func() {
		fmt.Println("feature_1 is enable")
	}, func() {
		fmt.Println("feature_1 is disable")
	}).If("feature_2", func() {
		fmt.Println("feature_2 is enable")
	}, func() {
		fmt.Println("feature_2 is disable")
	}).Exec()
}
