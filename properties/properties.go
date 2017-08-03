package properties

import (
	"log"

	"bytes"

	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

// Properties is a struct that allow to handle application
// properties seting and acessing. It wraps spf13/viper library
type Properties struct {
	*viper.Viper
	Config Config
}

// Properties constructor
// Settings default values if need
func New(config ...Config) *Properties {
	var c Config

	if len(config) == 0 {
		c = Config{}
	} else {
		c = config[0]
	}

	if c.ConfigName == "" {
		c.ConfigName = DefaultConfigName
	}

	if c.ConfigType == "" {
		c.ConfigType = DefaultConfigType
	}

	prop := Properties{Config: c, Viper: viper.New()}
	prop.init()

	return &prop
}

//initializes the properties instance - make calls to p.Viper library to initilize configuration.
func (p Properties) init() {

	//gives an instance of viper to Properties instance

	//Set config file name

	if p.Config.ConfigName != "" {
		p.Viper.SetConfigName(p.Config.ConfigName)
	}

	//set config file type for remote K/V stores
	if p.Config.ConfigType != "" {
		p.Viper.SetConfigType(p.Config.ConfigType)
	}

	//set the lookup pathes for config files
	if p.Config.ConfigPathes != nil && len(p.Config.ConfigPathes) > 0 {
		for _, path := range p.Config.ConfigPathes {
			p.Viper.AddConfigPath(path)
		}
		err := p.Viper.ReadInConfig()
		if err != nil {
			log.Panic(err)
		}
	}

	//Bind Env vars :
	if p.Config.EnvVars != nil && len(p.Config.EnvVars) > 0 {
		for _, envVar := range p.Config.EnvVars {
			p.Viper.BindEnv(envVar)
		}
	}

	//Set remote providers
	if p.Config.Providers != nil && len(p.Config.Providers) > 0 {
		for _, provider := range p.Config.Providers {
			if provider.KeyFile != "" {
				p.Viper.AddSecureRemoteProvider(provider.Name, provider.Url, provider.Path, provider.KeyFile)
			} else {
				p.Viper.AddRemoteProvider(provider.Name, provider.Url, provider.Path)
			}
		}
		err := p.Viper.ReadRemoteConfig()
		if err == nil {
			log.Panic(err)
		}
	}

	//Bind flags
	if p.Config.Flags != nil && len(p.Config.Flags) > 0 {
		for _, flag := range p.Config.Flags {
			pflag.String(flag.Name, flag.Default, flag.Usage)
			p.Viper.BindPFlag(flag.Name, pflag.Lookup(flag.Name))
		}
		pflag.Parse()
	}
}

// GetOrDie get key, if not found panic
func (p Properties) GetOrDie(key string) interface{} {
	if v := p.Get(key); v == nil {
		log.Panicf("Required property %s is not found!!", key)
	} else {
		return v
	}
	return nil
}

// GetStringOrDefault get string or a default value
func (p Properties) GetStringOrDefault(key string, dlft string) string {
	if v := p.GetString(key); v == "" {
		return dlft
	} else {
		return v
	}
}

// TryLoadRemoteProperties try load configuration from remote througth Viper
func (p Properties) TryLoadRemoteProperties() {
	var name = p.GetString("remote.name")
	var url = p.GetString("remote.url")
	var path = p.GetString("remote.path")
	var key = p.GetString("remote.key")
	if name != "" && url != "" && path != "" {
		if key != "" {
			p.AddSecureRemoteProvider(name, url, path, key)
		} else {
			p.AddRemoteProvider(name, url, path)
		}

		err := p.Viper.ReadRemoteConfig()
		if err != nil {
			// Handle errors reading the config file
			log.Panicf("Fatal error config Remote provider: %s \n", err)
		}
	}

}

// GetDefaultModeProperties get A Default Property set for classic app based on Flag with
// Config file : app.json in current directory from Flag
// With HOME and PWD from env
func GetDefaultModeProperties() *Properties {
	props := New(Config{
		EnvVars: []string{"HOME", "PWD"},
		Flags: []Flag{
			{ModeTag, "", "Execution mode: 'dev' or 'prod' or 'test'"},
			{ConfigDirTag, "", "Configuration directory"},
			{ConfigNameTag, "app", "Configuration name without extension"},
			{ConfigTypeTag, "json", "Configuration type, e.g.: json, yaml,..."},
		},
	})
	return props
}

// A Base Property set for quick init with
// Config file : app.json in current directory from Properties (no flag listener)
func GetSimpleProperties() *Properties {
	props := New(Config{})
	props.Set(ConfigDirTag, DefaultConfigDir)
	props.Set(ConfigNameTag, DefaultConfigName)
	props.Set(ConfigTypeTag, DefaultConfigType)
	return props
}

// Helper to Load Properties and merge it with mode related Properties
// defaultPath will be use by default if user not provide a ConfigDirTag in command line
func DefaultLoadModeProperties(defaultPath string) *Properties {
	return LoadModeProperties(defaultPath, DefaultConfigMode, GetDefaultModeProperties(), true)
}

// Helper to Load Properties and merge it with mode related Properties
// path will be use by default if user not provide a ConfigDirTag in command line
// defaultMode will be use by default if user not provide a ModeTag in command line
// props is used as properties base.
// panicOnModeLoad if true, when loading mode properties failed call "panic" otherwise "warning"
func LoadModeProperties(defaultPath string, defaultMode string, props *Properties, panicOnModeLoad bool) *Properties {

	var configName = props.GetStringOrDefault(ConfigNameTag, DefaultConfigName)
	var configDir = props.GetStringOrDefault(ConfigDirTag, defaultPath)
	var configType = props.GetStringOrDefault(ConfigTypeTag, DefaultConfigType)

	props.SetConfigType(configType)
	props.AddConfigPath(configDir)
	props.SetConfigName(configName)
	err := props.ReadInConfig() // Find and read the config file
	if err != nil {
		// Handle errors reading the config file
		log.Panicf("Fatal error config file %s.%s in %s : %s \n", configName, configType, configDir, err)
	}

	var modeStr = props.GetStringOrDefault(ModeTag, defaultMode)
	if modeStr == "" {
		log.Panic("Mode is not set !")
	}

	var modeConfigFilePath = configDir + "/" + modeStr + "." + configName + "." + configType
	file, err := afero.ReadFile(afero.NewOsFs(), modeConfigFilePath)
	if err != nil {
		if panicOnModeLoad {
			log.Panicf("Fatal error reading mode file %s : %s \n", modeConfigFilePath, err)
		} else {
			log.Printf("Error reading mode file %s : %s \n", modeConfigFilePath, err)
			return props
		}
	}

	err = props.MergeConfig(bytes.NewReader(file))
	if err != nil {
		log.Panicf("Fatal error merging mode file %s : %s \n", modeConfigFilePath, err)
	}
	return props

}
