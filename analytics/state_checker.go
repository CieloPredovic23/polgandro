package analytics

import (
	"github.com/bitrise-io/go-utils/v2/env"
)

const (
	// DisabledEnvKey controls both the old (analytics plugin) and new (v2) implementations
	DisabledEnvKey = "BITRISE_ANALYTICS_DISABLED"
	// V2DisabledEnvKey controls only the new (v2) implementation
	V2DisabledEnvKey = "BITRISE_ANALYTICS_V2_DISABLED"
	// V2AsyncEnvKey can be used to disable the default async queriess
	V2AsyncEnvKey = "BITRISE_ANALYTICS_V2_ASYNC"

	trueEnv = "true"
)

// StateChecker ...
type StateChecker interface {
	Enabled() bool
	UseAsync() bool
}

type stateChecker struct {
	envRepository env.Repository
}

// NewStateChecker ...
func NewStateChecker(repository env.Repository) StateChecker {
	return stateChecker{envRepository: repository}
}

// Enabled ...
func (s stateChecker) Enabled() bool {
	if s.envRepository.Get(V2DisabledEnvKey) == trueEnv {
		return false
	}

	return s.envRepository.Get(DisabledEnvKey) != trueEnv
}

func (s stateChecker) UseAsync() bool {
	return s.envRepository.Get(V2AsyncEnvKey) == "" || s.envRepository.Get(V2AsyncEnvKey) == trueEnv
}
