package config

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v4"
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

// func (d *dynamic_env) Fetch(ctx context.Context) (map[string]string, error) {
// 	return d.Vault.fetch(ctx)
// }

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

func (v *vaultEnvProvider) Fetch(ctx context.Context, conn *pgx.Conn) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}

	result := make(map[string]string, len(v.StructuredEnvVars))

	// Build list of remote names to query
	remoteNames := make([]string, len(v.StructuredEnvVars))
	for i, envVar := range v.StructuredEnvVars {
		remoteNames[i] = envVar.RemoteName
	}

	// Query vault for all secrets in one batch
	rows, err := conn.Query(ctx, `SELECT name, decrypted_secret FROM vault.decrypted_secrets WHERE name = ANY($1)`, remoteNames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Build map of remote name -> decrypted secret
	secrets := make(map[string]string)
	for rows.Next() {
		var name, secret string
		if err := rows.Scan(&name, &secret); err != nil {
			return nil, err
		}
		secrets[name] = secret
	}

	// Map remote secrets to local env var names
	for _, envVar := range v.StructuredEnvVars {
		if secret, ok := secrets[envVar.RemoteName]; ok {
			result[envVar.LocalName] = secret
		} else {
			// Secret not found in vault
			result[envVar.LocalName] = ""
		}
	}

	return result, nil
}
