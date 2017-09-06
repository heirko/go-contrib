package properties

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	var jsonExample = []byte(`{
"id": "0001",
"user": { "name" : "donut"},
"name": "Cake",
"ppu": 0.55,
"amiauth": {
		"baseurl": "http://myapp.com",
		"successurl" : "http://myapp.com/private",
        "batter": [
                { "type": "Regular" },
                { "type": "Chocolate" },
                { "type": "Blueberry" },
                { "type": "Devil's Food" }
            ]
    }
}`)

	//instanciate the Properties handler
	props := New(Config{
		ConfigType: "json",
		EnvVars:    []string{"HOME", "PWD"},
		//ConfigPathes: []string{"."},
		//Flags: []Flag{
		//	{"mode", "prod", "Execution mode: 'dev' or 'prod'"},
		//},
	})

	r := bytes.NewReader(jsonExample)
	props.ReadConfig(r)

	//read env var
	assert.NotEmpty(t, props.GetString("HOME"))
	assert.NotEmpty(t, props.GetString("PWD"))

	//Get some prop values from config.json
	assert.Equal(t, "donut", props.GetString("user.name"))
}

func TestOverrideFlagFromConfig(t *testing.T) {
	var jsonExample = []byte(`{
"id": "0001",
"user": { "name" : "donut"},
"name": "Cake",
"mode": "dev",
"amiauth": {
		"baseurl": "http://myapp.com",
		"successurl" : "http://myapp.com/private",
        "batter": [
                { "type": "Regular" },
                { "type": "Chocolate" },
                { "type": "Blueberry" },
                { "type": "Devil's Food" }
            ]
    }
}`)

	//instanciate the Properties handler
	props := New(Config{
		ConfigType: "json",
		EnvVars:    []string{"HOME", "PWD"},
		//ConfigPathes: []string{"."},
		Flags: []Flag{
			{"mode", "prod", "Execution mode: 'dev' or 'prod'"},
		},
	})

	assert.Equal(t, "prod", props.GetString("mode"))

	r := bytes.NewReader(jsonExample)
	props.ReadConfig(r)
	assert.Equal(t, "dev", props.GetString("mode"))
}

func TestUnMarshallingSimple(t *testing.T) {
	type DummyConfig struct {
		BaseUrl    string
		SuccessUrl string
	}
	//instanciate the Properties handler
	props := New(Config{
		ConfigType: "json",
		EnvVars:    []string{"HOME", "PWD"},
		//ConfigPathes: []string{"."},

	})

	props.Set("BaseUrl", "http://myapp.com")
	props.Set("successurl", "http://myapp.com/private")

	var C DummyConfig

	err := props.Unmarshal(&C)
	if err != nil {
		t.Fatal("UnmarshalExact should error when populating a struct from a conf that contains unused fields")
	}
	assert.Equal(t, "http://myapp.com", C.BaseUrl)
	assert.Equal(t, "http://myapp.com/private", C.SuccessUrl)
}

func TestUnMarshallingSimpleNested(t *testing.T) {
	type DummyConfig struct {
		BaseUrl    string
		SuccessUrl string
	}

	//instanciate the Properties handler
	props := New(Config{
		ConfigType: "json",
	})

	props.Set("amiauth", map[string]interface{}{
		"BaseUrl":    "http://myapp.com",
		"successUrl": "http://myapp.com/private",
	})
	var C DummyConfig

	err := props.UnmarshalKey("amiauth", &C)
	if err != nil {
		t.Fatal("UnmarshalExact should error when populating a struct from a conf that contains unused fields")
	}
	//fmt.Print(props.AllSettings())
	assert.Equal(t, "http://myapp.com", props.Get("amiauth.baseurl"))
	assert.Equal(t, "http://myapp.com", C.BaseUrl)
	assert.Equal(t, "http://myapp.com/private", C.SuccessUrl)
}

func TestReadSubBugReadingUserName(t *testing.T) {
	var jsonSrc = []byte(`{
"id": "0002",
"user": { "name" : "donut"},
"name": "Cake",
"ppu": 0.55,
"amiauth": {
		"baseurl": "http://myapp.com",
		"successurl" : "http://myapp.com/private",
        "batter": [
                { "type": "Regular" },
                { "type": "Chocolate" },
                { "type": "Blueberry" },
                { "type": "Devil's Food" }
            ]
    }
}`)

	var r = bytes.NewReader(jsonSrc)
	var src = viper.New()
	src.SetConfigType("json")
	var err = src.ReadConfig(r)
	assert.Empty(t, err)
	assert.Equal(t, "donut", src.GetString("user.name"))
	assert.Equal(t, "0002", src.GetString("id"))
	assert.Equal(t, 0.55, src.Get("ppu"))
	assert.Equal(t, "http://myapp.com", src.GetString("amiauth.baseurl"))

}

