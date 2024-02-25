package farming

import "time"

type birdhouse struct{}

func (b birdhouse) GetName() string {
	return "Birdhouse"
}

func (b birdhouse) GetTriggerTime(currentTime time.Time, cycleOffset time.Duration) (int64, error) {
	return currentTime.Add(30 * time.Minute).Unix(), nil
}
