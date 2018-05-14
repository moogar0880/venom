package venom

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type configStruct struct {
	Name    string `venom:"name"`
	Enabled bool   `venom:"enabled"`
	Delay   int    `venom:"delay"`
	Ignored int    `venom:"-"`

	// int
	Int8  int8
	Int16 int16 `venom:"int16"`
	Int32 int32 `venom:"int32"`
	Int64 int64 `venom:"int64"`

	// uint
	Uint   uint   `venom:"uint"`
	Uint8  uint8  `venom:"uint8"`
	Uint16 uint16 `venom:"uint16"`
	Uint32 uint32 `venom:"uint32"`
	Uint64 uint64 `venom:"uint64"`

	// floats
	Float32 float32 `venom:"float32"`
	Float64 float64 `venom:"float64"`

	Opts map[string]string `venom:"opts"`
}

type nestedConfig struct {
	SomethingSpecial string        `venom:"something_special"`
	Basic            configStruct  `venom:"basic"`
	BasicPtr         *configStruct `venom:"basic2"`
}

type sliceConfig struct {
	Strings []string
	Bools   []bool

	// ints
	Ints   []int
	Int8s  []int8
	Int16s []int16
	Int32s []int32
	Int64s []int64

	// uints
	Uints   []uint
	Uint8s  []uint8
	Uint16s []uint16
	Uint32s []uint32
	Uint64s []uint64

	// floats
	Float32s []float32
	Float64s []float64

	// Multi-D slices
	Int2D [][]int
	Int3D [][][]int
}

func TestUnmarshal(t *testing.T) {
	testIO := []struct {
		tc     string
		v      *Venom
		err    error
		expect *configStruct
	}{
		{
			tc:     "should unmarshal with empty venom",
			v:      New(),
			expect: &configStruct{},
		},
		{
			tc: "should ignore hyphen values",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("ignored", 12)
				return ven
			}(),
			expect: &configStruct{},
		},
		{
			tc: "should unmarshal values",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("delay", 12)
				ven.SetDefault("name", "foobar")
				ven.SetDefault("enabled", true)

				// int
				ven.SetDefault("int8", int8(8))
				ven.SetDefault("int16", int16(16))
				ven.SetDefault("int32", int32(32))
				ven.SetDefault("int64", int64(64))

				// uint
				ven.SetDefault("uint", uint(1))
				ven.SetDefault("uint8", uint8(8))
				ven.SetDefault("uint16", uint16(16))
				ven.SetDefault("uint32", uint32(32))
				ven.SetDefault("uint64", uint64(64))

				// floats
				ven.SetDefault("float32", float32(32.0))
				ven.SetDefault("float64", float64(64.0))
				return ven
			}(),
			expect: &configStruct{
				Delay:   12,
				Name:    "foobar",
				Enabled: true,
				Int8:    8,
				Int16:   16,
				Int32:   32,
				Int64:   64,
				Uint:    1,
				Uint8:   8,
				Uint16:  16,
				Uint32:  32,
				Uint64:  64,
				Float32: 32.0,
				Float64: 64.0,
			},
		},
		{
			tc: "should use global venom",
			v: func() *Venom {
				v.Clear()
				v.SetDefault("delay", 12)
				v.SetDefault("name", "foobar")
				return nil
			}(),
			expect: &configStruct{
				Delay: 12,
				Name:  "foobar",
			},
		},
		{
			tc: "should error coercing string from int",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("name", 12)
				return ven
			}(),
			err:    &CoerceErr{From: 12, To: "string"},
			expect: &configStruct{},
		},
		{
			tc: "should error coercing bool from string",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("enabled", "foobar")
				return ven
			}(),
			err:    &CoerceErr{From: "foobar", To: "bool"},
			expect: &configStruct{},
		},
		{
			tc: "should error coercing int from string",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("delay", "foobar")
				return ven
			}(),
			err:    &CoerceErr{From: "foobar", To: "int"},
			expect: &configStruct{},
		},
		{
			tc: "should error coercing int8 from string",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("int8", "foobar")
				return ven
			}(),
			err:    &CoerceErr{From: "foobar", To: "int8"},
			expect: &configStruct{},
		},
		{
			tc: "should error coercing int16 from string",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("int16", "foobar")
				return ven
			}(),
			err:    &CoerceErr{From: "foobar", To: "int16"},
			expect: &configStruct{},
		},
		{
			tc: "should error coercing int32 from string",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("int32", "foobar")
				return ven
			}(),
			err:    &CoerceErr{From: "foobar", To: "int32"},
			expect: &configStruct{},
		},
		{
			tc: "should error coercing int64 from string",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("int64", "foobar")
				return ven
			}(),
			err:    &CoerceErr{From: "foobar", To: "int64"},
			expect: &configStruct{},
		},
		{
			tc: "should error coercing uint from string",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uint", "foobar")
				return ven
			}(),
			err:    &CoerceErr{From: "foobar", To: "uint"},
			expect: &configStruct{},
		},
		{
			tc: "should error coercing uint8 from string",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uint8", "foobar")
				return ven
			}(),
			err:    &CoerceErr{From: "foobar", To: "uint8"},
			expect: &configStruct{},
		},
		{
			tc: "should error coercing uint16 from string",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uint16", "foobar")
				return ven
			}(),
			err:    &CoerceErr{From: "foobar", To: "uint16"},
			expect: &configStruct{},
		},
		{
			tc: "should error coercing uint32 from string",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uint32", "foobar")
				return ven
			}(),
			err:    &CoerceErr{From: "foobar", To: "uint32"},
			expect: &configStruct{},
		},
		{
			tc: "should error coercing uint64 from string",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uint64", "foobar")
				return ven
			}(),
			err:    &CoerceErr{From: "foobar", To: "uint64"},
			expect: &configStruct{},
		},
		{
			tc: "should error coercing float32 from string",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("float32", "foobar")
				return ven
			}(),
			err:    &CoerceErr{From: "foobar", To: "float32"},
			expect: &configStruct{},
		},
		{
			tc: "should error coercing float64 from string",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("float64", "foobar")
				return ven
			}(),
			err:    &CoerceErr{From: "foobar", To: "float64"},
			expect: &configStruct{},
		},
	}

	for _, test := range testIO {
		t.Run(test.tc, func(t *testing.T) {
			var conf configStruct

			err := Unmarshal(test.v, &conf)
			assertEqualErrors(t, test.err, err)
			assert.Equal(t, test.expect, &conf)
		})
	}
}

