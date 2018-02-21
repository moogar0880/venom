package venom_test

import (
	"flag"
	"fmt"
	"os"

	"github.com/moogar0880/venom"
)

func ExampleSetDefault() {
	venom.SetDefault("verbose", true)
	fmt.Println(venom.Get("verbose"))
	// Output: true
}

func ExampleSetOverride() {
	venom.SetDefault("verbose", true)
	venom.SetOverride("verbose", false)
	fmt.Println(venom.Get("verbose"))
	// Output: false
}

func ExampleEnvironmentVariableResolver_Resolve() {
	os.Setenv("LOG_LEVEL", "INFO")
	fmt.Println(venom.Get("log.level"))
	// Output: INFO
}

func ExampleFlagsetResolver_Resolve() {
	fs := flag.NewFlagSet("example", flag.ContinueOnError)
	fs.String("log-level", "WARNING", "set log level")

	flagResolver := &venom.FlagsetResolver{
		Flags:     fs,
		Arguments: []string{"-log-level=INFO"},
	}
	venom.RegisterResolver(venom.FlagLevel, flagResolver)
	fmt.Println(venom.Get("log.level"))
	// Output: INFO
}

func ExampleSetLevel() {
	var MySuperImportantLevel venom.ConfigLevel = venom.OverrideLevel + 1
	venom.SetLevel(MySuperImportantLevel, "verbose", true)
	venom.SetOverride("verbose", false)
	fmt.Println(venom.Get("verbose"))
	// Output: true
}

func ExampleFind() {
	key := "some.config"
	venom.SetDefault(key, 12)

	if val, ok := venom.Find(key); !ok {
		fmt.Printf("unable to find value for key %s", key)
	} else {
		fmt.Println(val)
	}
	// Output: 12
}

func ExampleGet() {
	venom.SetDefault("log.level", "INFO")
	fmt.Printf("%v\n", venom.Get("log"))
	fmt.Printf("%v\n", venom.Get("log.level")) // Output: INFO
	// Output: map[level:INFO]
	// INFO
}

func ExampleAlias() {
	venom.SetDefault("log.enabled", true)
	venom.Alias("verbose", "log.enabled")
	fmt.Println(venom.Get("verbose")) // Output: true
}
