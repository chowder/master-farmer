package farming

import (
	"github.com/chowder/master-farmer/internal/pkg/utils"
	"time"
)

var FarmingEpoch = time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC)

type _crop struct {
	name        string
	growthCycle int
	growthTime  int
}

func crop(name string, growthCycle int, growthTime int) _crop {
	return _crop{
		name:        name,
		growthCycle: growthCycle,
		growthTime:  growthTime,
	}
}

func (c _crop) GetTriggerTime(currentTime time.Time, cycleOffset time.Duration) (int64, error) {
	// Calculate day 1 in the farming cycle
	daysSinceEpoch := int(currentTime.Sub(FarmingEpoch).Hours() / 24)
	referenceDate := utils.StripTime(currentTime).Add(-time.Duration(daysSinceEpoch%4) * 24 * time.Hour)

	// Calculate when the next farming tick will be at
	elapsedMinutes := int(currentTime.Sub(referenceDate).Minutes())
	nextTickMinutes := ((elapsedMinutes / c.growthCycle) + 1) * c.growthCycle
	nextTickAt := referenceDate.Add(time.Duration(nextTickMinutes) * time.Minute)

	// Account for offset
	nextTickAt = nextTickAt.Add(-cycleOffset)
	if nextTickAt.Before(currentTime) {
		nextTickAt = nextTickAt.Add(time.Duration(c.growthCycle) * time.Minute)
	}

	// Calculate when the crop will finish growing
	completeAt := nextTickAt.Add(time.Duration(c.growthTime-c.growthCycle) * time.Minute)
	return completeAt.Unix(), nil
}

func (c _crop) GetName() string { return c.name }
