package selector

import (
	"math"
	"math/rand"
)

func Select(deps []string, ratio float64, seed int64) []string {
	if len(deps) == 0 {
		return []string{}
	}
	count := int(math.Ceil(ratio * float64(len(deps))))
	if count < 1 {
		count = 1
	}
	if count > len(deps) {
		count = len(deps)
	}
	selected := make([]string, len(deps))
	copy(selected, deps)
	rng := rand.New(rand.NewSource(seed))
	rng.Shuffle(len(selected), func(i, j int) {
		selected[i], selected[j] = selected[j], selected[i]
	})
	return selected[:count]
}
