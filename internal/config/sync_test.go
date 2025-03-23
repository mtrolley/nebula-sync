package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConfig_loadSync(t *testing.T) {
	conf := Config{}
	assert.Nil(t, conf.Sync)

	t.Setenv("FULL_SYNC", "true")
	t.Setenv("CRON", "* * * * *")
	t.Setenv("RUN_GRAVITY", "true")

	t.Setenv("SYNC_CONFIG_DNS", "true")
	t.Setenv("SYNC_CONFIG_DHCP", "true")
	t.Setenv("SYNC_CONFIG_NTP", "true")
	t.Setenv("SYNC_CONFIG_RESOLVER", "true")
	t.Setenv("SYNC_CONFIG_DATABASE", "true")
	t.Setenv("SYNC_CONFIG_MISC", "true")
	t.Setenv("SYNC_CONFIG_DEBUG", "true")

	t.Setenv("SYNC_GRAVITY_DHCP_LEASES", "true")
	t.Setenv("SYNC_GRAVITY_GROUP", "true")
	t.Setenv("SYNC_GRAVITY_AD_LIST", "true")
	t.Setenv("SYNC_GRAVITY_AD_LIST_BY_GROUP", "true")
	t.Setenv("SYNC_GRAVITY_DOMAIN_LIST", "true")
	t.Setenv("SYNC_GRAVITY_DOMAIN_LIST_BY_GROUP", "true")
	t.Setenv("SYNC_GRAVITY_CLIENT", "true")
	t.Setenv("SYNC_GRAVITY_CLIENT_BY_GROUP", "true")

	err := conf.loadSync()
	require.NoError(t, err)

	assert.Equal(t, true, conf.Sync.FullSync)
	assert.Equal(t, "* * * * *", *conf.Sync.Cron)
	assert.Equal(t, true, conf.Sync.RunGravity)

	assert.NotNil(t, conf.Sync.ConfigSettings)
	assert.NotNil(t, conf.Sync.GravitySettings)

	assert.True(t, conf.Sync.ConfigSettings.DNS.Enabled)
	assert.True(t, conf.Sync.ConfigSettings.DHCP.Enabled)
	assert.True(t, conf.Sync.ConfigSettings.NTP.Enabled)
	assert.True(t, conf.Sync.ConfigSettings.Resolver.Enabled)
	assert.True(t, conf.Sync.ConfigSettings.Database.Enabled)
	assert.True(t, conf.Sync.ConfigSettings.Misc.Enabled)
	assert.True(t, conf.Sync.ConfigSettings.Debug.Enabled)

	assert.True(t, conf.Sync.GravitySettings.DHCPLeases)
	assert.True(t, conf.Sync.GravitySettings.Group)
	assert.True(t, conf.Sync.GravitySettings.Adlist)
	assert.True(t, conf.Sync.GravitySettings.AdlistByGroup)
	assert.True(t, conf.Sync.GravitySettings.Domainlist)
	assert.True(t, conf.Sync.GravitySettings.DomainlistByGroup)
	assert.True(t, conf.Sync.GravitySettings.Client)
	assert.True(t, conf.Sync.GravitySettings.ClientByGroup)
}
