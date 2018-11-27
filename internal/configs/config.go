package configs

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

const privKeyEnvVar = "STAMP_PRIVATE_KEY"

type Config struct {
	Contract *ContractConfig `json:"contract"`
	Database *DatabaseConfig `json:"database"`
}

type ContractConfig struct {
	Address    string `json:"address"`
	Host       string `json:"host"`
	Port       string `json:"port"`
	Interval   int    `json:"interval"`
	PrivateKey string `json:"privateKey"`
	Password   string `json:"password"`
}

type DatabaseConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

func (nc ContractConfig) RawURL() string {
	return nc.Host + ":" + nc.Port
}

func ParseConfig(f io.Reader) (*Config, error) {
	config := new(Config)

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
