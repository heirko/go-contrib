# Logrus Helper

A Helper arround [Logrus](https://github.com/Sirupsen/logrus) to wrap with [spf13/Viper](https://github.com/spf13/viper") to load configuration with fangs!

And to simplify [Logrus](https://github.com/Sirupsen/logrus) configuration use some behavior of [Logrus_mate](https://github.com/gogap/logrus_mate)

## Why?

[Logrus](https://github.com/Sirupsen/logrus) is wonderful but miss some configuration helper.
[Logrus_mate](https://github.com/gogap/logrus_mate) is powerful, but bring some unecessary pattern.
[spf13/Viper](https://github.com/spf13/viper")  is simple, powerful and generic.

## Howto use it

```go
import(

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/heirko/go-contrib/logrusHelper"
)


func initLogger() {

    // ########## Init Viper  
	var viper = viper.New()

	viper.SetConfigName("mate") // name of config file (without extension), here we use some logrus_mate sample
	viper.AddConfigPath("/etc/appname/")   // path to look for the config file in
	viper.AddConfigPath("$HOME/.appname")  // call multiple times to add many search paths
	viper.AddConfigPath(".")               // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
    // ########### End Init Viper

    // Read configuration
	var c = logrusHelper.UnmarshalConfiguration(viper) // Unmarshal configuration from Viper
	logrusHelper.SetConfig(logrus.StandardLogger(), c) // for e.g. apply it to logrus default instance
	
	// ### End Read Configuration
	
	// ### Use logrus as normal
	logrus.WithFields(logrus.Fields{
		"animal": "walrus",
	}).Error("A walrus appears")
}

```

