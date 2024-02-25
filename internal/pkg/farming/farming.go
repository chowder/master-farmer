package farming

import (
	"time"
)

const (
	TREES          = "Trees"
	FRUIT_TREES    = "Fruit trees"
	HARDWOOD_TREES = "Hardwood trees"
	SPECIAL_TREES  = "Special trees"
	MISC           = "Misc"
)

type Timeable interface {
	GetName() string
	GetTriggerTime(currentTime time.Time, cycleOffset time.Duration) (int64, error)
}

var Categories = []string{TREES, FRUIT_TREES, HARDWOOD_TREES, SPECIAL_TREES, MISC}

var TimeablesByCategory = map[string][]Timeable{
	TREES: {
		crop("Acorn", 40, 5*40),
		crop("Willow", 40, 6*40),
		crop("Maple", 40, 8*40),
		crop("Yew", 40, 10*40),
		crop("Magic", 40, 12*40),
	},
	FRUIT_TREES: {
		crop("Apple", 160, 6*160),
		crop("Banana", 160, 6*160),
		crop("Orange", 160, 6*160),
		crop("Curry", 160, 6*160),
		crop("Pineapple", 160, 6*160),
		crop("Papaya", 160, 6*160),
		crop("Palm tree", 160, 6*160),
		crop("Dragonfruit", 160, 6*160),
	},
	HARDWOOD_TREES: {
		crop("Teak", 640, 7*640),
		crop("Mahogany", 640, 8*640),
	},
	SPECIAL_TREES: {
		crop("Calquat", 160, 8*160),
		crop("Crystal Acorn", 80, 6*80),
		crop("Celastrus", 160, 5*160),
		crop("Redwood", 640, 10*640),
	},
	MISC: {
		crop("Hespori", 640, 3*640),
		birdhouse{},
		dummy{},
	},
}
