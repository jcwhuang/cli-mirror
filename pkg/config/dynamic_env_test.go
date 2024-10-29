package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvVars(t *testing.T) {
	t.Run("simple env vars without rename", func(t *testing.T) {
		input := []string{"SIMPLE_VAR", "ANOTHER_VAR"}
		result := getEnvVars(input)

		assert.Equal(t, 2, len(result))
		assert.Equal(t, "SIMPLE_VAR", result[0].RemoteName)
		assert.Equal(t, "SIMPLE_VAR", result[0].LocalName)
		assert.Equal(t, "ANOTHER_VAR", result[1].RemoteName)
		assert.Equal(t, "ANOTHER_VAR", result[1].LocalName)
	})

	t.Run("env vars with rename", func(t *testing.T) {
		input := []string{"remote:LOCAL", "another-remote:LOCAL_NAME"}
		result := getEnvVars(input)

		assert.Equal(t, 2, len(result))
		assert.Equal(t, "remote", result[0].RemoteName)
		assert.Equal(t, "LOCAL", result[0].LocalName)
		assert.Equal(t, "another-remote", result[1].RemoteName)
		assert.Equal(t, "LOCAL_NAME", result[1].LocalName)
	})

	t.Run("env vars with escaped colons", func(t *testing.T) {
		input := []string{
			"remote\\:with\\:colons:LOCAL",
			"remote\\:colon:LOCAL\\:NAME",
		}
		result := getEnvVars(input)

		assert.Equal(t, 2, len(result))
		assert.Equal(t, "remote:with:colons", result[0].RemoteName)
		assert.Equal(t, "LOCAL", result[0].LocalName)
		assert.Equal(t, "remote:colon", result[1].RemoteName)
		assert.Equal(t, "LOCAL:NAME", result[1].LocalName)
	})

	t.Run("env vars with spaces in local name", func(t *testing.T) {
		input := []string{
			"remote:LOCAL NAME WITH SPACES",
			"another-remote:SPACED LOCAL NAME",
		}
		result := getEnvVars(input)

		assert.Equal(t, 2, len(result))
		assert.Equal(t, "remote", result[0].RemoteName)
		assert.Equal(t, "LOCAL NAME WITH SPACES", result[0].LocalName)
		assert.Equal(t, "another-remote", result[1].RemoteName)
		assert.Equal(t, "SPACED LOCAL NAME", result[1].LocalName)
	})

	t.Run("env vars with spaces in remote name", func(t *testing.T) {
		input := []string{
			"remote name with spaces:LOCAL",
			"another remote name:LOCAL_NAME",
		}
		result := getEnvVars(input)

		assert.Equal(t, 2, len(result))
		assert.Equal(t, "remote name with spaces", result[0].RemoteName)
		assert.Equal(t, "LOCAL", result[0].LocalName)
		assert.Equal(t, "another remote name", result[1].RemoteName)
		assert.Equal(t, "LOCAL_NAME", result[1].LocalName)
	})

	t.Run("empty input", func(t *testing.T) {
		result := getEnvVars([]string{})
		assert.Equal(t, 0, len(result))
	})
}

func TestDynamicEnvValidate(t *testing.T) {
	t.Run("valid vault config", func(t *testing.T) {
		config := &dynamic_env{
			Vault: &vaultEnvProvider{
				EnvVars: []string{
					"remote:LOCAL",
					"another:LOCAL_NAME",
				},
			},
		}
		err := config.validate()
		assert.NoError(t, err)
		assert.Equal(t, 2, len(config.Vault.StructuredEnvVars))
	})

	t.Run("nil vault config", func(t *testing.T) {
		config := &dynamic_env{
			Vault: nil,
		}
		err := config.validate()
		assert.NoError(t, err)
	})

}
