package engine

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"

	"github.com/vimal/sams/internal/models"
)

var (
	cachedRules []models.Rule
	rulesOnce   sync.Once
	rulesErr    error
)

// LoadRules loads all YAML rule files from rulesPath recursively.
// Rules are cached in memory after the first load.
func LoadRules(rulesPath string) ([]models.Rule, error) {
	rulesOnce.Do(func() {
		cachedRules, rulesErr = loadRulesFromDisk(rulesPath)
	})
	return cachedRules, rulesErr
}

// ReloadRules forces a fresh load (used in tests / hot-reload scenarios).
func ReloadRules(rulesPath string) ([]models.Rule, error) {
	rulesOnce = sync.Once{}
	return LoadRules(rulesPath)
}

func loadRulesFromDisk(rulesPath string) ([]models.Rule, error) {
	var rules []models.Rule

	err := filepath.WalkDir(rulesPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}

		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return fmt.Errorf("reading %s: %w", path, readErr)
		}

		// Each file may contain a single rule or a list.
		// Try list first, fall back to single.
		var ruleList []models.Rule
		if yamlErr := yaml.Unmarshal(data, &ruleList); yamlErr == nil && len(ruleList) > 0 {
			rules = append(rules, ruleList...)
			return nil
		}

		var single models.Rule
		if yamlErr := yaml.Unmarshal(data, &single); yamlErr != nil {
			return fmt.Errorf("parsing %s: %w", path, yamlErr)
		}
		if single.ID != "" {
			rules = append(rules, single)
		}
		return nil
	})

	return rules, err
}
