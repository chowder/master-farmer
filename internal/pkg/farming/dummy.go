package farming

import "time"

type dummy struct{}

func (d dummy) GetName() string {
	return "Dummy for debugging"
}

func (d dummy) GetTriggerTime(currentTime time.Time, cycleOffset time.Duration) (int64, error) {
	return currentTime.Add(5 * time.Second).Unix(), nil
}
