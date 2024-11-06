package main

import (
	"vms/appconfig"
	"vms/internal/app"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @contact.name   giangmt@ivi.work
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	appconfig.GetConfigEnv()
	// appconfig.AppConfig()
	// Initiate a simple logger
	log := logrus.New()

	// Setup Configs
	cfg := viper.New()

	// Load Config
	cfg.AddConfigPath("./conf")
	cfg.SetEnvPrefix("app")
	cfg.AllowEmptyEnv(true)
	cfg.AutomaticEnv()
	err := cfg.ReadInConfig()
	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Warnf("No Config file found, loaded config from Environment - Default path ./conf")
		default:
			log.Fatalf("Error when Fetching Configuration - %s", err)
		}
	}

	// Load Config from Consul
	if cfg.GetBool("use_consul") {
		log.Infof("Setting up Consul Config source - %s/%s", cfg.GetString("consul_addr"), cfg.GetString("consul_keys_prefix"))
		err = cfg.AddRemoteProvider("consul", cfg.GetString("consul_addr"), cfg.GetString("consul_keys_prefix"))
		if err != nil {
			log.Fatalf("Error adding Consul as a remote Configuration Provider - %s", err)
		}

		cfg.SetConfigType("json")
		err = cfg.ReadRemoteConfig()
		if err != nil {
			log.Fatalf("Error when Fetching Configuration from Consul - %s", err)
		}

		if cfg.GetBool("from_consul") {
			log.Infof("Successfully loaded configuration from consul")
		}
	}

	// Run application
	log.Info("Start running vms service............................................")
	err = app.Run(cfg)
	if err != nil && err != app.ErrShutdown {
		log.Fatalf("Service stopped - %s", err)
	}
	log.Infof("Service shutdown - %s", err)
}
