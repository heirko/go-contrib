package logrusHelper_Test
// there is a bug in go < 1.9 with test only package and name, bug on namespace, we need to break recognition...

import (
	"testing"

	"fmt"
	"os"

	"github.com/heirko/go-contrib/logrusHelper"
	_ "github.com/heralight/logrus_mate/hooks/file"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func initLogger() {

	// ########## Init Viper
	var viper = viper.New()

	viper.SetConfigName("mate")           // name of config file (without extension), here we use some logrus_mate sample
	viper.AddConfigPath("/etc/appname/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.appname") // call multiple times to add many search paths
	viper.AddConfigPath(".")              // optionally look for config in the working directory
	err := viper.ReadInConfig()           // Find and read the config file
	if err != nil {                       // Handle errors reading the config file
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
	}).Error("A walrus appears in " + viper.GetString("formatter.name"))
}

func TestInitLogger(t *testing.T) {
	os.Remove("mate.log")

	Convey("Check Configuration", t, func(c C) {
		initLogger()
		time.Sleep(1 * time.Second)
		existLogFile, err := os.Stat("mate.log")
		notExist := os.IsNotExist(err)

		c.Convey(`"mate.log" should exists`, func(c C) {

			c.So(notExist, ShouldBeFalse)
			c.So(existLogFile, ShouldNotBeNil)
		})
		// Reset(func() {
		// 	// This reset is run after each `Convey` at the same scope.
		// 	err = os.Remove("mate.log")
		// 	if err != nil {
		// 		fmt.Println(err)
		// 		return
		// 	}
		// })
	})
}
