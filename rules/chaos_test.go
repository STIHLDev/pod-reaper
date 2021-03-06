package rules

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
)

func TestChaosLoad(t *testing.T) {
	t.Run("load", func(t *testing.T) {
		os.Clearenv()
		os.Setenv(envChaosChance, "0.5")
		loaded, message, err := (&chaos{}).load()
		assert.NoError(t, err)
		assert.Equal(t, "chaos chance 0.5", message)
		assert.True(t, loaded)
	})
	t.Run("no load", func(t *testing.T) {
		os.Clearenv()
		loaded, message, err := (&chaos{}).load()
		assert.NoError(t, err)
		assert.Equal(t, "", message)
		assert.False(t, loaded)
	})
	t.Run("invalid chance", func(t *testing.T) {
		os.Clearenv()
		os.Setenv(envChaosChance, "not-a-number")
		loaded, message, err := (&chaos{}).load()
		assert.Error(t, err)
		assert.Equal(t, "", message)
		assert.False(t, loaded)
	})
	t.Run("explicit load without default", func(t *testing.T) {
		os.Clearenv()
		os.Setenv(envExplicitLoad, ruleChaos)
		loaded, message, err := (&chaos{}).load()
		assert.NoError(t, err)
		assert.Equal(t, "chaos (no default)", message)
		assert.True(t, loaded)
	})
}

func TestChaosShouldReap(t *testing.T) {
	t.Run("reap", func(t *testing.T) {
		os.Clearenv()
		os.Setenv(envChaosChance, "1.0") // always
		chaos := chaos{}
		chaos.load()
		shouldReap, message := chaos.ShouldReap(v1.Pod{})
		assert.True(t, shouldReap)
		assert.Equal(t, "was flagged for chaos", message)
	})
	t.Run("no reap", func(t *testing.T) {
		os.Clearenv()
		os.Setenv(envChaosChance, "0.0") // never
		chaos := chaos{}
		chaos.load()
		shouldReap, _ := chaos.ShouldReap(v1.Pod{})
		assert.False(t, shouldReap)
	})
	t.Run("annotation override reap", func(t *testing.T) {
		os.Clearenv()
		os.Setenv(envChaosChance, "0.0") // default never
		chaos := chaos{}
		chaos.load()
		pod := v1.Pod{}
		pod.Annotations = map[string]string{
			annotationChaosChance: "1.0", // override always
		}
		shouldReap, message := chaos.ShouldReap(pod)
		assert.True(t, shouldReap)
		assert.Equal(t, "was flagged for chaos", message)
	})
	t.Run("annotation override no reap", func(t *testing.T) {
		os.Clearenv()
		os.Setenv(envChaosChance, "1.0") // default always
		chaos := chaos{}
		chaos.load()
		pod := v1.Pod{}
		pod.Annotations = map[string]string{
			annotationChaosChance: "0.0", // override never
		}
		shouldReap, _ := chaos.ShouldReap(pod)
		assert.False(t, shouldReap)
	})
	t.Run("explicit load no annotation", func(t *testing.T) {
		os.Clearenv()
		os.Setenv(envExplicitLoad, ruleChaos)
		chaos := chaos{}
		chaos.load()
		shouldReap, _ := chaos.ShouldReap(v1.Pod{})
		assert.False(t, shouldReap)
	})
}
