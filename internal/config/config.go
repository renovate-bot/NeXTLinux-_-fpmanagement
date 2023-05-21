package config

import (
	"fmt"
	"path"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/adrg/xdg"
	"github.com/dakaneye/fpmanagement/internal"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const redacted = "******"

// Application is the main nextlinuxctl application configuration.
type Application struct {
	Nextlinux           Nextlinux          `yaml:"nextlinux" mapstructure:"nextlinux"`                     // Nextlinux provides connection details for requests to Nextlinux
}

// Information for posting in-use image details to Nextlinux (or any URL for that matter)
type Nextlinux struct {
	URL      string     `yaml:"url" mapstructure:"url"`
	User     string     `yaml:"user" mapstructure:"user"`
	Password string     `yaml:"password" mapstructure:"password"`
	Account  string     `yaml:"account" mapstructure:"account"`
	HTTP     HTTPConfig `yaml:"http" mapstructure:"http"`
}

// Configurations for the HTTP Client itself (net/http)
type HTTPConfig struct {
	Insecure       bool `yaml:"insecure" mapstructure:"insecure"`
	TimeoutSeconds int  `yaml:"timeoutSeconds" mapstructure:"timeoutSeconds"`
}

// LoadApplicationConfig populates the given viper object with application configuration discovered on disk
func LoadApplicationConfig(v *viper.Viper) (*Application, error) {
	// the user may not have a config, and this is OK, we can use the default config + default cobra cli values instead
	setNonCliDefaultValues(v)
	_ = readConfig(v, "")

	config := &Application{}
	err := v.Unmarshal(config)
	if err != nil {
		return nil, fmt.Errorf("unable to parse config: %w", err)
	}

	return config, nil
}

// readConfig attempts to read the given config path from disk or discover an alternate store location
func readConfig(v *viper.Viper, overridePath string) error {
	v.AutomaticEnv()
	v.SetEnvPrefix(internal.ApplicationName)
	// allow for nested options to be specified via environment variables
	// e.g. pod.context = APPNAME_POD_CONTEXT
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// use explicitly the given user config
	if overridePath != "" {
		v.SetConfigFile(overridePath)
		if err := v.ReadInConfig(); err == nil {
			return nil
		}
		// don't fall through to other options if this fails
		return fmt.Errorf("unable to read config: %v", overridePath)
	}

	// start searching for valid configs in order...

	// 1. look for .<appname>.yaml (in the current directory)
	v.AddConfigPath(".")
	v.SetConfigName(internal.ApplicationName)
	if err := v.ReadInConfig(); err == nil {
		return nil
	}

	// 2. look for .<appname>/config.yaml (in the current directory)
	v.AddConfigPath("." + internal.ApplicationName)
	v.SetConfigName("config")
	if err := v.ReadInConfig(); err == nil {
		return nil
	}

	// 3. look for ~/.<appname>.yaml
	home, err := homedir.Dir()
	if err == nil {
		v.AddConfigPath(home)
		v.SetConfigName("." + internal.ApplicationName)
		if err := v.ReadInConfig(); err == nil {
			return nil
		}
	}

	// 4. look for <appname>/config.yaml in xdg locations (starting with xdg home config dir, then moving upwards)
	v.AddConfigPath(path.Join(xdg.ConfigHome, internal.ApplicationName))
	for _, dir := range xdg.ConfigDirs {
		v.AddConfigPath(path.Join(dir, internal.ApplicationName))
	}
	v.SetConfigName("config")
	if err := v.ReadInConfig(); err == nil {
		return nil
	}

	return fmt.Errorf("application config not found")
}

// setNonCliDefaultValues ensures that there are sane defaults for values that do not have CLI equivalent options (where there would already be a default value)
func setNonCliDefaultValues(v *viper.Viper) {
	v.SetDefault("nextlinux.account", "admin")
	v.SetDefault("nextlinux.http.insecure", false)
	v.SetDefault("nextlinux.http.timeoutSeconds", 10)
}

func (cfg Application) String() string {
	// redact sensitive information
	// Note: If the configuration grows to have more redacted fields it would be good to refactor this into something that
	// is more dynamic based on a property or list of "sensitive" fields
	if cfg.Nextlinux.Password != "" {
		cfg.Nextlinux.Password = redacted
	}

	// yaml is pretty human friendly (at least when compared to json)
	appCfgStr, err := yaml.Marshal(&cfg)

	if err != nil {
		return err.Error()
	}

	return string(appCfgStr)
}
