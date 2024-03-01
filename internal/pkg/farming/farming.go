package farming

import (
	"time"
)

const (
	TREES          = "Trees"
	HARDWOOD_TREES = "Hardwood trees"
	SPECIAL_TREES  = "Special trees"
	OTHER          = "Other"
)

type Timeable interface {
	GetName() string
	GetTriggerTime(currentTime time.Time, cycleOffset time.Duration) (int64, error)
}

var Categories = []string{TREES, HARDWOOD_TREES, SPECIAL_TREES, OTHER}

var TimeablesByCategory = map[string][]Timeable{
	TREES: {
		crop("Acorn", 40, 5*40),
		crop("Willow", 40, 6*40),
		crop("Maple", 40, 8*40),
		crop("Yew", 40, 10*40),
		crop("Magic", 40, 12*40),
	},
	HARDWOOD_TREES: {
		crop("Teak", 640, 7*640),
		crop("Mahogany", 640, 8*640),
	},
	SPECIAL_TREES: {
		crop("Calquat", 160, 8*160),
		crop("Crystal", 80, 6*80),
		crop("Spirit", 320, 12*320),
		crop("Celastrus", 160, 5*160),
		crop("Redwood", 640, 10*640),
	},
	OTHER: {
		crop("Herbs", 20, 4*20),
		crop("Fruit trees", 160, 6*160),
		crop("Hespori", 640, 3*640),
		birdhouse{},
	},
}
