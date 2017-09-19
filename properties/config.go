package properties

const (
	// Default config type
	DefaultConfigType = "json"

	// Default config file name
	DefaultConfigName = "app"

	// Default config storage directory
	DefaultConfigDir = ""

	// Default config mode
	DefaultConfigMode = "prod"
	// Default test config mode spelling
	DefaultTestModeTag = "test"
)

// Flags Tag referrer
// e.g.
// ConfigNameTag imply => myexec --config-name "xxx"
const (
	ConfigNameTag = "config-name"
	ConfigDirTag  = "config-dir"
	ModeTag       = "mode"
	ConfigTypeTag = "config-type"
)

// Config is a struct that allows to initialize the Properties type with
// user defined values
type Config struct {

	// Define the config files type for K/V stores, allowed types are :
	// ["json", "toml", "yaml", "yml", "properties", "props", "prop"]
	// Default: json
	ConfigType string

	// Define to Environement variables to look up
	EnvVars []string

	// Define the name of config files (without extension) to look up for.
	// Default: "config"
	ConfigName string

	// Define the pathes where to lookup for config files
	ConfigPathes []string

	// Define the remote provides names:
	Providers []RemoteProvider

	// Define the flags to lookup for
	Flags []Flag

	// Overridable Mode Tag to use for test session by default set to DefaultTestModeTag
	TestModeTag string

	// Overridable default mode
	DefaultConfigMode string
}

func NewConfig() Config {
	return Config{
		ConfigType:        DefaultConfigType,
		ConfigName:        DefaultConfigName,
		TestModeTag:       DefaultTestModeTag,
		DefaultConfigMode: DefaultConfigMode,
	}
}

// DefaultConfig return a default configuration with Flags and Env already set...
func DefaultConfig() (c Config) {
	c = Config{
		EnvVars: []string{"HOME", "PWD"},
		Flags: []Flag{
			{ModeTag, "", "Execution mode: 'dev' or 'prod' or 'test'"},
			{ConfigDirTag, "", "Configuration directory"},
			{ConfigNameTag, "app", "Configuration name without extension"},
			{ConfigTypeTag, "json", "Configuration type, e.g.: json, yaml,..."},
		},
	}
	c.InitConfig()
	return
}

// InitConfig init config with default value if not set
func (c *Config) InitConfig() {
	if c.ConfigName == "" {
		c.ConfigName = DefaultConfigName
	}

	if c.ConfigType == "" {
		c.ConfigType = DefaultConfigType
	}
	if c.TestModeTag == "" {
		c.TestModeTag = DefaultTestModeTag
	}

	if c.DefaultConfigMode == "" {
		c.DefaultConfigMode = DefaultConfigMode
	}
	return
}

// Privider is a struct that hold remote providers data
type RemoteProvider struct {

	// Set the provider's name | must be "ectd" or "consul"
	Name string

	// Set the provider's url : "http://ip:port" for "etcd", "ip:port" for "consul"
	Url string

	// Set the path in the k/v store to retrieve configuration
	Path string

	// If set check the remote provider with encryption using the defined keyFile
	KeyFile string
}

// Flag is a struct that stores flags configuration
type Flag struct {

	// Flag Name used in command line
	Name string

	// When gives the default flag value
	Default string

	// Usage string shown in the help
	Usage string
}
