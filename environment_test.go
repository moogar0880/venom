package venom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironment(t *testing.T) {
	testIO := []struct {
		tc       string
		key      string
		envVar   string
		value    string
		expect   interface{}
		ok       bool
		resolver Resolver
	}{
		{
			tc:       "should retrieve a basic environment variable",
			key:      "foo",
			envVar:   "FOO",
			value:    "bar",
			resolver: defaultEnvResolver,
			expect:   "bar",
			ok:       true,
		},
		{
			tc:     "should retrieve a prefix environment variable",
			key:    "timeout",
			envVar: "MY_SERVICE_TIMEOUT",
			value:  "10",
			resolver: func() Resolver {
				return &EnvironmentVariableResolver{
					Prefix: "MY_SERVICE",
				}
			}(),
			expect: "10",
			ok:     true,
		},
		{
			tc:     "should fail to retrieve non-prefixed environment variable",
			key:    "timeout",
			envVar: "TIMEOUT",
			value:  "10",
			resolver: func() Resolver {
				return &EnvironmentVariableResolver{
					Prefix: "MY_SERVICE",
				}
			}(),
			expect: nil,
			ok:     false,
		},
	}

	for _, test := range testIO {
		t.Run(test.key, func(t *testing.T) {
			v := New()
			v.RegisterResolver(EnvironmentLevel, test.resolver)

			// set the test value into the environment
			t.Setenv(test.envVar, test.value)

			// ensure we get the expected value back from the environment
			actual, ok := v.Find(test.key)
			assert.Equal(t, test.ok, ok)
			assert.Equal(t, test.expect, actual)
		})
	}
}