func TestInvalidUnmarshalError(t *testing.T) {
	testIO := []struct {
		tc   string
		v    *Venom
		conf interface{}
		err  error
	}{
		{
			tc: "should err due to nil config",
			v: func() *Venom {
				return New()
			}(),
			conf: nil,
			err:  &InvalidUnmarshalError{reflect.TypeOf(nil)},
		},
		{
			tc: "should err due to nil config",
			v: func() *Venom {
				return New()
			}(),
			conf: configStruct{},
			err:  &InvalidUnmarshalError{reflect.TypeOf(configStruct{})},
		},
		{
			tc: "should err due to nil config",
			v: func() *Venom {
				return New()
			}(),
			conf: func() interface{} {
				var c *configStruct
				return c
			}(),
			err: func() error {
				var c *configStruct
				return &InvalidUnmarshalError{reflect.TypeOf(c)}
			}(),
		},
	}

	for _, test := range testIO {
		t.Run(test.tc, func(t *testing.T) {
			err := Unmarshal(test.v, test.conf)
			assertEqualErrors(t, test.err, err)
		})
	}
}

func TestUnmarshalNested(t *testing.T) {
	testIO := []struct {
		tc     string
		v      *Venom
		err    error
		expect *nestedConfig
	}{
		{
			tc: "should load nested",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("something_special", "special")
				ven.SetDefault("basic.uint8", uint8(8))
				ven.SetDefault("basic2.uint16", uint16(16))
				return ven
			}(),
			expect: &nestedConfig{
				SomethingSpecial: "special",
				Basic: configStruct{
					Uint8: 8,
				},
				BasicPtr: &configStruct{
					Uint16: 16,
				},
			},
		},
	}

	for _, test := range testIO {
		t.Run(test.tc, func(t *testing.T) {
			conf := nestedConfig{
				BasicPtr: &configStruct{},
			}

			err := Unmarshal(test.v, &conf)
			assertEqualErrors(t, test.err, err)
			assert.Equal(t, test.expect, &conf)
		})
	}
}

