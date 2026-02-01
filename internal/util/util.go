package util

import (
	"hash/fnv"
	"os"
	"time"
)

func DepguardDir() string {
	return ".depguard"
}

func ResolveSeed() int64 {
	sha := commitSHA()
	if sha != "" {
		hasher := fnv.New64a()
		_, _ = hasher.Write([]byte(sha))
		return int64(hasher.Sum64())
	}
	return time.Now().UnixNano()
}

func commitSHA() string {
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
	return ""
}

func CurrentUser() string {
	if value := os.Getenv("USER"); value != "" {
		return value
	}
	if value := os.Getenv("USERNAME"); value != "" {
		return value
	}
	return "unknown"
}
