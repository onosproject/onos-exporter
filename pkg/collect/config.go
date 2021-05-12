// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package collect

import (
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	configDir  = ".onos"
	addressKey = "service-address"

	tlsCertPathKey = "tls.certPath"
	tlsKeyPathKey  = "tls.keyPath"
	noTLSKey       = "no-tls"
	authHeaderKey  = "auth-header"
)

var configOptions = []string{
	addressKey,     // The gRPC endpoint
	tlsCertPathKey, // The path to the TLS certificate
	tlsKeyPathKey,  // The path to the TLS key
	noTLSKey,       // If present, do not use TLS
	authHeaderKey,  // Auth header in the form 'Bearer <base64>'
}

type config struct {
	subsystem string
	options   map[string]string
}

type Configuration interface {
	Init()
	Set(map[string]string) error
	getAddress() string
	getCertPath() string
	getKeyPath() string
	noTLS() bool
}

func NewConfig(subsystem string) Configuration {
	opts := make(map[string]string)
	for _, optName := range configOptions {
		opts[optName] = ""
	}

	return &config{
		subsystem: subsystem,
		options:   opts,
	}
}

func (c *config) Init() {
	for opt := range c.options {
		c.options[opt] = viper.GetString(opt)
	}
}

func (c *config) Set(options map[string]string) error {
	for opt, value := range options {
		if _, ok := c.options[opt]; ok {
			c.options[opt] = value
			viper.Set(opt, value)
		}
	}
	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}

func (c *config) getAddress() string {
	address := c.options[addressKey]
	if address == "" {
		return viper.GetString(addressKey)
	}
	return address
}

func (c *config) getCertPath() string {
	certPath := c.options[tlsCertPathKey]
	return certPath
}

func (c *config) getKeyPath() string {
	keyPath := c.options[tlsKeyPathKey]
	return keyPath
}

func (c *config) noTLS() bool {
	tls := c.options[noTLSKey]

	if tls == "" {
		return false
	} else {
		return true
	}
}

func runConfigInitCommand(configName string) error {
	if err := viper.ReadInConfig(); err == nil {
		return nil
	}

	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(home+"/"+configDir, 0777); err != nil {
		return err
	}

	filePath := home + "/" + configDir + "/" + configName + ".yaml"
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	_ = f.Close()

	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}

// InitConfig :
func InitConfig(configNameInit string) Configuration {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	viper.SetConfigName(configNameInit)
	viper.AddConfigPath(home + "/" + configDir)
	viper.AddConfigPath("/etc/onos")
	viper.AddConfigPath(".")

	err = runConfigInitCommand(configNameInit)
	if err != nil {
		panic(err)
	}

	_ = viper.ReadInConfig()

	config := NewConfig(configNameInit)
	config.Init()
	return config
}
