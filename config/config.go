package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/andscoop/ancli/card"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	// ConfigFileName is the name of the config files (home / project)
	ConfigFileName = "config"
	// ConfigFileType is the config file extension
	ConfigFileType = "yaml"
	// ConfigFileNameWithExt is the config filename with extension
	ConfigFileNameWithExt = fmt.Sprintf("%s.%s", ConfigFileName, ConfigFileType)
	// ConfigDir is the directory for astro files
	ConfigDir = ".ancli"

	// HomePath is the path to a users home directory
	HomePath, _ = homedir.Dir()
	// HomeConfigPath is the path to the users global config directory
	HomeConfigPath = filepath.Join(HomePath, ConfigDir)
	// HomeConfigFile is the global config file
	HomeConfigFile = filepath.Join(HomeConfigPath, ConfigFileNameWithExt)
)

// Init viper for config file in home directory
func Init() {
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType(ConfigFileType)
	viper.SetConfigFile(HomeConfigFile)

	viper.AutomaticEnv()

	viper.SetDefault("lastIndexed", "1900-01-01T00:00:00.00000-00:00")
	viper.SetDefault("deckPrefix", "#ancli")
	viper.SetDefault("cardFileExt", "md")
	// wasd default
	viper.SetDefault("cmdShortcuts.next", "d")
	viper.SetDefault("cmdShortcuts.back", "a")
	viper.SetDefault("cmdShortcuts.pass", "w")
	viper.SetDefault("cmdShortcuts.fail", "s")

	// Read in home config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok || strings.HasSuffix(err.Error(), "no such file or directory") {
			// No config file found.
			// That's okay we'll save when we save.
		} else {
			panic(err)
		}
	}
}

// GetConfig gets the viper object for direct manipulation
func GetConfig() *viper.Viper {
	return viper.GetViper()
}

// SetAndSave sets a string value in home config
func SetAndSave(path string, value interface{}) {
	viper.Set(path, value)
	err := Save()
	if err != nil {
		panic(err)
	}
}

// Save writes viper config to specified dir
func Save() error {
	if _, err := os.Stat(HomeConfigPath); os.IsNotExist(err) {
		os.Mkdir(HomeConfigPath, 0755)
	}

	err := viper.WriteConfigAs(HomeConfigFile)
	if err != nil {
		return err
	}
	return nil
}

// GetString will return config from home string
func GetString(path string) string {
	return viper.GetString(path)
}

// GetSavedCards will unmarshal an index object and return
func GetSavedCards() (map[string]card.Card, error) {
	var cs = make(map[string]card.Card)
	// todo probably not necessary to pass address of cs?
	err := viper.UnmarshalKey("decks", &cs)
	if err != nil {
		panic(err)
	}
	return cs, nil

}
