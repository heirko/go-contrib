# Properties Module

[![GoDoc](https://godoc.org/github.com/heirko/go-contrib/properties?status.png)](https://godoc.org/github.com/heirko/go-contrib)

A centralized configuration endpoint.

## Why ?

We have a lot configuration possibilities with Golang, we need something open with large capabilities like Merge, subset, provider order, flags; and completly integrate with standard tools like Consul and etcd.

## How ?

To simplify, it's a [spf13/Viper](https://github.com/spf13/viper) extension and helper.

## Usages

Take a look at functional properties_test.go in properties_test and properties_test.go in properties unit test.

### Sample usage can be resume at :

To read Mode from command line and by default "prod", where resx is configuration directory store with app.json and prod.app.json .
app.json and prod.app.json (or any profile as command line argument) will be merge.

'''
props := properties.DefaultLoadModeProperties("./resx")
'''

To have more control and keep mode manager:
'''
props := properties.LoadModeProperties("./resx", "test", properties.GetSimpleProperties(), true)
'''

Can be use to unmarshal any configuration from [spf13/Viper](https://github.com/spf13/viper):

```golang

type DummyConfig struct {
		BaseUrl    string
		SuccessUrl string
	}

    var C DummyConfig

	err := props.Unmarshal(&C)
	
```


Add some functions, like GetOrDie, GetStringOrDefault, ...etc.
