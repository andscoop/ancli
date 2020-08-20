package config

import (
	"fmt"
	"os"
	"path/filepath"

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

	CFG = cfgs{
		LastIndexed: newCfg("lastIndexed", "1900-01-01T00:00:00.00000-00:00"),
	}

	// CFGStrMap maintains string to cfg mapping
	CFGStrMap = make(map[string]cfg)
)

type cfg struct {
	Path    string
	Default string
}

// cfgs houses all configurations for an astro project
type cfgs struct {
	LastIndexed cfg
}

type Index struct {
	FilePath    string
	LastIndexed string
	LastTested  string
}

// Creates a new cfg struct
func newCfg(path string, dflt string) cfg {
	cfg := cfg{path, dflt}
	CFGStrMap[path] = cfg
	return cfg
}

// SetHomeString sets a string value in home config
func (c cfg) SetAndSave(value string) {
	viper.Set(c.Path, value)
	SaveConfig(viper.GetViper())
}

// saveConfig will save the config to a file
func SaveConfig(v *viper.Viper) error {
	err := v.WriteConfigAs(HomeConfigFile)
	if err != nil {
		return err
	}
	return nil
}

// GetConfig gets the viper object for direct manipulation
func GetConfig() *viper.Viper {
	return viper.GetViper()
}

// GetHomeString will return config from home string
func (c cfg) GetString() string {
	return viper.GetString(c.Path)
}

// GetIndex will unmarshal an index object and return
func GetIndex() (map[string]Index, error) {
	var i = make(map[string]Index)
	err := viper.UnmarshalKey("decks", &i)
	if err != nil {
		panic(err)
	}
	return i, nil

}

// Init viper for config file in home directory
func Init() {
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType(ConfigFileType)
	viper.SetConfigFile(HomeConfigFile)

	viper.AutomaticEnv()

	// If home config does not exist, create it
	homeConfigExists, _ := Exists(HomeConfigFile)
	if !homeConfigExists {
		err := os.Mkdir(HomeConfigPath, 0755)
		if err != nil {
			fmt.Println(err)
		}

		err = viper.SafeWriteConfigAs(HomeConfigFile)
		if err != nil {
			fmt.Println(err)
		}
	}

	for _, cfg := range CFGStrMap {
		if len(cfg.Default) > 0 {
			viper.SetDefault(cfg.Path, cfg.Default)
		}
	}

	// Read in home config
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("ERROR READING CONFIG")
		return
	}
}

// Exists returns a boolean indicating if the given path already exists
func Exists(path string) (bool, error) {
	if path == "" {
		return false, nil
	}

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if !os.IsNotExist(err) {
		return false, err
	}

	return false, nil
}
