package models

import (
	"fmt"
	"time"
)

const (
	MinTime = 0
	MaxTime = 1440
)

type Slot struct {
	ID        string    `json:"id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	CreatedBy string    `json:"created_by"`
}

type CreateSlotParams struct {
	StartTime int `json:"start_time"`
	EndTime   int `json:"end_time"`
}

func (csp CreateSlotParams) Validate() map[string]string {
	errs := map[string]string{}
	if csp.StartTime < MinTime {
		errs["start_time"] = fmt.Sprintf("start time should be minimum of %d", MinTime)
	}
	if csp.StartTime > MaxTime {
		errs["start_time"] = fmt.Sprintf("start time should be maximum of %d", MaxTime)
	}
	if csp.EndTime < MinTime {
		errs["start_time"] = fmt.Sprintf("start time should be minimum of %d", MinTime)
	}
	if csp.EndTime > MaxTime {
		errs["start_time"] = fmt.Sprintf("start time should be maximum of %d", MaxTime)
	}
	if csp.StartTime >= csp.EndTime {
		errs["start_time"] = fmt.Sprintf("start time: %d, should be lessthan endtime: %d", csp.StartTime, csp.EndTime)
	}

	return errs
}
