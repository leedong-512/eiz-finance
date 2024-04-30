package config
import (
	"github.com/spf13/viper"
	"sync"
)

var (
	configOnce sync.Once
	config     *Config
	err        error
)

type Parameters struct {
	Exp string
	Dir string
	Cpl string
	R1 string
	Sml string
	Sl string
	//Mysql    map[string]string
}
type Config struct {
	Parameters *Parameters
}

func newConfig() *Config {
	return &Config{
		Parameters: &Parameters{},
	}
}

func GetConfig() (*Config, error) {
	configOnce.Do(func() {
		config, err = prepareConfig()
	})
	return config, err
}

func Initialize(cfgFile string) {
	viper.SetConfigFile(cfgFile)
}

func prepareConfig() (*Config, error) {
	var (
		err error
	)

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := newConfig()

	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}