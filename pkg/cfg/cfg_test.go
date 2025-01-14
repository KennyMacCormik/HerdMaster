package cfg

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

type TestConfig struct {
	Field1 string `mapstructure:"field1"`
	Field2 int    `mapstructure:"field2"`
}

func cleanConfigEntries() {
	mtx.Lock()
	defer mtx.Unlock()
	configEntries = make(map[string]ConfigEntry)
}

func TestRegisterConfig_Success(t *testing.T) {
	defer cleanConfigEntries()
	entry := ConfigEntry{
		Config: &TestConfig{},
		BindArray: []BindValue{
			{ValName: "field1"},
			{ValName: "field2"},
		},
	}

	err := RegisterConfig("test", entry)
	assert.NoError(t, err, "expected no error when registering a valid config")

	configs := ListConfigs()
	assert.Contains(t, configs, "test", "expected registered config to be listed")
}

func TestRegisterConfig_ValidationError(t *testing.T) {
	defer cleanConfigEntries()
	entry := ConfigEntry{
		Config: "invalid",
		BindArray: []BindValue{
			{ValName: "field1"},
		},
	}

	err := RegisterConfig("invalid", entry)
	assert.Error(t, err, "expected an error for invalid config registration")
	assert.Contains(t, err.Error(), "config validation failed", "expected validation error message")
}

func TestListConfigs(t *testing.T) {
	defer cleanConfigEntries()
	entry1 := ConfigEntry{Config: &TestConfig{}, BindArray: []BindValue{{ValName: "field1"}}}
	entry2 := ConfigEntry{Config: &TestConfig{}, BindArray: []BindValue{{ValName: "field2"}}}

	_ = RegisterConfig("config1", entry1)
	_ = RegisterConfig("config2", entry2)

	configs := ListConfigs()
	assert.ElementsMatch(t, configs, []string{"config1", "config2"}, "expected all registered configs to be listed")
}

func TestGetConfig(t *testing.T) {
	defer cleanConfigEntries()
	entry := ConfigEntry{
		Config: &TestConfig{},
		BindArray: []BindValue{
			{ValName: "field1"},
		},
	}

	_ = RegisterConfig("test", entry)

	config, ok := GetConfig("test")
	assert.True(t, ok, "expected to find registered config")
	assert.IsType(t, &TestConfig{}, config, "expected config to be of type TestConfig")
}

func TestNewConfig_BindAndUnmarshal(t *testing.T) {
	defer cleanConfigEntries()
	entry := ConfigEntry{
		Config: &TestConfig{},
		BindArray: []BindValue{
			{ValName: "field1", DefaultVal: "default1"},
			{ValName: "field2", DefaultVal: 123},
		},
	}

	_ = RegisterConfig("test", entry)

	err := NewConfig()
	assert.NoError(t, err, "expected no error when initializing configs")

	config, ok := GetConfig("test")
	assert.True(t, ok, "expected to find registered config")

	typedConfig, ok := config.(*TestConfig)
	assert.True(t, ok, "expected config to be of type TestConfig")
	assert.Equal(t, "default1", typedConfig.Field1, "expected default value for Field1")
	assert.Equal(t, 123, typedConfig.Field2, "expected default value for Field2")
}

func TestValidateBindArray_Error(t *testing.T) {
	defer cleanConfigEntries()
	entry := ConfigEntry{
		Config: &TestConfig{},
		BindArray: []BindValue{
			{ValName: ""},
		},
	}

	err := RegisterConfig("test", entry)
	assert.Error(t, err, "expected an error for invalid BindValue")
	assert.Contains(t, err.Error(), "BindValue.ValName cannot be empty", "expected validation error for empty ValName")
}

func TestEmptyConfig(t *testing.T) {
	defer cleanConfigEntries()

	entry := ConfigEntry{
		Config:    &TestConfig{},
		BindArray: []BindValue{},
	}

	err := RegisterConfig("empty", entry)
	assert.NoError(t, err, "expected no error when registering an empty config entry")
}

func TestDuplicateRegistration(t *testing.T) {
	defer cleanConfigEntries()

	entry := ConfigEntry{
		Config:    &TestConfig{},
		BindArray: []BindValue{{ValName: "field1"}},
	}

	err := RegisterConfig("duplicate", entry)
	assert.NoError(t, err, "expected no error when registering a valid config")

	// Register again with the same name
	err = RegisterConfig("duplicate", entry)
	assert.NoError(t, err, "expected no error when re-registering the same config")
}

func TestEmptyConfigName(t *testing.T) {
	defer cleanConfigEntries()

	entry := ConfigEntry{
		Config:    &TestConfig{},
		BindArray: []BindValue{{ValName: "field1"}},
	}

	err := RegisterConfig("", entry)
	assert.Error(t, err, "expected an error for an empty config name")
}

func TestTypeMismatch(t *testing.T) {
	defer cleanConfigEntries()

	entry := ConfigEntry{
		Config:    &TestConfig{},
		BindArray: []BindValue{{ValName: "field1"}},
	}

	err := RegisterConfig("test", entry)
	assert.NoError(t, err, "expected no error when registering a valid config")

	config, ok := GetConfig("test")
	assert.True(t, ok, "expected to find registered config")

	_, ok = config.(*struct{})
	assert.False(t, ok, "expected type assertion to fail for incorrect type")
}

func TestConcurrentRegisterConfig(t *testing.T) {
	defer cleanConfigEntries()

	const numConfigs = 100
	var wg sync.WaitGroup
	wg.Add(numConfigs)

	for i := 0; i < numConfigs; i++ {
		go func(i int) {
			defer wg.Done()
			name := fmt.Sprintf("config%d", i)
			entry := ConfigEntry{
				Config:    &TestConfig{},
				BindArray: []BindValue{{ValName: fmt.Sprintf("field%d", i)}},
			}
			_ = RegisterConfig(name, entry)
		}(i)
	}

	wg.Wait()
	assert.Equal(t, numConfigs, len(ListConfigs()), "expected all concurrently registered configs to be listed")
}

func TestEnvironmentVariableMocking(t *testing.T) {
	defer cleanConfigEntries()

	entry := ConfigEntry{
		Config: &TestConfig{},
		BindArray: []BindValue{
			{ValName: "field1", DefaultVal: "default1"},
			{ValName: "field2", DefaultVal: 123},
		},
	}

	err := RegisterConfig("test", entry)
	assert.NoError(t, err, "expected no error when registering a valid config")

	t.Setenv("HM_FIELD1", "env_value1")
	t.Setenv("HM_FIELD2", "456")

	err = NewConfig()
	assert.NoError(t, err, "expected no error when initializing configs")

	config, ok := GetConfig("test")
	assert.True(t, ok, "expected to find registered config")

	typedConfig, ok := config.(*TestConfig)
	assert.True(t, ok, "expected config to be of type TestConfig")
	assert.Equal(t, "env_value1", typedConfig.Field1, "expected Field1 to be overridden by environment variable")
	assert.Equal(t, 456, typedConfig.Field2, "expected Field2 to be overridden by environment variable")
}
