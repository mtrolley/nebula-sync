package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/lovelaze/nebula-sync/internal/filter"
	"github.com/lovelaze/nebula-sync/internal/pihole/model"
)

type Config struct {
	Primary  model.PiHole   `required:"true" envconfig:"PRIMARY"`
	Replicas []model.PiHole `required:"true" envconfig:"REPLICAS"`
	Client   *Client        `ignored:"true"`
	Sync     *Sync          `ignored:"true"`
}

type ConfigSettings struct {
	DNS       *ConfigSetting
	DHCP      *ConfigSetting
	NTP       *ConfigSetting
	Resolver  *ConfigSetting
	Database  *ConfigSetting
	Webserver *ConfigSetting
	Files     *ConfigSetting
	Misc      *ConfigSetting
	Debug     *ConfigSetting
}

type RawConfigSettings struct {
	DNS             bool     `default:"false" envconfig:"SYNC_CONFIG_DNS"`
	DNSInclude      []string `envconfig:"SYNC_CONFIG_DNS_INCLUDE"`
	DNSExclude      []string `envconfig:"SYNC_CONFIG_DNS_EXCLUDE"`
	DHCP            bool     `default:"false" envconfig:"SYNC_CONFIG_DHCP"`
	DHCPInclude     []string `envconfig:"SYNC_CONFIG_DHCP_INCLUDE"`
	DHCPExclude     []string `envconfig:"SYNC_CONFIG_DHCP_EXCLUDE"`
	NTP             bool     `default:"false" envconfig:"SYNC_CONFIG_NTP"`
	NTPInclude      []string `envconfig:"SYNC_CONFIG_NTP_INCLUDE"`
	NTPExclude      []string `envconfig:"SYNC_CONFIG_NTP_EXCLUDE"`
	Resolver        bool     `default:"false" envconfig:"SYNC_CONFIG_RESOLVER"`
	ResolverInclude []string `envconfig:"SYNC_CONFIG_RESOLVER_INCLUDE"`
	ResolverExclude []string `envconfig:"SYNC_CONFIG_RESOLVER_EXCLUDE"`
	Database        bool     `default:"false" envconfig:"SYNC_CONFIG_DATABASE"`
	DatabaseInclude []string `envconfig:"SYNC_CONFIG_DATABASE_INCLUDE"`
	DatabaseExclude []string `envconfig:"SYNC_CONFIG_DATABASE_EXCLUDE"`
	Webserver       bool     `default:"false" ignored:"true"` // ignore for now
	Files           bool     `default:"false" ignored:"true"` // ignore for now
	Misc            bool     `default:"false" envconfig:"SYNC_CONFIG_MISC"`
	MiscInclude     []string `envconfig:"SYNC_CONFIG_MISC_INCLUDE"`
	MiscExclude     []string `envconfig:"SYNC_CONFIG_MISC_EXCLUDE"`
	Debug           bool     `default:"false" envconfig:"SYNC_CONFIG_DEBUG"`
	DebugInclude    []string `envconfig:"SYNC_CONFIG_DEBUG_INCLUDE"`
	DebugExclude    []string `envconfig:"SYNC_CONFIG_DEBUG_EXCLUDE"`
}

func (raw *RawConfigSettings) Parse() *ConfigSettings {
	return &ConfigSettings{
		DNS:       NewConfigSetting(raw.DNS, raw.DNSInclude, raw.DNSExclude),
		DHCP:      NewConfigSetting(raw.DHCP, raw.DHCPInclude, raw.DHCPExclude),
		NTP:       NewConfigSetting(raw.NTP, raw.NTPExclude, raw.NTPExclude),
		Resolver:  NewConfigSetting(raw.Resolver, raw.ResolverExclude, raw.ResolverExclude),
		Database:  NewConfigSetting(raw.Database, raw.DatabaseExclude, raw.DatabaseExclude),
		Webserver: NewConfigSetting(raw.Webserver, nil, nil),
		Files:     NewConfigSetting(raw.Files, nil, nil),
		Misc:      NewConfigSetting(raw.Misc, raw.MiscExclude, raw.MiscExclude),
		Debug:     NewConfigSetting(raw.Debug, raw.DebugExclude, raw.DebugExclude),
	}
}

type ConfigSetting struct {
	Enabled bool
	Filter  *ConfigFilter
}

type ConfigFilter struct {
	Type filter.FilterType
	Keys []string
}

func newConfigFilter(filterType filter.FilterType, keys []string) *ConfigFilter {
	return &ConfigFilter{
		Type: filterType,
		Keys: keys,
	}
}

func NewConfigSetting(enabled bool, included, excluded []string) *ConfigSetting {
	var configFilter *ConfigFilter

	if included != nil {
		configFilter = newConfigFilter(filter.Include, included)
	} else if excluded != nil {
		configFilter = newConfigFilter(filter.Exclude, excluded)
	} else {
		configFilter = nil
	}

	return &ConfigSetting{
		Enabled: enabled,
		Filter:  configFilter,
	}
}

func (c *Config) Load() error {
	if err := envconfig.Process("", c); err != nil {
		return fmt.Errorf("env vars: %w", err)
	}

	if err := c.loadClient(); err != nil {
		return err
	}

	if err := c.loadSync(); err != nil {
		return err
	}

	return nil
}

func (sync *Sync) loadConfigSettings() error {
	cs := RawConfigSettings{}

	if err := envconfig.Process("", &cs); err != nil {
		return fmt.Errorf("config settings env vars: %w", err)
	}

	sync.ConfigSettings = cs.Parse()
	return nil
}

func (c *Config) String() string {
	replicas := make([]string, len(c.Replicas))
	for _, replica := range c.Replicas {
		replicas = append(replicas, replica.Url.String())
	}

	cron := ""
	if c.Sync.Cron != nil {
		cron = *c.Sync.Cron
	}

	sync := ""
	if c.Sync != nil {
		if mc := c.Sync.ConfigSettings; mc != nil {
			sync += fmt.Sprintf("config=%+v", *mc)
		}
		if gc := c.Sync.GravitySettings; gc != nil {
			sync += fmt.Sprintf(", gravity=%+v", *gc)
		}
	}

	return fmt.Sprintf("primary=%s, replicas=%s, fullSync=%t, cron=%s, sync=%s", c.Primary.Url, replicas, c.Sync.FullSync, cron, sync)
}