func TestUnmarshalSlice(t *testing.T) {
	testIO := []struct {
		tc     string
		v      *Venom
		err    error
		expect *sliceConfig
	}{
		{
			tc: "should unmarshal string slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("strings", []string{"foo", "bar"})
				return ven
			}(),
			expect: &sliceConfig{
				Strings: []string{"foo", "bar"},
			},
		},
		{
			tc: "should unmarshal nil string slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("strings", nil)
				return ven
			}(),
			expect: &sliceConfig{
				Strings: nil,
			},
		},
		{
			tc: "should unmarshal string interface{} slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("strings", []interface{}{"foo", "bar"})
				return ven
			}(),
			expect: &sliceConfig{
				Strings: []string{"foo", "bar"},
			},
		},
		{
			tc: "should unmarshal bool slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("bools", []bool{true, true, false})
				return ven
			}(),
			expect: &sliceConfig{
				Bools: []bool{true, true, false},
			},
		},
		{
			tc: "should unmarshal nil bool slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("bools", nil)
				return ven
			}(),
			expect: &sliceConfig{
				Bools: nil,
			},
		},
		{
			tc: "should unmarshal bool interface{} slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("bools", []interface{}{true, true, false})
				return ven
			}(),
			expect: &sliceConfig{
				Bools: []bool{true, true, false},
			},
		},
		{
			tc: "should unmarshal int slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("ints", []int{0, 1})
				return ven
			}(),
			expect: &sliceConfig{
				Ints: []int{0, 1},
			},
		},
		{
			tc: "should unmarshal nil int slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("ints", nil)
				return ven
			}(),
			expect: &sliceConfig{
				Ints: nil,
			},
		},
		{
			tc: "should unmarshal int interface{} slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("ints", []interface{}{0, 1})
				return ven
			}(),
			expect: &sliceConfig{
				Ints: []int{0, 1},
			},
		},
		{
			tc: "should unmarshal int8 slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("int8s", []int8{2, 3})
				return ven
			}(),
			expect: &sliceConfig{
				Int8s: []int8{2, 3},
			},
		},
		{
			tc: "should unmarshal nil int8 slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("int8s", nil)
				return ven
			}(),
			expect: &sliceConfig{
				Int8s: nil,
			},
		},
		{
			tc: "should unmarshal int8 interface{} slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("int8s", []interface{}{int8(2), int8(3)})
				return ven
			}(),
			expect: &sliceConfig{
				Int8s: []int8{2, 3},
			},
		},
		{
			tc: "should unmarshal int16 slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("int16s", []int16{4, 5})
				return ven
			}(),
			expect: &sliceConfig{
				Int16s: []int16{4, 5},
			},
		},
		{
			tc: "should unmarshal nil int16 slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("int16s", nil)
				return ven
			}(),
			expect: &sliceConfig{
				Int16s: nil,
			},
		},
		{
			tc: "should unmarshal int16 interface{} slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("int16s", []interface{}{int16(2), int16(3)})
				return ven
			}(),
			expect: &sliceConfig{
				Int16s: []int16{2, 3},
			},
		},
		{
			tc: "should unmarshal int32 slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("int32s", []int32{6, 7})
				return ven
			}(),
			expect: &sliceConfig{
				Int32s: []int32{6, 7},
			},
		},
		{
			tc: "should unmarshal nil int32 slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("int32s", nil)
				return ven
			}(),
			expect: &sliceConfig{
				Int32s: nil,
			},
		},
		{
			tc: "should unmarshal int32 interface{} slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("int32s", []interface{}{int32(2), int32(3)})
				return ven
			}(),
			expect: &sliceConfig{
				Int32s: []int32{2, 3},
			},
		},
		{
			tc: "should unmarshal int64 slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("int64s", []int64{8, 9})
				return ven
			}(),
			expect: &sliceConfig{
				Int64s: []int64{8, 9},
			},
		},
		{
			tc: "should unmarshal nil int64 slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("int64s", nil)
				return ven
			}(),
			expect: &sliceConfig{
				Int64s: nil,
			},
		},
		{
			tc: "should unmarshal int64 interface{} slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("int64s", []interface{}{int64(2), int64(3)})
				return ven
			}(),
			expect: &sliceConfig{
				Int64s: []int64{2, 3},
			},
		},
		{
			tc: "should unmarshal uint slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uints", []uint{0, 1})
				return ven
			}(),
			expect: &sliceConfig{
				Uints: []uint{0, 1},
			},
		},
		{
			tc: "should unmarshal nil uint slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uints", nil)
				return ven
			}(),
			expect: &sliceConfig{
				Uints: nil,
			},
		},
		{
			tc: "should unmarshal uint interface{} slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uints", []interface{}{uint(2), uint(3)})
				return ven
			}(),
			expect: &sliceConfig{
				Uints: []uint{2, 3},
			},
		},
		{
			tc: "should unmarshal uint8 slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uint8s", []uint8{2, 3})
				return ven
			}(),
			expect: &sliceConfig{
				Uint8s: []uint8{2, 3},
			},
		},
		{
			tc: "should unmarshal nil uint8 slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uint8s", nil)
				return ven
			}(),
			expect: &sliceConfig{
				Uint8s: nil,
			},
		},
		{
			tc: "should unmarshal uint8 interface{} slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uint8s", []interface{}{uint8(2), uint8(3)})
				return ven
			}(),
			expect: &sliceConfig{
				Uint8s: []uint8{2, 3},
			},
		},
		{
			tc: "should unmarshal uint16 slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uint16s", []uint16{4, 5})
				return ven
			}(),
			expect: &sliceConfig{
				Uint16s: []uint16{4, 5},
			},
		},
		{
			tc: "should unmarshal nil uint16 slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uint16s", nil)
				return ven
			}(),
			expect: &sliceConfig{
				Uint16s: nil,
			},
		},
		{
			tc: "should unmarshal uint16 interface{} slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uint16s", []interface{}{uint16(2), uint16(3)})
				return ven
			}(),
			expect: &sliceConfig{
				Uint16s: []uint16{2, 3},
			},
		},
		{
			tc: "should unmarshal uint32 slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uint32s", []uint32{6, 7})
				return ven
			}(),
			expect: &sliceConfig{
				Uint32s: []uint32{6, 7},
			},
		},
		{
			tc: "should unmarshal nil uint32 slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uint32s", nil)
				return ven
			}(),
			expect: &sliceConfig{
				Uint32s: nil,
			},
		},
		{
			tc: "should unmarshal uint32 interface{} slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uint32s", []interface{}{uint32(2), uint32(3)})
				return ven
			}(),
			expect: &sliceConfig{
				Uint32s: []uint32{2, 3},
			},
		},
		{
			tc: "should unmarshal uint64 slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uint64s", []uint64{8, 9})
				return ven
			}(),
			expect: &sliceConfig{
				Uint64s: []uint64{8, 9},
			},
		},
		{
			tc: "should unmarshal nil uint64 slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uint64s", nil)
				return ven
			}(),
			expect: &sliceConfig{
				Uint64s: nil,
			},
		},
		{
			tc: "should unmarshal uint64 interface{} slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uint64s", []interface{}{uint64(2), uint64(3)})
				return ven
			}(),
			expect: &sliceConfig{
				Uint64s: []uint64{2, 3},
			},
		},
		{
			tc: "should unmarshal float32 slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("float32s", []float32{0.0, 8675.309})
				return ven
			}(),
			expect: &sliceConfig{
				Float32s: []float32{0.0, 8675.309},
			},
		},
		{
			tc: "should unmarshal nil float32 slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("float32s", nil)
				return ven
			}(),
			expect: &sliceConfig{
				Float32s: nil,
			},
		},
		{
			tc: "should unmarshal float32 interface{} slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("float32s", []interface{}{float32(0.0), float32(8675.309)})
				return ven
			}(),
			expect: &sliceConfig{
				Float32s: []float32{0.0, 8675.309},
			},
		},
		{
			tc: "should unmarshal float64 slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("float64s", []float64{0.0, 8675.309})
				return ven
			}(),
			expect: &sliceConfig{
				Float64s: []float64{0.0, 8675.309},
			},
		},
		{
			tc: "should unmarshal nil float64 slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("float64s", nil)
				return ven
			}(),
			expect: &sliceConfig{
				Float64s: nil,
			},
		},
		{
			tc: "should unmarshal float64 interface{} slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("float64s", []interface{}{float64(0.0), float64(8675.309)})
				return ven
			}(),
			expect: &sliceConfig{
				Float64s: []float64{0.0, 8675.309},
			},
		},
		{
			tc: "should unmarshal 2d slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("int2d", [][]int{
					{0, 1, 2, 3, 4, 5},
					{6, 7, 8, 9, 10, 11},
					{12, 13, 14, 15, 16, 17},
				})
				return ven
			}(),
			expect: &sliceConfig{
				Int2D: [][]int{
					{0, 1, 2, 3, 4, 5},
					{6, 7, 8, 9, 10, 11},
					{12, 13, 14, 15, 16, 17},
				},
			},
		},
		{
			tc: "should unmarshal 3d slice",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("int3d", [][][]int{
					{
						{0, 1, 2, 3, 4, 5},
						{6, 7, 8, 9, 10, 11},
						{12, 13, 14, 15, 16, 17},
					},
					{
						{18, 19, 20, 21, 22, 23},
						{24, 25, 26, 27, 28, 29},
						{30, 31, 32, 33, 34, 35},
					},
				})
				return ven
			}(),
			expect: &sliceConfig{
				Int3D: [][][]int{
					{
						{0, 1, 2, 3, 4, 5},
						{6, 7, 8, 9, 10, 11},
						{12, 13, 14, 15, 16, 17},
					},
					{
						{18, 19, 20, 21, 22, 23},
						{24, 25, 26, 27, 28, 29},
						{30, 31, 32, 33, 34, 35},
					},
				},
			},
		},
		{
			tc: "should error coercing string to slice of strings",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("strings", "foobar")
				return ven
			}(),
			expect: new(sliceConfig),
			err:    &CoerceErr{From: "foobar", To: "[]string"},
		},
		{
			tc: "should error coercing interface{} slice with invalid values to slice of strings",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("strings", []interface{}{"foo", true})
				return ven
			}(),
			expect: new(sliceConfig),
			err: &CoerceErr{
				From: []interface{}{"foo", true},
				To:   "[]string",
				Err: &CoerceErr{
					From: true,
					To:   "string",
				},
			},
		},
		{
			tc: "should error coercing string to slice of bools",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("bools", "foobar")
				return ven
			}(),
			expect: new(sliceConfig),
			err:    &CoerceErr{From: "foobar", To: "[]bool"},
		},
		{
			tc: "should error coercing interface{} slice with invalid values to slice of bools",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("bools", []interface{}{true, "false"})
				return ven
			}(),
			expect: new(sliceConfig),
			err: &CoerceErr{
				From: []interface{}{true, "false"},
				To:   "[]bool",
				Err: &CoerceErr{
					From: "false",
					To:   "bool",
				},
			},
		},
		{
			tc: "should error coercing string to slice of ints",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("ints", "foobar")
				return ven
			}(),
			expect: new(sliceConfig),
			err:    &CoerceErr{From: "foobar", To: "[]int"},
		},
		{
			tc: "should error coercing interface{} slice with invalid values to slice of ints",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("ints", []interface{}{0, "foobar"})
				return ven
			}(),
			expect: new(sliceConfig),
			err: &CoerceErr{
				From: []interface{}{0, "foobar"},
				To:   "[]int",
				Err: &CoerceErr{
					From: "foobar",
					To:   "int",
				},
			},
		},
		{
			tc: "should error coercing string to slice of int8s",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("int8s", "foobar")
				return ven
			}(),
			expect: new(sliceConfig),
			err:    &CoerceErr{From: "foobar", To: "[]int8"},
		},
		{
			tc: "should error coercing interface{} slice with invalid values to slice of int8s",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("int8s", []interface{}{int8(0), "foobar"})
				return ven
			}(),
			expect: new(sliceConfig),
			err: &CoerceErr{
				From: []interface{}{int8(0), "foobar"},
				To:   "[]int8",
				Err: &CoerceErr{
					From: "foobar",
					To:   "int8",
				},
			},
		},
		{
			tc: "should error coercing string to slice of int16s",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("int16s", "foobar")
				return ven
			}(),
			expect: new(sliceConfig),
			err:    &CoerceErr{From: "foobar", To: "[]int16"},
		},
		{
			tc: "should error coercing interface{} slice with invalid values to slice of int16s",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("int16s", []interface{}{int16(0), "foobar"})
				return ven
			}(),
			expect: new(sliceConfig),
			err: &CoerceErr{
				From: []interface{}{int16(0), "foobar"},
				To:   "[]int16",
				Err: &CoerceErr{
					From: "foobar",
					To:   "int16",
				},
			},
		},
		{
			tc: "should error coercing string to slice of int32s",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("int32s", "foobar")
				return ven
			}(),
			expect: new(sliceConfig),
			err:    &CoerceErr{From: "foobar", To: "[]int32"},
		},
		{
			tc: "should error coercing interface{} slice with invalid values to slice of int32s",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("int32s", []interface{}{int32(0), "foobar"})
				return ven
			}(),
			expect: new(sliceConfig),
			err: &CoerceErr{
				From: []interface{}{int32(0), "foobar"},
				To:   "[]int32",
				Err: &CoerceErr{
					From: "foobar",
					To:   "int32",
				},
			},
		},
		{
			tc: "should error coercing string to slice of int64s",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("int64s", "foobar")
				return ven
			}(),
			expect: new(sliceConfig),
			err:    &CoerceErr{From: "foobar", To: "[]int64"},
		},
		{
			tc: "should error coercing interface{} slice with invalid values to slice of int64s",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("int64s", []interface{}{int64(0), "foobar"})
				return ven
			}(),
			expect: new(sliceConfig),
			err: &CoerceErr{
				From: []interface{}{int64(0), "foobar"},
				To:   "[]int64",
				Err: &CoerceErr{
					From: "foobar",
					To:   "int64",
				},
			},
		},
		{
			tc: "should error coercing string to slice of uints",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uints", "foobar")
				return ven
			}(),
			expect: new(sliceConfig),
			err:    &CoerceErr{From: "foobar", To: "[]uint"},
		},
		{
			tc: "should error coercing interface{} slice with invalid values to slice of uints",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uints", []interface{}{uint(0), "foobar"})
				return ven
			}(),
			expect: new(sliceConfig),
			err: &CoerceErr{
				From: []interface{}{uint(0), "foobar"},
				To:   "[]uint",
				Err: &CoerceErr{
					From: "foobar",
					To:   "uint",
				},
			},
		},
		{
			tc: "should error coercing string to slice of uint8s",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uint8s", "foobar")
				return ven
			}(),
			expect: new(sliceConfig),
			err:    &CoerceErr{From: "foobar", To: "[]uint8"},
		},
		{
			tc: "should error coercing interface{} slice with invalid values to slice of uint8s",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uint8s", []interface{}{uint8(0), "foobar"})
				return ven
			}(),
			expect: new(sliceConfig),
			err: &CoerceErr{
				From: []interface{}{uint8(0), "foobar"},
				To:   "[]uint8",
				Err: &CoerceErr{
					From: "foobar",
					To:   "uint8",
				},
			},
		},
		{
			tc: "should error coercing string to slice of uint16s",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uint16s", "foobar")
				return ven
			}(),
			expect: new(sliceConfig),
			err:    &CoerceErr{From: "foobar", To: "[]uint16"},
		},
		{
			tc: "should error coercing interface{} slice with invalid values to slice of uint16s",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uint16s", []interface{}{uint16(0), "foobar"})
				return ven
			}(),
			expect: new(sliceConfig),
			err: &CoerceErr{
				From: []interface{}{uint16(0), "foobar"},
				To:   "[]uint16",
				Err: &CoerceErr{
					From: "foobar",
					To:   "uint16",
				},
			},
		},
		{
			tc: "should error coercing string to slice of uint32s",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uint32s", "foobar")
				return ven
			}(),
			expect: new(sliceConfig),
			err:    &CoerceErr{From: "foobar", To: "[]uint32"},
		},
		{
			tc: "should error coercing interface{} slice with invalid values to slice of uint32s",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uint32s", []interface{}{uint32(0), "foobar"})
				return ven
			}(),
			expect: new(sliceConfig),
			err: &CoerceErr{
				From: []interface{}{uint32(0), "foobar"},
				To:   "[]uint32",
				Err: &CoerceErr{
					From: "foobar",
					To:   "uint32",
				},
			},
		},
		{
			tc: "should error coercing string to slice of uint64s",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uint64s", "foobar")
				return ven
			}(),
			expect: new(sliceConfig),
			err:    &CoerceErr{From: "foobar", To: "[]uint64"},
		},
		{
			tc: "should error coercing interface{} slice with invalid values to slice of uint64s",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("uint64s", []interface{}{uint64(0), "foobar"})
				return ven
			}(),
			expect: new(sliceConfig),
			err: &CoerceErr{
				From: []interface{}{uint64(0), "foobar"},
				To:   "[]uint64",
				Err: &CoerceErr{
					From: "foobar",
					To:   "uint64",
				},
			},
		},
		{
			tc: "should error coercing string to slice of float32s",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("float32s", "foobar")
				return ven
			}(),
			expect: new(sliceConfig),
			err:    &CoerceErr{From: "foobar", To: "[]float32"},
		},
		{
			tc: "should error coercing interface{} slice with invalid values to slice of float32s",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("float32s", []interface{}{float32(0), "foobar"})
				return ven
			}(),
			expect: new(sliceConfig),
			err: &CoerceErr{
				From: []interface{}{float32(0), "foobar"},
				To:   "[]float32",
				Err: &CoerceErr{
					From: "foobar",
					To:   "float32",
				},
			},
		},
		{
			tc: "should error coercing string to slice of float64s",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("float64s", "foobar")
				return ven
			}(),
			expect: new(sliceConfig),
			err:    &CoerceErr{From: "foobar", To: "[]float64"},
		},
		{
			tc: "should error coercing interface{} slice with invalid values to slice of float64s",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("float64s", []interface{}{float64(0), "foobar"})
				return ven
			}(),
			expect: new(sliceConfig),
			err: &CoerceErr{
				From: []interface{}{float64(0), "foobar"},
				To:   "[]float64",
				Err: &CoerceErr{
					From: "foobar",
					To:   "float64",
				},
			},
		},
		{
			tc: "should error coercing non-matching slices",
			v: func() *Venom {
				ven := New()
				ven.SetDefault("strings", []float64{0.0})
				return ven
			}(),
			expect: new(sliceConfig),
			err:    &CoerceErr{From: []float64{0}, To: "[]string"},
		},
	}

	for _, test := range testIO {
		t.Run(test.tc, func(t *testing.T) {
			conf := new(sliceConfig)

			err := Unmarshal(test.v, conf)
			assertEqualErrors(t, test.err, err)
			assert.Equal(t, test.expect, conf)
		})
	}
}

type Config struct {
	Ports   []int
	Storage StorageConfig
	Dir     string `venom:"dir"`
}

type StorageConfig struct {
	Driver     string
	DriverOpts map[string]string `venom:"driver_opts"`
}

func TestUnmarshal_Hyper(t *testing.T) {
	inp := Config{
		Ports: []int{8443},
		Storage: StorageConfig{
			DriverOpts: make(map[string]string),
		},
	}

	dir, err := ioutil.TempDir("", "tmp_test_data")
	assert.Nil(t, err)
	defer os.RemoveAll(dir)

	content := `{
    "dir": "/some/path",
    "ports": [443],
    "storage": {
        "driver": "bolt",
        "driver_opts": {
            "path": "/some/path"
        }
    }
}`

	tmpfile := filepath.Join(dir, "test.json")
	err = ioutil.WriteFile(filepath.Join(dir, "test.json"), []byte(content), 0644)
	assert.Nil(t, err)

	v := New()
	assert.Nil(t, v.LoadFile(tmpfile))
	err = Unmarshal(v, &inp)
	assert.Nil(t, err)

	assert.Equal(t, []int{443}, inp.Ports)
	assert.Equal(t, "/some/path", inp.Dir)
}
