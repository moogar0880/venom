package venom

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironment(t *testing.T) {
	testIO := []struct {
		key   string
		value string
	}{
		{"foo", "bar"},
		{"timeout", "10"},
	}

	for _, test := range testIO {
		t.Run(test.key, func(t *testing.T) {
			v := New()

			// set the test value into the environment
			os.Setenv(test.key, test.value)

			// register our test's environment variable and load the environment
			v.RegisterEnv(test.key)
			v.LoadEnvironment()

			// ensure the EnvironmentLevel was properly injected into our config map
			assert.Contains(t, v.config, EnvironmentLevel)

			// ensure we get the expected value back from the environment
			assert.Equal(t, test.value, v.Get(test.key))

			// unset our test key from the environment to keep the next test run
			// clean
			os.Unsetenv(test.key)
		})
	}
}
