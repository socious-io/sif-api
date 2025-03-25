package config

import (
	"fmt"
	"os"

	"github.com/socious-io/goaccount"
	"github.com/socious-io/gopay"
	"gopkg.in/yaml.v2"
)

var Config *ConfigType

type ConfigType struct {
	Env      string `mapstructure:"env"`
	Port     int    `mapstructure:"port"`
	Debug    bool   `mapstructure:"debug"`
	Secret   string `mapstructure:"secret"`
	Host     string `mapstructure:"host"`
	Database struct {
		URL        string `mapstructure:"url"`
		SqlDir     string `mapstructure:"sqldir"`
		Migrations string `mapstructure:"migrations"`
	} `mapstructure:"database"`
	Cors struct {
		Origins []string `mapstructure:"origins"`
	} `mapstructure:"cors"`
	Nats struct {
		Url   string `mapstructure:"url"`
		Token string `mapstructure:"token"`
	} `mapstructure:"nats"`

	Payment struct {
		Chains gopay.Chains `mapstructure:"chains"`
		Fiats  gopay.Fiats  `mapstructure:"fiats"`
	} `mapstructure:"payment"`

	GoAccounts goaccount.Config `mapstructure:"goaccounts"`
	Upload     struct {
		Bucket      string `mapstructure:"bucket"`
		CDN         string `mapstructure:"cdn"`
		Credentials string `mapstructure:"credentials"`
	} `mapstructure:"upload"`

	AdminToken string `mapstructure:"admintoken"`

	Discord struct {
		Channel string `mapstructure:"channel"`
	} `mapstructure:"discord"`
}

func Init(filename string) (*ConfigType, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	conf := new(ConfigType)
	if err := decoder.Decode(conf); err != nil {
		return nil, err
	}
	Config = conf
	fmt.Println(Config.AdminToken, "---------------------------")
	return conf, err
}
