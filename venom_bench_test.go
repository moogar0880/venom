package venom

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

var readResult interface{}

func benchmarkVenomRead(ven *Venom, key string, b *testing.B) {
	var v interface{}
	for n := 0; n < b.N; n++ {
		v = ven.Get(key)
	}
	readResult = v
}

type keyer func(i int) string
type leveler func(i int) ConfigLevel

func benchmarkVenomWrite(ven *Venom, l leveler, key keyer, val interface{}, b *testing.B) {
	for n := 0; n < b.N; n++ {
		ven.SetLevel(l(n), key(n), val)
	}
}

func BenchmarkVenomGet(b *testing.B) {
	envVarsToUnset := make([]string, 0)

	defer func() {
		// clean up any env vars
		for _, v := range envVarsToUnset {
			os.Unsetenv(v)
		}
		envVarsToUnset = []string{}
	}()

	testIO := []struct {
		tc  string
		ven *Venom
		key string
	}{
		{
			tc: "single ConfigLevel with one key/value pair",
			ven: func() *Venom {
				ven := New()
				ven.SetDefault("test", 100)
				return ven
			}(),
			key: "test",
		},
		{
			tc: "many key/value pairs in a single ConfigLevel",
			ven: func() *Venom {
				ven := New()
				for i := 0; i < 10000; i++ {
					ven.SetDefault(fmt.Sprintf("test_%d", i), i)
				}
				return ven
			}(),
			key: "test_5000",
		},
		{
			tc: "many key/value pairs spread across multiple ConfigLevels",
			ven: func() *Venom {
				ven := New()
				var key string
				for i := 0; i < 10000; i++ {
					key = fmt.Sprintf("test_%d", i)
					switch {
					case i%int(FileLevel) == 0:
						ven.SetLevel(FileLevel, key, i)
					case i%int(EnvironmentLevel) == 0:
						os.Setenv(strings.ToUpper(key), fmt.Sprintf("%d", i))
						envVarsToUnset = append(envVarsToUnset, strings.ToUpper(key))
					case i%int(OverrideLevel) == 0:
						ven.SetOverride(key, i)
					default:
						ven.SetDefault(key, i)
					}
				}
				return ven
			}(),
			key: "test_5000",
		},
	}

	for _, test := range testIO {
		b.Run(test.tc, func(b *testing.B) {
			benchmarkVenomRead(test.ven, test.key, b)
		})
	}
}

func BenchmarkVenomWrite(b *testing.B) {
	envVarsToUnset := make([]string, 0)

	defer func() {
		// clean up any env vars
		for _, v := range envVarsToUnset {
			os.Unsetenv(v)
		}
		envVarsToUnset = []string{}
	}()

	testIO := []struct {
		tc    string
		ven   *Venom
		level leveler
		key   keyer
		value interface{}
	}{
		{
			tc:    "single key/value pair in one ConfigLevel",
			ven:   New(),
			level: func(int) ConfigLevel { return DefaultLevel },
			key:   func(int) string { return "test" },
			value: 1000,
		},
		{
			tc:    "many key/value pairs in one ConfigLevel",
			ven:   New(),
			level: func(int) ConfigLevel { return DefaultLevel },
			key:   func(i int) string { return fmt.Sprintf("test_%d", i) },
			value: 1000,
		},
		{
			tc:    "many nested key/value pairs in one ConfigLevel",
			ven:   New(),
			level: func(int) ConfigLevel { return DefaultLevel },
			key:   func(i int) string { return fmt.Sprintf("test_%d.test.%d", i, i) },
			value: 1000,
		},
		{
			tc:  "many key/value pairs in many ConfigLevels",
			ven: New(),
			level: func(i int) ConfigLevel {
				switch {
				case i%int(FileLevel) == 0:
					return FileLevel
				case i%int(EnvironmentLevel) == 0:
					envVarsToUnset = append(envVarsToUnset, strings.ToUpper(fmt.Sprintf("test_%d", i)))
					return EnvironmentLevel
				case i%int(OverrideLevel) == 0:
					return OverrideLevel
				default:
					return DefaultLevel
				}
			},
			key:   func(i int) string { return fmt.Sprintf("test_%d", i) },
			value: 1000,
		},
		{
			tc:  "many nested key/value pairs in many ConfigLevels",
			ven: New(),
			level: func(i int) ConfigLevel {
				switch {
				case i%int(FileLevel) == 0:
					return FileLevel
				case i%int(EnvironmentLevel) == 0:
					envVarsToUnset = append(envVarsToUnset, strings.ToUpper(fmt.Sprintf("test_%d", i)))
					return EnvironmentLevel
				case i%int(OverrideLevel) == 0:
					return OverrideLevel
				default:
					return DefaultLevel
				}
			},
			key:   func(i int) string { return fmt.Sprintf("test_%d.test.%d", i, i) },
			value: 1000,
		},
	}

	for _, test := range testIO {
		b.Run(test.tc, func(b *testing.B) {
			benchmarkVenomWrite(test.ven, test.level, test.key, test.value, b)
		})
	}
}
