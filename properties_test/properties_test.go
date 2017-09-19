package propertiestest

import (
	"testing"

	"github.com/heirko/go-contrib/properties"
	"github.com/stretchr/testify/assert"
)

func TestModeLoadConfig(t *testing.T) {
	c := properties.NewConfig()
	c.ConfigPathes = []string{"./resx"}
	c.DefaultConfigMode = "test"

	props := properties.New(c)
	props.LoadModeProperties(true)

	assert.Equal(t, "http://tapp.test.me", props.GetString("app.plateform.baseUrl"))
	assert.Equal(t, "http://tapp.me", props.GetString("app.plateform.baseurlapp"))
	assert.Equal(t, 3, props.GetInt("app.plateform.val.t1"))
}

func TestModeLoadConfigPanic(t *testing.T) {
	assert.Panics(t, func() {
		c := properties.NewConfig()
		c.ConfigPathes = []string{"./resx"}
		c.DefaultConfigMode = "testNotExistMode"

		properties.New(c).LoadModeProperties(true)
	},
		"Mode not exists and should throw a panic")

	assert.Panics(t, func() {
		c := properties.NewConfig()
		c.ConfigPathes = []string{"./resx"}
		c.DefaultConfigMode = ""
		properties.New(c).LoadModeProperties(true)
	},
		"Mode not set and should throw a panic")

	assert.NotPanics(t, func() {
		c := properties.NewConfig()
		c.ConfigPathes = []string{"./resx"}
		c.DefaultConfigMode = "testNotExistMode"

		properties.New(c).LoadModeProperties(false)
	},
		"Mode not exists and should not throw a panic")

	assert.Panics(t, func() {
		c := properties.NewConfig()
		c.ConfigPathes = []string{"./resx"}
		c.DefaultConfigMode = "test"
		props := properties.New(c)
		props.Set(properties.ConfigNameTag, "notexistconfigfilename")
		props.LoadModeProperties(true)
	},
		"Config file not exists and should throw a panic")
}

func TestModeLoadConfigBuggyFile(t *testing.T) {

	assert.Panics(t, func() {
		c := properties.NewConfig()
		c.ConfigPathes = []string{"./resx"}
		c.DefaultConfigMode = "test"
		c.ConfigName = "appbuggy"
		props := properties.New(c)
		props.LoadModeProperties(false)
	},
		"Config file is corrupted and should throw a panic")

	assert.Panics(t, func() {
		c := properties.NewConfig()
		c.ConfigPathes = []string{"./resx"}
		c.DefaultConfigMode = "test"
		c.DefaultConfigMode = "testNotExistMode"

		props := properties.New(c)
		props.LoadModeProperties(true)
	},
		"mode Config file is corrupted and should throw a panic")

	assert.NotPanics(t, func() {
		c := properties.NewConfig()
		c.ConfigPathes = []string{"./resx"}
		c.DefaultConfigMode = "test"
		c.DefaultConfigMode = "testNotExistMode"

		props := properties.New(c)
		props.LoadModeProperties(false)
	},
		"mode Config file is corrupted and should throw a panic")
}
