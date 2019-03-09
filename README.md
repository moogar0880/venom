# venom

A pluggable hierarchical configuration management library with zero
dependencies, for golang.

[![Build Status](https://travis-ci.org/moogar0880/venom.svg?branch=master)](https://travis-ci.org/moogar0880/venom)
[![Go Report Card](https://goreportcard.com/badge/github.com/moogar0880/venom)](https://goreportcard.com/report/github.com/moogar0880/venom)
[![GoDoc](https://godoc.org/github.com/moogar0880/venom?status.svg)](https://godoc.org/github.com/moogar0880/venom)

## Install

With dep:

`dep ensure -add github.com/moogar0880/venom`

Or with go get:

`go get github.com/moogar0880/venom`

## Why Venom?

This library aims to provide the basic building blocks of a configuration
management library. It exposes as many aspects as the standard lib will
reasonably support, such as:

1. Explicitly set defaults or overrides.
1. Find and load JSON config files.
1. Load configuration files from a directory, with the option to recursively
   load configs from sub-directories.
1. Load configs from environment variables.
1. Load command line arguments using the `flags` library.

These mechanisms are exposed in a completely extendible way which will allow 
you to easily perform the following actions:

1. Specify your own config level precedence (eg, environment variables can be
   made to override command line arguments).
1. Default behavior can be replaced with custom behavior by writing a custom
   `Resolver`.
1. Custom config levels can be easily implemented by writing a custom `Resolver`.

## Loading Configs

### Setting Defaults

You can easily set default values for configuration options that are not 
required to be specified by config levels with a higher precedence.

```go
// simple values can easily be set as defaults
venom.SetDefault("verbose", false)
venom.SetDefault("hostname", "localhost")

// more complex objects can also be set
venom.SetDefault("log", map[string]interface{
    "level": "INFO",
    "fields": map[string]interface{
        "origin": "my_awesome_app", 
        "category": "categories",
    },
})
```

### Loading Configs From Files

Venom allows you to specify custom file loaders for specific file types. By
default the only implementation is the `JSONLoader` which loads any file with
an extension of `.json` using `json.Unmarshal`.

If you wish to implement your own type of config file reader you need only to
implement the `IOFileLoader` interface:

```go
type IOFileLoader func(io.Reader) (map[string]interface{}, error)
```

Once you have any custom `IOFileLoader`s registered via `RegisterExtension`,
you can easily load a single file or all files in a directory.

```go
// load a single file
venom.LoadFile("config.json")

// load all files within a directory
venom.LoadDirectory("/etc/conf.d", false)

// recursively load all files within a directory or it's sub-directories
venom.LoadDirectory("/etc/conf.d", true)
```

### Setting Overrides

You can easily set values which overrides all other values for a single 
configuration.

```go
venom.SetOverride("verbose", true)
```

### Environment Variables

Configuration values may be loaded from any set environment variables, which
is enabled by default in the global venom instance and when creating a custom
venom instance via `venom.Default()`. 

**Note:** Environment variables are case sensitive and should be set as 
uppercase strings with words that are separated by an underscore ("_") 
character.

```go
os.Setenv("LOG_LEVEL", "INFO")

fmt.Println(venom.Get("log.level"))  // Output: "INFO"
```

To specify a custom environment variable prefix, you can simply create and 
register a new `EnvironmentVariableResolver`.

```go
envVarResolver := &venom.EnvironmentVariableResolver{
    Prefix: "MYSERVICE",
}

venom.RegisterResolver(venom.EnvironmentLevel, envVarResolver)

os.Setenv("MYSERVICE_LOG_LEVEL", "INFO")

fmt.Println(venom.Get("log.level"))  // Output: "INFO"
```

### Flags

By default commandline flags can be parsed using the standard lib `flag` 
library. There are two different ways that you can approach loading command 
line configs, and they both involve creating and registering a new 
`FlagsetResolver` instance with Venom.

**Note:** Venom will parse flags if they have not already been parsed, but 
Venom will not attempt to re-parse an already parsed `FlagSet`.

#### Straight From os.Args

If you configure your commandline options to be loaded from the global `flags`
FlagSet, you can pass a zero-value `FlagsetResolver` when registering the
resolver with Venom.

```go
fs.String("log-level", "WARNING", "set log level")

flagResolver := &venom.FlagsetResolver{}

venom.RegisterResolver(venom.FlagLevel, flagResolver)
```

#### Provide a Custom `flag.FlagSet`

You also have the option to provide a custom `flag.FlagSet` instance via a
`FlagsetResolver`.

**Note:** When providing a `FlagSet`, you may also provide the arguments to be
parsed if the `FlagSet` has not already been parsed. The arguments will default
to `os.Args[1:]` if none are specified. 

```go
fs := flag.NewFlagSet("example", flag.ContinueOnError)
fs.String("log-level", "WARNING", "set log level")

flagResolver := &venom.FlagsetResolver{
    Flags: fs,
}
venom.RegisterResolver(venom.FlagLevel, flagResolver)
```

### Custom ConfigLevels

Consuming applications are able to define their own `ConfigLevel`s in order to
define configuration values with higher or lower precedence. For example, to
set configs at levels above `Override` (not recommended), you could do 
something similar to the following:

```go
const MySuperImportantLevel venom.ConfigLevel = venom.OverrideLevel + 1

venom.SetLevel(MySuperImportantLevel, "verbose", true)
```

## Reading Config Values

There are several ways to access config values from Venom:

* `Get(key string) interface{}`
* `Find(key string) (interface{}, bool)`

If you're unsure whether or not a config value has been set, `Find` will return
an optional boolean value indicating whether or not the value has been 
specified. Otherwise, `Get` will return `nil` in the event that a config has not
been specified.

## Key Management

Venom automatically nests config values that are specified as separated by the
value of `Delim`, which defaults to `"."`. This means that you can express
more complex config structures when setting and reading variables.

```go
venom.SetDefault("log.level", "INFO")
fmt.Printf("%v", venom.Get("log"))  // Output: map[string]interface{"level": "INFO"}
fmt.Printf("%v", venom.Get("log.level"))  // Output: "INFO"
```

## Aliasing Keys

Venom exposes the ability to alias one key to another. This allows applications
to more easily modify their configuration without breaking backwards 
compatibility when doing so.

```go
venom.SetDefault("log.enabled", true)
venom.Alias("verbose", "log.enabled")
fmt.Println(venom.Get("verbose"))  // Output: true
```

## Unmarshal Configs

Venom supports the ability to unmarshal configuration data into struct values
and also introduces the new `venom` struct field tag, which allows callers
to specify different configuration fields to set to fields. Nested structs

**Note:** that errors are returned if the stored type and the type of the struct 
field do not match.

**Note:** Nested structs are supported and will carry the context of the name
of the parent config. See the following example for more on this

```go
type LoggingConfig struct {
    // this field will be loaded from log.level if Unmarshalled via the 
    // following Config struct
    Level string `venom:"level"`
}

type Config struct {
    Log LoggingConfig `venom:"log"` 
}

var c Config
err := venom.Unmarshal(nil, &c)
```

## Safety

As of 0.2.0, venom now exposes optional functions for defining Venom instances 
that are safe for concurrent goroutine access.

Two functions are exposed to provide new goroutine safe Venom instances, the 
behavior of which mimics the functions unsafe counterparts:

```go
// generate a new, empty, venom instance that is safe for concurrent access 
ven := venom.NewSafe()

// or generate a new venom instance with some default levels applied to it
ven = venom.DefaultSafe()
```

## Custom Venom Behavior

The above examples show how to use the global venom instance for the sake of 
brevity, but you can create your own venom instances to use directly.

```go
ven := venom.New()

ven.SetDefault("verbose", false)
```

Additionally, as of 0.2.0, you can define your own `ConfigStore` to control how
venom manages it's underlying configuration storage. This can be achieved by 
creating a new Venom instance with a predefined `ConfigStore` via the
`NewWithStore` function. 

```go
ven := venom.NewWithStore(venom.NewSafeConfigStore())

ven.SetDefault("verbose", false)
````

## Benchmarks

```
goos: darwin
goarch: amd64
pkg: github.com/moogar0880/venom
BenchmarkVenomGet/single_ConfigLevel_with_one_key/value_pair-8         	20000000	        84.6 ns/op	      16 B/op	       1 allocs/op
BenchmarkVenomGet/many_key/value_pairs_in_a_single_ConfigLevel-8       	20000000	       102 ns/op	      16 B/op	       1 allocs/op
BenchmarkVenomGet/many_key/value_pairs_spread_across_multiple_ConfigLevels-8         	20000000	        99.2 ns/op	      16 B/op	       1 allocs/op
BenchmarkVenomWrite/single_key/value_pair_in_one_ConfigLevel-8                       	20000000	        89.0 ns/op	      16 B/op	       1 allocs/op
BenchmarkVenomWrite/many_key/value_pairs_in_one_ConfigLevel-8                        	 2000000	       648 ns/op	     121 B/op	       3 allocs/op
BenchmarkVenomWrite/many_nested_key/value_pairs_in_one_ConfigLevel-8                 	 1000000	      2055 ns/op	     922 B/op	       7 allocs/op
BenchmarkVenomWrite/many_key/value_pairs_in_many_ConfigLevels-8                      	 2000000	       589 ns/op	     121 B/op	       3 allocs/op
BenchmarkVenomWrite/many_nested_key/value_pairs_in_many_ConfigLevels-8               	 1000000	      1679 ns/op	     922 B/op	       7 allocs/op
```
