package cfg

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestConfig struct {
	Field1 string `mapstructure:"field1" validate:"required"`
	Field2 int    `mapstructure:"field2" validate:"numeric"`
}

func resetConfigEntries() {
	mtx.Lock()
	defer mtx.Unlock()
	configEntries = make(map[string]ConfigEntry)
}

func TestRegisterConfig_ValidConfig(t *testing.T) {
	defer resetConfigEntries()

	entry := ConfigEntry{
		Config: &TestConfig{},
		BindArray: []BindValue{
			{ValName: "field1", DefaultVal: "default1"},
			{ValName: "field2", DefaultVal: 123},
		},
	}

	err := RegisterConfig("testConfig", entry)
	assert.NoError(t, err, "expected no error for valid config registration")

	entries := ListConfigs()
	assert.Contains(t, entries, "testConfig", "expected config name to be listed")
}

func TestRegisterConfig_InvalidConfig(t *testing.T) {
	defer resetConfigEntries()

	entry := ConfigEntry{
		Config: "invalid",
		BindArray: []BindValue{
			{ValName: "field1"},
		},
	}

	err := RegisterConfig("invalidConfig", entry)
	assert.Error(t, err, "expected error for invalid config struct")
	assert.Contains(t, err.Error(), "config validation failed", "expected config validation error")
}

func TestRegisterConfig_EmptyName(t *testing.T) {
	defer resetConfigEntries()

	entry := ConfigEntry{
		Config:    &TestConfig{},
		BindArray: []BindValue{{ValName: "field1"}},
	}

	err := RegisterConfig("", entry)
	assert.Error(t, err, "expected error for empty config name")
	assert.Contains(t, err.Error(), "name cannot be empty", "expected error for empty name")
}

func TestListConfigs(t *testing.T) {
	defer resetConfigEntries()

	entry1 := ConfigEntry{
		Config:    &TestConfig{},
		BindArray: []BindValue{{ValName: "field1"}},
	}
	entry2 := ConfigEntry{
		Config:    &TestConfig{},
		BindArray: []BindValue{{ValName: "field2"}},
	}

	RegisterConfig("config1", entry1)
	RegisterConfig("config2", entry2)

	entries := ListConfigs()
	assert.ElementsMatch(t, entries, []string{"config1", "config2"}, "expected all registered configs to be listed")
}

func TestGetConfig(t *testing.T) {
	defer resetConfigEntries()

	entry := ConfigEntry{
		Config:    &TestConfig{},
		BindArray: []BindValue{{ValName: "field1"}},
	}

	RegisterConfig("testConfig", entry)

	config, ok := GetConfig("testConfig")
	assert.True(t, ok, "expected to retrieve registered config")
	assert.IsType(t, &TestConfig{}, config, "expected config to match registered type")
}

func TestGetConfig_Nonexistent(t *testing.T) {
	defer resetConfigEntries()

	config, ok := GetConfig("nonexistentConfig")
	assert.False(t, ok, "expected false for nonexistent config")
	assert.Nil(t, config, "expected nil config for nonexistent entry")
}

func TestBindActualValue_Valid(t *testing.T) {
	defer resetConfigEntries()

	entry := ConfigEntry{
		Config: &TestConfig{},
		BindArray: []BindValue{
			{ValName: "field1", DefaultVal: "default1"},
			{ValName: "field2", DefaultVal: 123},
		},
	}
	RegisterConfig("testConfig", entry)

	err := newConfigWithLock()
	assert.NoError(t, err, "expected no error when binding actual values")

	config, ok := GetConfig("testConfig")
	assert.True(t, ok, "expected to retrieve registered config")

	typedConfig, ok := config.(*TestConfig)
	assert.True(t, ok, "expected type assertion to succeed")
	assert.Equal(t, "default1", typedConfig.Field1, "expected default value for Field1")
	assert.Equal(t, 123, typedConfig.Field2, "expected default value for Field2")
}

func TestBindActualValue_MissingField(t *testing.T) {
	defer resetConfigEntries()

	entry := ConfigEntry{
		Config: &TestConfig{},
		BindArray: []BindValue{
			{ValName: ""},
		},
	}
	err := RegisterConfig("testConfig", entry)
	assert.Error(t, err, "expected error for missing ValName")
	assert.Contains(t, err.Error(), "BindValue.ValName cannot be empty", "expected validation error for empty ValName")
}

func TestWithSetEnvPrefix_Valid(t *testing.T) {
	err := NewConfig(WithSetEnvPrefix("myapp"))
	assert.NoError(t, err, "expected no error for valid environment prefix")
}

func TestWithSetEnvPrefix_Empty(t *testing.T) {
	err := NewConfig(WithSetEnvPrefix(""))
	assert.Error(t, err, "expected error for empty environment prefix")
	assert.Contains(t, err.Error(), "incorrect env prefix", "expected error for invalid prefix")
}

func TestNewConfig_MultipleOptions(t *testing.T) {
	err := NewConfig(WithSetEnvPrefix("test"), func() error {
		viper.Set("customOption", true)
		return nil
	})
	assert.NoError(t, err, "expected no error for valid functional options")
}
