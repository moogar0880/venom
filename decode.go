package venom

//go:generate go run generate_coercers.go

import (
	"reflect"
	"strings"
)

const tag = "venom"

// Unmarshal unmarshals the provided Venom config into the provided interface
func Unmarshal(data *Venom, dst interface{}) error {
	// default to the global venom config if the provided Venom is nil
	if data == nil {
		data = v
	}

	var d decoder
	d.init(data)
	return d.unmarshal(dst)
}

// InvalidUnmarshalError is an error returned when an an invalid destination
// value is provided to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "venom: can not Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "venom: can not Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "venom: can not Unmarshal(nil " + e.Type.String() + ")"
}

type namespace struct {
	namespaces []string
}

func (n *namespace) add(s string) {
	n.namespaces = append(n.namespaces, s)
}

func (n *namespace) pop() {
	if len(n.namespaces) == 0 {
		return
	}
	n.namespaces = n.namespaces[:len(n.namespaces)-1]
}

func (n *namespace) String() string {
	return strings.Join(n.namespaces, Delim)
}

type decoder struct {
	data *Venom
	ns   *namespace
}

func (d *decoder) init(data *Venom) *decoder {
	d.data = data
	d.ns = &namespace{
		namespaces: make([]string, 0),
	}
	return d
}

func (d *decoder) unmarshal(dst interface{}) error {
	rv := reflect.ValueOf(dst)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(dst)}
	}
	return d.value(rv)
}

func (d *decoder) value(val reflect.Value) error {
	var err error
	elem := val.Elem()
	typ := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		// pull out the venom struct tag
		elemField := elem.Field(i)
		typField := typ.Field(i)

		fieldTag := typField.Tag.Get(tag)
		if fieldTag == "" {
			// resolvable fields must have at least the `flag` struct tag
			fieldTag = strings.ToLower(typField.Name)
		}

		// determine if this is an unsettable field or was explicitly set to be
		// ignored
		if !elemField.CanSet() || fieldTag == "-" {
			continue
		}

		// add the current namespace to the namespace context
		d.ns.add(fieldTag)

		// fail early if the specified config doesn't exist
		config, ok := d.data.Find(d.ns.String())
		switch {
		case ok:
			if err = d.coerce(config, typField.Type.Kind(), elemField); err != nil {
				return err
			}
		case !ok && elemField.Kind() == reflect.Struct:
			if err := d.value(elemField.Addr()); err != nil {
				return err
			}
		case !ok && elemField.Kind() == reflect.Ptr:
			if err := d.value(elemField); err != nil {
				return err
			}
		}

		// reset the namespace before iterating to the next field
		d.ns.pop()
	}
	return nil
}

// coerce converts the provided query parameter slice into the proper type for
// the target field. this coerced value is then assigned to the current field
func (d *decoder) coerce(val interface{}, to reflect.Kind, field reflect.Value) error {
	var err error

	switch to {
	case reflect.String:
		var actual string
		actual, err = coerceString(val)
		if err != nil {
			return err
		}
		field.SetString(actual)
	case reflect.Bool:
		var actual bool
		actual, err = coerceBool(val)
		if err != nil {
			return err
		}
		field.SetBool(actual)
	case reflect.Int:
		var actual int
		actual, err = coerceInt(val)
		if err != nil {
			return err
		}
		field.SetInt(int64(actual))
	case reflect.Int8:
		var actual int8
		actual, err = coerceInt8(val)
		if err != nil {
			return err
		}
		field.SetInt(int64(actual))
	case reflect.Int16:
		var actual int16
		actual, err = coerceInt16(val)
		if err != nil {
			return err
		}
		field.SetInt(int64(actual))
	case reflect.Int32:
		var actual int32
		actual, err = coerceInt32(val)
		if err != nil {
			return err
		}
		field.SetInt(int64(actual))
	case reflect.Int64:
		var actual int64
		actual, err = coerceInt64(val)
		if err != nil {
			return err
		}
		field.SetInt(actual)
	case reflect.Uint:
		var actual uint
		actual, err = coerceUint(val)
		if err != nil {
			return err
		}
		field.SetUint(uint64(actual))
	case reflect.Uint8:
		var actual uint8
		actual, err = coerceUint8(val)
		if err != nil {
			return err
		}
		field.SetUint(uint64(actual))
	case reflect.Uint16:
		var actual uint16
		actual, err = coerceUint16(val)
		if err != nil {
			return err
		}
		field.SetUint(uint64(actual))
	case reflect.Uint32:
		var actual uint32
		actual, err = coerceUint32(val)
		if err != nil {
			return err
		}
		field.SetUint(uint64(actual))
	case reflect.Uint64:
		var actual uint64
		actual, err = coerceUint64(val)
		if err != nil {
			return err
		}
		field.SetUint(actual)
	case reflect.Float32:
		var actual float32
		actual, err = coerceFloat32(val)
		if err != nil {
			return err
		}
		field.SetFloat(float64(actual))
	case reflect.Float64:
		var actual float64
		actual, err = coerceFloat64(val)
		if err != nil {
			return err
		}
		field.SetFloat(actual)
	case reflect.Struct:
		if field.CanAddr() {
			err = d.value(field.Addr())
		}
	case reflect.Ptr:
		err = d.value(field)
	case reflect.Slice:
		err = d.coerceSlice(val, to, field)
	}
	return err
}

