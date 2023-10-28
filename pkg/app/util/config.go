package util

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	supererrors "github.com/sarumaj/go-super/errors"
	viper "github.com/spf13/viper"
	keyring "github.com/zalando/go-keyring"
)

const configFileName = "ldap-cli.conf"
const keyringServiceName = "ldap-cli"

var config = sync.Pool{New: func() any {
	confCfg := viper.New()
	confCfg.SetConfigFile(configFileName)
	confCfg.SetConfigType("yaml")
	confCfg.AddConfigPath(ConfigDir())
	supererrors.Except(confCfg.ReadRemoteConfig(), viper.ConfigFileNotFoundError{})

	var conf Configuration
	supererrors.Except(confCfg.Unmarshal(&conf))
	conf.Viper = confCfg

	if conf.User != "" {
		conf.Password = supererrors.ExceptFn(supererrors.W(keyring.Get(keyringServiceName, conf.User)), keyring.ErrNotFound)
	}

	return &conf
}}

type Configuration struct {
	*viper.Viper `yaml:"-"`
	Password     string `yaml:"-"`
	User         string
	Timeout      time.Duration
}

func Config() *Configuration {
	return config.Get().(*Configuration)
}

// Config path precedence: XDG_CONFIG_HOME, AppData (windows only), HOME.
func ConfigDir() string {
	var path string
	if b := os.Getenv(xdgConfigHome); b != "" {
		path = filepath.Join(b, "gh")
	} else if c := os.Getenv(appData); runtime.GOOS == "windows" && c != "" {
		path = filepath.Join(c, "ldap-cli")
	} else {
		d, _ := os.UserHomeDir()
		path = filepath.Join(d, ".config", "ldap-cli")
	}
	return path
}
