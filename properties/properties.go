package properties

import (
	"log"

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
		c = NewConfig()
	} else {
		c = config[0]
	}
	c.InitConfig()

	prop := Properties{Config: c, Viper: viper.New()}
	prop.init()

	return &prop
}

//initializes the properties instance - make calls to p.Viper library to initilize configuration.
func (p Properties) init() {

	//gives an instance of viper to Properties instance

	//Bind flags
	if p.Config.Flags != nil && len(p.Config.Flags) > 0 {
		for _, flag := range p.Config.Flags {
			pflag.String(flag.Name, flag.Default, flag.Usage)
			p.Viper.BindPFlag(flag.Name, pflag.Lookup(flag.Name))
		}
		pflag.Parse()
	}

	//Bind Env vars :
	if p.Config.EnvVars != nil && len(p.Config.EnvVars) > 0 {
		for _, envVar := range p.Config.EnvVars {
			p.Viper.BindEnv(envVar)
		}
	}

	//Set config file name
	var configName = p.GetStringOrDefault(ConfigNameTag, p.Config.ConfigName)
	var configType = p.GetStringOrDefault(ConfigTypeTag, p.Config.ConfigType)

	p.Viper.SetConfigName(configName)
	p.Viper.SetConfigType(configType)

	//set the lookup pathes for config files from overloading flags and env
	var configDir = p.GetStringOrDefault(ConfigDirTag, "")
	if configDir != "" {
		p.Config.ConfigPathes = []string{configDir}
	}

	if p.Config.ConfigPathes != nil && len(p.Config.ConfigPathes) > 0 {
		for _, path := range p.Config.ConfigPathes {
			p.Viper.AddConfigPath(path)
		}
		err := p.Viper.ReadInConfig()
		if err != nil {
			log.Panic(err)
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

// GetOrDie get key, if not found panic
func (p Properties) GetSubOrDie(key string) *viper.Viper {
	if v := p.Sub(key); v == nil {
		log.Panicf("Required sub properties \"%s\" are not found!!", key)
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
func (c Config) GetDefaultModeProperties() *Properties {
	props := New(c)
	return props
}

// CheckRunInTestEnvironment return true if this application is running with 'go test'
func CheckRunInTestEnvironment() bool {
	if pflag.Lookup("test.v") == nil {
		return false
	} else {
		return true
	}
}

// Helper to Load Properties and merge it with mode related Properties
// path will be use by default if user not provide a ConfigDirTag in command line
// defaultMode will be use by default if user not provide a ModeTag in command line
// props is used as properties base.
// panicOnModeLoad if true, when loading mode properties failed call "panic" otherwise "warning"
func (props *Properties) LoadModeProperties(panicOnModeLoad bool) *Properties {

	var configName = props.GetStringOrDefault(ConfigNameTag, props.Config.ConfigName)
	var configType = props.GetStringOrDefault(ConfigTypeTag, props.Config.ConfigType)

	props.SetConfigType(configType)
	var modeStr = props.GetString(ModeTag)
	if modeStr == "" {
		if isInTest := CheckRunInTestEnvironment(); isInTest == true {
			modeStr = props.Config.TestModeTag
		} else {
			modeStr = props.Config.DefaultConfigMode
		}
		if modeStr == "" {
			log.Panic("Mode is not set !")
		}
	}
	props.Set(ModeTag, modeStr)
	modeConfigName := modeStr + "." + configName
	props.SetConfigName(modeConfigName)

	err := props.MergeInConfig() // Find and read the config file
	if err != nil {
		if panicOnModeLoad {
			log.Panicf("Fatal error config mode %s : %s \n", modeConfigName, err)
		} else {
			log.Printf("Fatal error config mode %s : %s \n", modeConfigName, err)
			return props
		}
	}

	return props

}