// coerceSlice creates a new slice of the appropriate type for the target field
// and coerces each of the query parameter values into the destination type.
// Should any of the provided query parameters fail to be coerced, an error is
// returned and the entire slice will not be applied
func (d *decoder) coerceSlice(val interface{}, to reflect.Kind, field reflect.Value) error {
	var err error
	sliceType := field.Type().Elem()
	coerceKind := sliceType.Kind()

	switch coerceKind {
	case reflect.String:
		var actual []string
		actual, err = coerceStringSlice(val)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(actual))
	case reflect.Bool:
		var actual []bool
		actual, err = coerceBoolSlice(val)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(actual))
	case reflect.Int:
		var actual []int
		actual, err = coerceIntSlice(val)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(actual))
	case reflect.Int8:
		var actual []int8
		actual, err = coerceInt8Slice(val)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(actual))
	case reflect.Int16:
		var actual []int16
		actual, err = coerceInt16Slice(val)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(actual))
	case reflect.Int32:
		var actual []int32
		actual, err = coerceInt32Slice(val)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(actual))
	case reflect.Int64:
		var actual []int64
		actual, err = coerceInt64Slice(val)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(actual))
	case reflect.Uint:
		var actual []uint
		actual, err = coerceUintSlice(val)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(actual))
	case reflect.Uint8:
		var actual []uint8
		actual, err = coerceUint8Slice(val)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(actual))
	case reflect.Uint16:
		var actual []uint16
		actual, err = coerceUint16Slice(val)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(actual))
	case reflect.Uint32:
		var actual []uint32
		actual, err = coerceUint32Slice(val)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(actual))
	case reflect.Uint64:
		var actual []uint64
		actual, err = coerceUint64Slice(val)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(actual))
	case reflect.Float32:
		var actual []float32
		actual, err = coerceFloat32Slice(val)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(actual))
	case reflect.Float64:
		var actual []float64
		actual, err = coerceFloat64Slice(val)
		if err != nil {
			return err
		}
		field.Set(reflect.ValueOf(actual))
	case reflect.Struct:
		unAddressableSlice := reflect.MakeSlice(field.Type(), field.Len(), field.Cap())
		actual := reflect.New(unAddressableSlice.Type())

		sliceVal := reflect.ValueOf(val)
		for i := 0; i < sliceVal.Len(); i++ {
			item := reflect.New(sliceType)
			if item.CanAddr() {
				if err = d.value(item.Addr()); err != nil {
					return err
				}
			}
			reflectAppend(actual, item)
		}
		field.Set(reflect.Indirect(actual))
	case reflect.Slice:
		// we've hit a multidimensional slice so we need to recursively build
		// up the inner slices
		slice := reflect.ValueOf(val)

		fld := field.Addr().Elem()
		innerSliceType := reflect.TypeOf(field.Interface()).Elem()
		innerSliceKind := innerSliceType.Kind()

		for i := 0; i < slice.Len(); i++ {
			// allocate new inner slice
			fld.Set(reflect.Append(fld, reflect.Indirect(reflect.New(innerSliceType))))
			if err = d.coerceSlice(slice.Index(i).Interface(), innerSliceKind, fld.Index(i)); err != nil {
				return err
			}
		}
	}
	return err
}

func reflectAppend(slice reflect.Value, item reflect.Value) {
	sliceVal := reflect.ValueOf(slice.Interface())
	sliceVal.Elem().Set(reflect.Append(sliceVal.Elem(), reflect.Indirect(item)))
}
