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
	return 100 * float32(s.Current.Load()) / float32(s.Num)
}

func (s *BuildStat) PercentageString() string {
	percent := s.Percentage()
	if percent == 100 {
		return "100.%"
	}
	percentInt := uint8(percent)
	percentDec := uint8((percent - float32(percentInt)) * 10)
	return fmt.Sprintf("%02d.%01d%%", percentInt, percentDec)
}

func (s *BuildStat) String() string {
	return fmt.Sprintf("%d/%d", s.Current.Load(), s.Num)
}
