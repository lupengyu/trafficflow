package helper

import (
	"github.com/lupengyu/trafficflow/constant"
	"math"
)

type TrackSorter struct {
	tracks []*constant.Track
}

func (s *TrackSorter) Len() int {
	return len(s.tracks)
}

func (s *TrackSorter) Swap(i, j int) {
	s.tracks[i], s.tracks[j] = s.tracks[j], s.tracks[i]
}

func (s *TrackSorter) Less(i, j int) bool {
	if math.Abs(float64(s.tracks[i].Deviation)) != math.Abs(float64(s.tracks[j].Deviation)) {
		return math.Abs(float64(s.tracks[i].Deviation)) < math.Abs(float64(s.tracks[j].Deviation))
	}
	return s.tracks[i].Deviation < s.tracks[j].Deviation
}

func (s *TrackSorter) DeWeighting() {
	if s.Len() == 0 {
		return
	}
	pre := s.tracks[0]
	newTracks := []*constant.Track{pre}
	for i := 1; i < s.Len(); i++ {
		if s.tracks[i].Deviation != pre.Deviation {
			newTracks = append(newTracks, s.tracks[i])
		}
		pre = s.tracks[i]
	}
	s.tracks = newTracks
}
