package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConfig_Load(t *testing.T) {
	conf := Config{}

	t.Setenv("PRIMARY", "http://localhost:1337|asdf")
	t.Setenv("REPLICAS", "http://localhost:1338|qwerty")
	t.Setenv("FULL_SYNC", "false")

	err := conf.Load()
	require.NoError(t, err)

	assert.Equal(t, "http://localhost:1337", conf.Primary.Url.String())
	assert.Equal(t, "asdf", conf.Primary.Password)
	assert.Len(t, conf.Replicas, 1)
	assert.Equal(t, "http://localhost:1338", conf.Replicas[0].Url.String())
	assert.Equal(t, "qwerty", conf.Replicas[0].Password)
	assert.Equal(t, false, conf.Sync.FullSync)
}

func TestConfig_Validate_Both(t *testing.T) {
	settings := RawConfigSettings{
		DNSInclude: []string{"a"},
		DNSExclude: []string{"b"},
	}
	assert.Error(t, settings.Validate())
}

func TestConfig_Validate_Single(t *testing.T) {
	include := RawConfigSettings{
		DNSInclude: []string{"a"},
		DNSExclude: nil,
	}
	exclude := RawConfigSettings{
		DNSInclude: nil,
		DNSExclude: []string{"a"},
	}
	assert.NoError(t, include.Validate())
	assert.NoError(t, exclude.Validate())
}

func TestConfig_Validate_None(t *testing.T) {
	settings := RawConfigSettings{
		DNSInclude: nil,
		DNSExclude: nil,
	}
	assert.NoError(t, settings.Validate())
}
