package config

import (
	"fmt"
	"os"
	"path/filepath"

	yaml3 "gopkg.in/yaml.v3"
)

// UpdateRecordedDelay reads the existing keploy.yml config file at the given configPath,
// updates the record.recordedDelay field with the measured startup delay, and writes it back.
// This preserves all other user-configured values in the file.
func UpdateRecordedDelay(configPath string, recordedDelay uint64) error {
	configFilePath := filepath.Join(configPath, "keploy.yml")

	// Read existing config file
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", configFilePath, err)
	}

	// Unmarshal into a generic map to preserve all existing fields
	var configMap map[string]interface{}
	if err := yaml3.Unmarshal(data, &configMap); err != nil {
		return fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	// Get or create the record section
	recordSection, ok := configMap["record"]
	if !ok {
		recordSection = make(map[string]interface{})
		configMap["record"] = recordSection
	}

	recordMap, ok := recordSection.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected type for 'record' section in config")
	}

	// Update the recordedDelay value
	recordMap["recordedDelay"] = recordedDelay

	// Write back
	updatedData, err := yaml3.Marshal(configMap)
	if err != nil {
		return fmt.Errorf("failed to marshal updated config: %w", err)
	}

	if err := os.WriteFile(configFilePath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write updated config file: %w", err)
	}

	return nil
}
