package config

import (
	"log"

	. "github.com/C3S/certinfo/internal/globals"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func InitViperConfig(cmd *cobra.Command) error {
	v := viper.New()

	customConfig := cmd.Flags().Lookup("config").Value.String()

	if customConfig != "" {
		v.SetConfigFile(customConfig)
	} else {
		v.SetConfigName("config.yaml")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("$HOME/.config/certinfo/")
	}
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Panicf("Fatal error in config file: %s\n", err)
		} else {
		}
		log.Print("Using fallback configuration (no config file found)")
	} else {
		err = v.UnmarshalKey("hosts", &AllHosts)
	}

	return nil
}
