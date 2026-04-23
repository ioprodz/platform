package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// translations[lang][file] = parsed JSON tree for that file.
var translations = map[Lang]map[string]interface{}{}

// Load reads locales/<lang>/*.json into memory. Call once at boot.
// Missing directories are tolerated (e.g. during phase-by-phase rollout).
func Load() error {
	translations = map[Lang]map[string]interface{}{}
	base := "locales"
	for _, lang := range AllLangs {
		dir := filepath.Join(base, string(lang))
		entries, err := os.ReadDir(dir)
		if err != nil {
			if os.IsNotExist(err) {
				translations[lang] = map[string]interface{}{}
				continue
			}
			return fmt.Errorf("i18n: reading %s: %w", dir, err)
		}
		merged := map[string]interface{}{}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			if filepath.Ext(name) != ".json" {
				continue
			}
			fileKey := name[:len(name)-len(".json")]
			data, err := os.ReadFile(filepath.Join(dir, name))
			if err != nil {
				return fmt.Errorf("i18n: reading %s: %w", name, err)
			}
			var parsed map[string]interface{}
			if err := json.Unmarshal(data, &parsed); err != nil {
				return fmt.Errorf("i18n: parsing %s: %w", name, err)
			}
			merged[fileKey] = parsed
		}
		translations[lang] = merged
	}
	return nil
}
