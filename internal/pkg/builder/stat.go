package builder

import (
	"fmt"
	"sync/atomic"
)

type BuildStat struct {
	Num     uint32
	Current atomic.Uint32
}

func (s *BuildStat) SetNum(num uint32) {
	s.Num = num
	s.Current.Store(uint32(0))
}

func (s *BuildStat) Done() {
	s.Current.Add(uint32(1))
}

func (s *BuildStat) Percentage() float32 {
	return float32(s.Current.Load()) / float32(s.Num)
}

func (s *BuildStat) String() string {
	return fmt.Sprintf("%d/%d", s.Current.Load(), s.Num)
}
