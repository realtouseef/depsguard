package selection

import (
	"hash/fnv"
	"math"
	"math/rand"
	"os"
	"sort"
	"time"
)

func SelectDependencies(deps []string, percentage float64) []string {
	if len(deps) == 0 {
		return []string{}
	}
	sorted := make([]string, len(deps))
	copy(sorted, deps)
	sort.Strings(sorted)

	rng := rand.New(rand.NewSource(hashSeed(getSeed())))
	rng.Shuffle(len(sorted), func(i, j int) {
		sorted[i], sorted[j] = sorted[j], sorted[i]
	})

	count := int(math.Ceil(float64(len(sorted)) * percentage))
	if count < 1 {
		count = 1
	}
	if count > len(sorted) {
		count = len(sorted)
	}
	return sorted[:count]
}

func getSeed() string {
	for _, key := range []string{
		"GITHUB_SHA",
		"CI_COMMIT_SHA",
		"CIRCLE_SHA1",
		"TRAVIS_COMMIT",
		"BUILDKITE_COMMIT",
	} {
		if value := os.Getenv(key); value != "" {
			return value
		}
	}
	return time.Now().Format("2006-01-02")
}

func hashSeed(seed string) int64 {
	hasher := fnv.New64a()
	_, _ = hasher.Write([]byte(seed))
	return int64(hasher.Sum64())
}
