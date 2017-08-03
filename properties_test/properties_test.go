package propertiestest

import (
	"testing"

	"github.com/heirko/go-contrib/properties"
	"github.com/stretchr/testify/assert"
)

func TestModeLoadConfig(t *testing.T) {
	props := properties.LoadModeProperties("./resx", "test", properties.GetSimpleProperties(), true)

	assert.Equal(t, "http://tapp.test.me", props.GetString("app.plateform.baseUrl"))
	assert.Equal(t, "http://tapp.me", props.GetString("app.plateform.baseurlapp"))
	assert.Equal(t, 3, props.GetInt("app.plateform.val.t1"))
}

func TestModeLoadConfigPanic(t *testing.T) {
	assert.Panics(t, func() {
		properties.LoadModeProperties("./resx", "testNotExistMode", properties.GetSimpleProperties(), true)
	},
		"Mode not exists and should throw a panic")

	assert.Panics(t, func() {
		properties.LoadModeProperties("./resx", "", properties.GetSimpleProperties(), false)
	},
		"Mode not set and should throw a panic")

	assert.NotPanics(t, func() {
		properties.LoadModeProperties("./resx", "testNotExistMode", properties.GetSimpleProperties(), false)
	},
		"Mode not exists and should not throw a panic")

	wrongprops := properties.GetSimpleProperties()
	wrongprops.Set(properties.ConfigNameTag, "notexistconfigfilename")
	assert.Panics(t, func() {
		properties.LoadModeProperties("./resx", "test", wrongprops, false)
	},
		"Config file not exists and should throw a panic")
}

func TestModeLoadConfigBuggyFile(t *testing.T) {

	wrongprops := properties.GetSimpleProperties()
	wrongprops.Set(properties.ConfigNameTag, "appbuggy")
	assert.Panics(t, func() {
		properties.LoadModeProperties("./resx", "test", wrongprops, false)
	},
		"Config file is corrupted and should throw a panic")

	wrongprops = properties.GetSimpleProperties()
	assert.Panics(t, func() {
		properties.LoadModeProperties("./resx", "testbuggy", wrongprops, false)
	},
		"mode Config file is corrupted and should throw a panic")
}
