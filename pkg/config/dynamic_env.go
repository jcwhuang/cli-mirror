package config

import (
	"strings"
)

type structuredEnv struct {
	RemoteName string
	LocalName  string
}

type (
	dynamic_env struct {
		Vault *vaultEnvProvider `toml:"vault"`
	}
	vaultEnvProvider struct {
		EnvVars           []string        `toml:"env_vars"`
		StructuredEnvVars []structuredEnv `toml:"-"`
	}
)

// Get a list of env variables with possible renames like so:
// "some_env:RENAME_NAME"
// And extract them all in a structured { remoteName: "some_env", local_name: "RENAME_NAME" }
func getEnvVars(EnvVars []string) []structuredEnv {
	result := make([]structuredEnv, len(EnvVars))

	for i, envVar := range EnvVars {
		// Find first : that's not escaped with \
		colonIndex := -1
		for j := 0; j < len(envVar); j++ {
			if envVar[j] == ':' && (j == 0 || envVar[j-1] != '\\') {
				colonIndex = j
				break
			}
		}

		if colonIndex > 0 {
			// Has rename format "remote:local"
			remoteName := envVar[:colonIndex]
			localName := envVar[colonIndex+1:]

			// Unescape any \: in both names
			remoteName = strings.ReplaceAll(remoteName, "\\:", ":")
			localName = strings.ReplaceAll(localName, "\\:", ":")

			result[i] = structuredEnv{
				RemoteName: remoteName,
				LocalName:  localName,
			}
		} else {
			// No rename or escaped :, use same name for both
			// Unescape any \: in the name
			name := strings.ReplaceAll(envVar, "\\:", ":")
			result[i] = structuredEnv{
				RemoteName: name,
				LocalName:  name,
			}
		}
	}

	return result
}

func (d *dynamic_env) validate() error {
	if err := d.Vault.validate(); err != nil {
		return err
	}
	return nil
}

func (v *vaultEnvProvider) validate() error {
	if v == nil {
		return nil
	}
	v.StructuredEnvVars = getEnvVars(v.EnvVars)
	return nil
}
