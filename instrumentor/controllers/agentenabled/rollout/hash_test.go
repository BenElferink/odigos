package rollout

import (
	"testing"

	"github.com/stretchr/testify/assert"

	odigosv1alpha1 "github.com/odigos-io/odigos/api/odigos/v1alpha1"
)

func TestHashForContainersConfig(t *testing.T) {
	// empty containers config
	hashText, err := HashForContainersConfig([]odigosv1alpha1.ContainerAgentConfig{})
	assert.NoError(t, err)
	assert.Equal(t, len(hashText), 0)

	containersConfig := []odigosv1alpha1.ContainerAgentConfig{
		{
			ContainerName:  "container1",
			AgentEnabled:   true,
			OtelDistroName: "otel-distro-1",
		},
		{
			ContainerName:  "container2",
			AgentEnabled:   false,
			OtelDistroName: "otel-distro-2",
		},
	}

	hashText, err = HashForContainersConfig(containersConfig)
	assert.NoError(t, err)
	assert.Greater(t, len(hashText), 0)

	// flip the order of the containers and check if the hash is the same

	containersConfig = []odigosv1alpha1.ContainerAgentConfig{
		{
			ContainerName:  "container2",
			AgentEnabled:   false,
			OtelDistroName: "otel-distro-2",
		},
		{
			ContainerName:  "container1",
			AgentEnabled:   true,
			OtelDistroName: "otel-distro-1",
		},
	}

	hash2, err := HashForContainersConfig(containersConfig)
	assert.NoError(t, err)
	assert.Greater(t, len(hash2), 0)
	assert.Equal(t, hashText, hash2)

	// flip the false instrumented flag and check if the hash is different
	containersConfig[1].AgentEnabled = true
	hash3, err := HashForContainersConfig(containersConfig)
	assert.NoError(t, err)
	assert.Greater(t, len(hash3), 0)
	assert.NotEqual(t, hashText, hash3)

	// change the distro name and check if the hash is different
	containersConfig[1].OtelDistroName = "otel-distro-3"
	hash4, err := HashForContainersConfig(containersConfig)
	assert.NoError(t, err)
	assert.Greater(t, len(hash4), 0)
	assert.NotEqual(t, hashText, hash4)
}