func TestReadSubBug(t *testing.T) {
	var jsonSrc = []byte(`{
"id": "0001",
"user": {
 "name" : "donut"
},
"name": "Cake",
"plus": "efefe",
"amiauth": {
		"baseurl": "toto",
		"eee": "ttt",
        "batter": [
                { "type": "Regular" },
                { "type": "Chocolate" },
                { "type": "Blueberry" },
                { "type": "Devil's Food" }
            ]
    }
}`)

	var r = bytes.NewReader(jsonSrc)
	var src = viper.New()
	src.SetConfigType("json")
	var err = src.ReadConfig(r)
	assert.Empty(t, err)
	assert.Equal(t, "donut", src.GetString("user.name"))
	assert.Equal(t, "0001", src.GetString("id"))
	assert.Equal(t, "efefe", src.Get("plus"))
	assert.Equal(t, "toto", src.GetString("amiauth.baseurl"))

}

func TestMerge(t *testing.T) {
	var jsonDst = []byte(`{
"id": "0001",
"user": {
 "name" : "donut"
},
"name": "Cake",
"plus": "efefe",
"amiauth": {
		"baseurl": "toto",
		"eee": "ttt",
        "batter": [
                { "type": "Regular" },
                { "type": "Chocolate" },
                { "type": "Blueberry" },
                { "type": "Devil's Food" }
            ]
    }
}`)
	var jsonSrc = []byte(`{
"id": "0002",
"user": { "name" : "donut"},
"name": "Cake",
"ppu": 0.55,
"amiauth": {
		"baseurl": "http://myapp.com",
		"successurl" : "http://myapp.com/private",
        "batter": [
                { "type": "Regular" },
                { "type": "Chocolate" },
                { "type": "Blueberry" },
                { "type": "Devil's Food" }
            ]
    }
}`)

	var r = bytes.NewReader(jsonDst)
	var dst = viper.New()
	dst.SetConfigType("json")
	var err = dst.ReadConfig(r)
	assert.Empty(t, err)
	assert.Equal(t, "donut", dst.GetString("user.name"))
	assert.Equal(t, "0001", dst.GetString("id"))
	assert.Equal(t, "efefe", dst.GetString("plus"))
	assert.Equal(t, "ttt", dst.GetString("amiauth.eee"))

	r = bytes.NewReader(jsonSrc)
	var src = viper.New()
	src.SetConfigType("json")
	err = src.ReadConfig(r)
	assert.Empty(t, err)

	r = bytes.NewReader(jsonSrc)
	err = dst.MergeConfig(r)
	assert.Empty(t, err)
	assert.Equal(t, "0002", dst.GetString("id"))
	assert.Equal(t, "donut", dst.GetString("user.name"))
	assert.Equal(t, 0.55, dst.GetFloat64("ppu"))
	assert.Equal(t, "efefe", dst.GetString("plus"))
	assert.Equal(t, "http://myapp.com", dst.GetString("amiauth.baseurl"))
	assert.Equal(t, "ttt", dst.GetString("amiauth.eee"))
}

func TestGetOrDie(t *testing.T) {
	var jsonSrc = []byte(`{
"id": "0001",
"user": {
 "name" : "donut"
},
"name": "Cake",
"plus": "efefe",
"amiauth": {
		"baseurl": "toto",
		"eee": "ttt",
        "batter": [
                { "type": "Regular" },
                { "type": "Chocolate" },
                { "type": "Blueberry" },
                { "type": "Devil's Food" }
            ]
    }
}`)

	var r = bytes.NewReader(jsonSrc)
	var src = New()
	src.SetConfigType("json")
	var err = src.ReadConfig(r)
	assert.Empty(t, err)
	assert.Equal(t, "donut", src.GetOrDie("user.name"))
	assert.Panics(t, func() {
		src.GetOrDie("user.notExist")
	}, "Calling GetOrDie() should panic")

}