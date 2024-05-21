package integration

import (
	"testing"

	"github.com/bitrise-io/go-utils/command"
	"github.com/stretchr/testify/require"
)

const configPath = "offline_mode.yml"

func preloadSteps(t *testing.T) {
	cmd := command.New(binPath(), "steps", "preload", "--majors=1", "--minors=1", "--minors-since=0", "--patches-since=0")
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()

	require.NoError(t, err, "Preload failed, output: %s", out)
}

func Test_GivenOfflineMode_WhenStepNotCached_ThenFails(t *testing.T) {
	preloadSteps(t)

	cmd := command.New(binPath(), "run", "not_cached", "--config", configPath)
	cmd.AppendEnvs("BITRISE_OFFLINE_MODE=true")
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()

	require.Error(t, err, "Bitrise CLI failed, output: %s", out)
	require.Contains(t, out, "Other versions available in the local cache:")
}

func Test_GivenOnlineMode_WhenStepNotCached_ThenSucceeds(t *testing.T) {
	preloadSteps(t)

	cmd := command.New(binPath(), "run", "not_cached", "--config", configPath)
	cmd.AppendEnvs("BITRISE_OFFLINE_MODE=false")
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()

	require.NoError(t, err, "Bitrise CLI failed, output: %s", out)
}

func Test_GivenOfflineMode_WhenStepCached_ThenSuceeds(t *testing.T) {
	preloadSteps(t)

	cmd := command.New(binPath(), "run", "cached", "--config", configPath)
	cmd.AppendEnvs("BITRISE_OFFLINE_MODE=true")
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()

	require.NoError(t, err, "Bitrise CLI failed, output: %s", out)
	t.Log(out)
}
