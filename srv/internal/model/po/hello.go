package po

import (
	"time"
)

type Hello struct {
	Time time.Time
}

func (h *Hello) GetTime() time.Time {
	return h.Time
}

func (h *Hello) SetTime(t time.Time) {
	h.Time = t
}
