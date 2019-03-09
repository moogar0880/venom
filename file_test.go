package venom

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterExtension(t *testing.T) {
	testIO := []struct {
		tc     string
		key    string
		loader IOFileLoader
	}{
		{
			tc:     "should load JSONLoader",
			key:    jsonKey,
			loader: JSONLoader,
		},
		{
			tc:     "should load XMLLoader",
			key:    jsonKey,
			loader: JSONLoader,
		},
	}

	for _, test := range testIO {
		t.Run(test.tc, func(t *testing.T) {
			RegisterExtension(test.key, test.loader)
			assert.Contains(t, extensionMap, test.key)
		})
	}
}

func getJSONFileErr(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	d := make(map[string]interface{})
	return json.Unmarshal(data, &d)
}

func TestLoadFile(t *testing.T) {
	testIO := []struct {
		tc       string
		filename string
		err      error
		expect   ConfigMap
	}{
		{
			tc:       "should load JSON file",
			filename: "testdata/config.json",
			expect: ConfigMap{
				"foo":   "bar",
				"level": 5.0,
			},
		},
		{
			tc:       "should error on non-existent file",
			filename: "testdata/missing.config.json",
			err: &os.PathError{
				Op:   "open",
				Path: "testdata/missing.config.json",
				Err:  syscall.ENOENT,
			},
			expect: nil,
		},
		{
			tc:       "should error on unknown file extension",
			filename: "testdata/config.xml",
			err:      ErrNoFileLoader{ext: "xml"},
			expect:   nil,
		},
		{
			tc:       "should error on invalid file contents",
			filename: "testdata/invalid/config.bad.json",
			err:      getJSONFileErr("testdata/invalid/config.bad.json"),
			expect:   nil,
		},
	}

	for _, test := range testIO {
		t.Run(test.tc, func(t *testing.T) {
			v := New()
			err := v.LoadFile(test.filename)

			assertEqualErrors(t, test.err, err)
			st := v.Store.(*DefaultConfigStore)
			assert.EqualValues(t, st.config[FileLevel], test.expect)
		})
	}
}

func TestLoadDirectory(t *testing.T) {
	testIO := []struct {
		tc      string
		dir     string
		recurse bool
		err     error
		expect  ConfigMap
	}{
		{
			tc:      "should load single directory",
			dir:     "testdata",
			recurse: false,
			expect: ConfigMap{
				"foo":   "bar",
				"level": 5.0,
			},
		},
		{
			tc:      "should recursively load directories",
			dir:     "testdata/sub",
			recurse: true,
			expect: ConfigMap{
				"foo":   "baz",
				"level": 5.0,
				"and":   "another",
			},
		},
		{
			tc:      "should error if directory contains invalid files",
			dir:     "testdata/invalid",
			recurse: false,
			err:     getJSONFileErr("testdata/invalid/config.bad.json"),
		},
	}

	for _, test := range testIO {
		t.Run(test.tc, func(t *testing.T) {
			v := New()
			err := v.LoadDirectory(test.dir, test.recurse)

			assertEqualErrors(t, test.err, err)
			st := v.Store.(*DefaultConfigStore)
			assert.EqualValues(t, st.config[FileLevel], test.expect)
		})
	}
}
